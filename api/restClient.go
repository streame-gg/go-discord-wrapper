package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/util"
)

// routePathKey is the context key used to carry the relative route path through
// the request lifecycle so the rate limiter can key on it without re-parsing the URL.
type routePathKey struct{}

type RestEventType string

const (
	RestEventRequest   RestEventType = "REQUEST"
	RestEventResponse  RestEventType = "RESPONSE"
	RestEventRetry     RestEventType = "RETRY"
	RestEventRateLimit RestEventType = "RATE_LIMIT"
	RestEventError     RestEventType = "ERROR"
)

type RestEvent struct {
	Type       RestEventType
	Request    *http.Request
	Response   *http.Response
	Attempt    int
	RetryAfter time.Duration
	Err        error
}

type RestEventHandler func(*RestClient, RestEvent)

type RestClient struct {
	BaseURL string
	token   string
	Version common.APIVersion

	httpClient *http.Client

	retryOptions options.RetryOptions

	// minRequestInterval enforces a crude global delay between all requests.
	// This is separate from — and complementary to — the proactive rate limiter.
	minRequestInterval time.Duration
	rateLimitMu        sync.Mutex
	nextRequestAt      time.Time

	// rateLimiter implements proactive per-route and global rate limiting.
	// nil when proactive throttling has been disabled via WithRateLimiting.
	rateLimiter *rateLimiter

	eventEmitter *util.EventEmitter[RestEventType, RestEventHandler]
}

// NewRestClient creates a REST client for the Discord API.
// Configure it using options from the options package:
//
//	rc, err := api.NewRestClient(token,
//	    options.WithRetry(options.RetryOptions{MaxRetries: 5}),
//	    options.WithRateLimiting(options.RateLimiterOptions{SafetyMargin: 1}),
//	)
func NewRestClient(token string, opts ...options.Option) (*RestClient, error) {
	cfg := options.Build(options.Config{
		BaseURL:    "https://discord.com/api",
		APIVersion: common.APIVersion10,
		Retry: options.RetryOptions{
			MaxRetries:          3,
			BaseBackoff:         500 * time.Millisecond,
			MaxBackoff:          5 * time.Second,
			RetryOnRateLimit:    true,
			RetryOnServerErrors: true,
		},
	}, opts)
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("go-discord-wrapper: %w", err)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	// Build rate limiter. Proactive throttling is on by default unless explicitly
	// disabled via WithRateLimiting(RateLimiterOptions{Disabled: true}).
	var rl *rateLimiter
	if cfg.RateLimiter == nil || !cfg.RateLimiter.Disabled {
		safetyMargin := 0
		if cfg.RateLimiter != nil {
			safetyMargin = cfg.RateLimiter.SafetyMargin
		}
		rl = newRateLimiter(safetyMargin)
	}

	return &RestClient{
		BaseURL:            cfg.BaseURL,
		token:              token,
		Version:            cfg.APIVersion,
		httpClient:         httpClient,
		retryOptions:       cfg.Retry,
		minRequestInterval: cfg.MinRequestInterval,
		rateLimiter:        rl,
		eventEmitter:       util.NewEventEmitter[RestEventType, RestEventHandler](),
	}, nil
}

func (c *RestClient) buildURL() string {
	return c.BaseURL + "/v" + c.Version.ToString()
}

// OnEvent registers a handler for REST lifecycle events (request, response, retry, etc.).
func (c *RestClient) OnEvent(eventType RestEventType, handler RestEventHandler) {
	if c.eventEmitter == nil {
		c.eventEmitter = util.NewEventEmitter[RestEventType, RestEventHandler]()
	}

	c.eventEmitter.On(eventType, handler)
}

// redactWebhookToken returns a clone of req with the webhook token replaced by
// "[REDACTED]" so that lifecycle hook payloads never expose credentials.
// Paths of the form /webhooks/{id}/{token}/... are detected by structure.
func redactWebhookToken(req *http.Request) *http.Request {
	if req == nil || req.URL == nil {
		return req
	}
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) < 3 || parts[0] != "webhooks" {
		return req
	}
	clone := req.Clone(req.Context())
	urlCopy := *req.URL
	redacted := make([]string, len(parts))
	copy(redacted, parts)
	redacted[2] = "[REDACTED]"
	urlCopy.Path = "/" + strings.Join(redacted, "/")
	clone.URL = &urlCopy
	return clone
}

func (c *RestClient) emitEvent(event RestEvent) {
	if c.eventEmitter == nil {
		return
	}

	for _, handler := range c.eventEmitter.Handlers(event.Type) {
		handler(c, event)
	}
}

// generateRequest builds an authenticated HTTP request and embeds the relative
// route path into the request context for the rate limiter to consume.
func (c *RestClient) generateRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	// Store the relative path (e.g. "/channels/111/messages") so that the
	// rate limiter can normalise it into a bucket key without re-parsing the URL.
	ctx = context.WithValue(ctx, routePathKey{}, path)

	req, err := http.NewRequestWithContext(ctx, method, c.buildURL()+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bot "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("GoDiscordWrapper (%s@%s)", common.RepositoryURL, common.RepositoryVersion))

	return req, nil
}

