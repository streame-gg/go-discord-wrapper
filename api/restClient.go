// Package api implements the Discord REST API client with proactive rate limiting,
// configurable retry logic, and full coverage of Discord API v10 endpoints.
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

// Close stops the background rate-limiter cleanup goroutine. Call it when the
// RestClient is no longer needed to prevent a goroutine leak.
func (c *RestClient) Close() {
	if c.rateLimiter != nil {
		c.rateLimiter.close()
	}
}

// OnEvent registers a handler for REST lifecycle events (request, response, retry, etc.).
func (c *RestClient) OnEvent(eventType RestEventType, handler RestEventHandler) {
	if c.eventEmitter == nil {
		c.eventEmitter = util.NewEventEmitter[RestEventType, RestEventHandler]()
	}

	c.eventEmitter.On(eventType, handler)
}

// redactRequestSecrets returns a clone of req with the webhook token replaced by
// "[REDACTED]" so that lifecycle hook payloads never expose credentials.
// Paths of the form /webhooks/{id}/{token}/... are detected by structure.
func redactRequestSecrets(req *http.Request) *http.Request {
	if req == nil || req.URL == nil {
		return req
	}

	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(parts) == 0 {
		return req
	}

	tokenIndex := -1

	for i, part := range parts {
		switch part {
		case "webhooks":
			// /webhooks/{webhook_id}/{webhook_token}
			if i+2 < len(parts) {
				tokenIndex = i + 2
			}
		case "interactions":
			// /interactions/{interaction_id}/{interaction_token}/callback
			if i+3 < len(parts) && parts[i+3] == "callback" {
				tokenIndex = i + 2
			}
		}

		if tokenIndex != -1 {
			break
		}
	}

	if tokenIndex == -1 {
		return req
	}

	clone := req.Clone(req.Context())
	urlCopy := *req.URL

	redacted := make([]string, len(parts))
	copy(redacted, parts)
	redacted[tokenIndex] = "[REDACTED]"

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

// validateAPIPath rejects any digit-only path segment that is not a valid Discord Snowflake
// (15–20 decimal digits). This stops the most common URL-injection pattern where user-supplied
// short numeric input (e.g. "123") is concatenated directly into a path without sanitisation.
// Slash-injection (e.g. Snowflake("123456789012345/evil")) is caught earlier by Snowflake.String(),
// which panics when the value contains a '/'.
func validateAPIPath(path string) error {
	for i, seg := range strings.Split(strings.Trim(path, "/"), "/") {
		if seg == "" {
			return fmt.Errorf("go-discord-wrapper: path %q contains an empty segment at position %d", path, i)
		}
		if isAllDigits(seg) && !isSnowflakeID(seg) {
			return fmt.Errorf("go-discord-wrapper: path segment %q is not a valid Discord Snowflake (must be 15–20 decimal digits)", seg)
		}
	}
	return nil
}

// isAllDigits reports whether s is non-empty and consists only of ASCII decimal digits.
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (c *RestClient) WithBotAuthorization() string {
	return "Bot " + c.token
}

func WithUserAuthorization(token string) string {
	return "Bearer " + token
}

func WithoutAuthorization() string {
	return ""
}

// generateRequest builds an authenticated HTTP request and embeds the relative
// route path into the request context for the rate limiter to consume.
func (c *RestClient) generateRequest(ctx context.Context, method, path string, body io.Reader, auth string) (*http.Request, error) {
	if err := validateAPIPath(path); err != nil {
		return nil, err
	}
	// Store the relative path (e.g. "/channels/111/messages") so that the
	// rate limiter can normalise it into a bucket key without re-parsing the URL.
	ctx = context.WithValue(ctx, routePathKey{}, path)

	req, err := http.NewRequestWithContext(ctx, method, c.buildURL()+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("GoDiscordWrapper (%s@%s)", common.RepositoryURL, common.RepositoryVersion))

	if auth != "" {
		req.Header.Set("Authorization", "Bot "+c.token)
	}

	return req, nil
}

func (c *RestClient) includesInIntList(actual int, list []int) bool {
	for _, code := range list {
		if code == actual {
			return true
		}
	}

	return false
}

// doRequest executes req and decodes a successful response into v (when v != nil).
// The returned *http.Response has its Body already closed; callers must not
// read from it. Inspect headers and status codes via the returned value, but
// doRequest not call resp.Body.Read or resp.Body.Close again.

// successResponseCodeData is based on map[statusCode]returnRequestBody
// if statusCode is 204 (No Content), no request body will be returned, so you do not have to set it to false

// Waiting for Go 1.27 to implement this in the RestClient;
// https://github.com/golang/go/issues/77273

func doRequest[T any](c *RestClient, req *http.Request, successResponseCodeData map[int]bool) (*T, error) {
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
			if err := c.rateLimiter.wait(req.Context(), req.Method, routePath); err != nil {
				return nil, err
			}
		}

		attemptReq, err := cloneRequest(req, bodyBytes)
		if err != nil {
			return nil, err
		}

		safeReq := redactRequestSecrets(attemptReq)
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

		if entry, ok := successResponseCodeData[resp.StatusCode]; ok {
			if resp.StatusCode == http.StatusNoContent || !entry {
				return nil, nil
			}

			var returnType T

			if json.NewDecoder(resp.Body).Decode(&returnType) != nil {
				return nil, err
			}

			_ = resp.Body.Close()

			return &returnType, nil
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