func (c *RestClient) do(req *http.Request, successResponseCode int, v interface{}) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("request must not be nil")
	}

	// Extract the relative path stored by generateRequest for rate limiting.
	routePath, _ := req.Context().Value(routePathKey{}).(string)

	bodyBytes, err := captureRequestBody(req)
	if err != nil {
		return nil, err
	}

	maxAttempts := c.retryOptions.MaxRetries + 1

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Crude global request interval (legacy option, complementary to rate limiter).
		if err := c.waitForMinInterval(req.Context()); err != nil {
			return nil, err
		}

		// Proactive rate limit check — blocks if the route's bucket is exhausted
		// or the global rate limit is active.
		if c.rateLimiter != nil && routePath != "" {
			c.rateLimiter.wait(req.Method, routePath)
		}

		attemptReq, err := cloneRequest(req, bodyBytes)
		if err != nil {
			return nil, err
		}

		safeReq := redactWebhookToken(attemptReq)
		c.emitEvent(RestEvent{Type: RestEventRequest, Request: safeReq, Attempt: attempt})

		resp, reqErr := c.httpClient.Do(attemptReq)
		if reqErr != nil {
			c.emitEvent(RestEvent{Type: RestEventError, Request: safeReq, Attempt: attempt, Err: reqErr})
			if attempt == maxAttempts {
				return nil, reqErr
			}

			delay := c.retryDelay(attempt, 0)
			c.emitEvent(RestEvent{Type: RestEventRetry, Request: safeReq, Attempt: attempt, RetryAfter: delay, Err: reqErr})
			select {
			case <-time.After(delay):
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
			continue
		}

		// Update rate limit state from response headers on every response,
		// including 429s, so bucket state stays accurate for future requests.
		if c.rateLimiter != nil && routePath != "" {
			c.rateLimiter.update(req.Method, routePath, resp)
		}

		c.emitEvent(RestEvent{Type: RestEventResponse, Request: safeReq, Response: resp, Attempt: attempt})

		if resp.StatusCode == successResponseCode {
			if v != nil && resp.StatusCode != http.StatusNoContent {
				if decodeErr := json.NewDecoder(resp.Body).Decode(v); decodeErr != nil {
					_ = resp.Body.Close()
					return nil, decodeErr
				}
			}

			_ = resp.Body.Close()
			return resp, nil
		}

		if attempt < maxAttempts && c.shouldRetry(resp.StatusCode) {
			retryAfter := parseRetryAfter(resp)
			delay := c.retryDelay(attempt, retryAfter)
			if retryAfter > 0 {
				c.emitEvent(RestEvent{Type: RestEventRateLimit, Request: safeReq, Response: resp, Attempt: attempt, RetryAfter: retryAfter})
			}

			c.emitEvent(RestEvent{Type: RestEventRetry, Request: safeReq, Response: resp, Attempt: attempt, RetryAfter: delay})
			_ = resp.Body.Close()
			select {
			case <-time.After(delay):
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
			continue
		}

		respErr := decodeGatewayError(resp)
		_ = resp.Body.Close()
		c.emitEvent(RestEvent{Type: RestEventError, Request: safeReq, Response: resp, Attempt: attempt, Err: respErr})
		return nil, respErr
	}

	return nil, errors.New("request failed after retries")
}

func (c *RestClient) shouldRetry(statusCode int) bool {
	if statusCode == http.StatusTooManyRequests {
		return c.retryOptions.RetryOnRateLimit
	}

	return c.retryOptions.RetryOnServerErrors && statusCode >= http.StatusInternalServerError
}

func (c *RestClient) retryDelay(attempt int, retryAfter time.Duration) time.Duration {
	if retryAfter > 0 {
		return retryAfter
	}

	delay := c.retryOptions.BaseBackoff
	if delay <= 0 {
		delay = 100 * time.Millisecond
	}

	for i := 1; i < attempt; i++ {
		delay *= 2
		if c.retryOptions.MaxBackoff > 0 && delay > c.retryOptions.MaxBackoff {
			return c.retryOptions.MaxBackoff
		}
	}

	if c.retryOptions.MaxBackoff > 0 && delay > c.retryOptions.MaxBackoff {
		return c.retryOptions.MaxBackoff
	}

	return delay
}

// waitForMinInterval enforces the optional MinRequestInterval global delay.
func (c *RestClient) waitForMinInterval(ctx context.Context) error {
	if c.minRequestInterval <= 0 {
		return nil
	}

	c.rateLimitMu.Lock()
	wait := time.Until(c.nextRequestAt)
	if wait < 0 {
		wait = 0
	}
	c.nextRequestAt = time.Now().Add(c.minRequestInterval)
	c.rateLimitMu.Unlock()

	if wait > 0 {
		c.emitEvent(RestEvent{Type: RestEventRateLimit, RetryAfter: wait})
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func decodeGatewayError(resp *http.Response) error {
	var body common.GatewayError
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return &Error{
		HTTPStatus: resp.StatusCode,
		Code:       common.GatewayErrorCode(body.Code),
		Message:    body.Message,
		Errors:     body.Errors,
	}
}

func captureRequestBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	if err := req.Body.Close(); err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}

	return body, nil
}

func cloneRequest(req *http.Request, body []byte) (*http.Request, error) {
	cloned := req.Clone(req.Context())
	if len(body) == 0 {
		return cloned, nil
	}

	cloned.Body = io.NopCloser(bytes.NewReader(body))
	cloned.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}
	cloned.ContentLength = int64(len(body))

	return cloned, nil
}
