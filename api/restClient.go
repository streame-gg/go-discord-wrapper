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
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/util"
)

// routePathKey is the context key used to carry the relative route path through
// the request lifecycle so the rate limiter can key on it without re-parsing the URL.
type routePathKey struct{}

type RestEventType string

type NoReturnData struct{}

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
	// BaseURL is the root of the Discord REST API (default: "https://discord.com/api").
	// It is read by buildURL on every request without a lock. If you need to
	// override it, do so before the client is used for the first time; mutating
	// it concurrently with in-flight requests is a data race.
	BaseURL string
	token   string
	Version discord.APIVersion

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

	maxResponseBodySize int64
}

// NewRestClient creates a REST client for the Discord API.
// Configure it using options from the options package:
//
//	rc, err := api.NewRestClient(token,
//	    options.WithRetry(options.RetryOptions{MaxRetries: 5}),
//	    options.WithRateLimiting(options.RateLimiterOptions{SafetyMargin: 1}),
//	)
func NewRestClient(token string, opts ...options.Option) (*RestClient, error) {
	token = strings.TrimSpace(token)
	if len(token) == 0 {
		return nil, errors.New("token required")
	}

	if parts := strings.SplitN(token, " ", 2); strings.ToLower(parts[0]) == "bot" {
		slog.Info("NewRestClient: token appears to be a bot token with 'Bot' prefix; the library will add this prefix automatically, so you should pass the raw token without 'Bot '")
		if len(parts) == 2 {
			token = strings.TrimSpace(parts[1])
		} else {
			token = ""
		}
	}

	if token == "" {
		return nil, errors.New("go-discord-wrapper: token must not be empty")
	}

	cfg := options.Build(options.Config{
		BaseURL:    "https://discord.com/api",
		APIVersion: discord.APIVersion10,
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

	maxBody := cfg.MaxResponseBodySize
	if maxBody == 0 {
		maxBody = 50 << 20
	}

	return &RestClient{
		BaseURL:             cfg.BaseURL,
		token:               token,
		Version:             cfg.APIVersion,
		httpClient:          httpClient,
		retryOptions:        cfg.Retry,
		minRequestInterval:  cfg.MinRequestInterval,
		rateLimiter:         rl,
		eventEmitter:        util.NewEventEmitter[RestEventType, RestEventHandler](),
		maxResponseBodySize: maxBody,
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

// redactRequestSecrets returns a clone of req with credentials removed so that
// lifecycle hook payloads never expose secrets. It redacts:
//   - The Authorization header (contains the bot token)
//   - Webhook/interaction tokens embedded in URL paths
func redactRequestSecrets(req *http.Request) *http.Request {
	if req == nil || req.URL == nil {
		return req
	}

	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")

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

	hasAuth := req.Header.Get("Authorization") != ""
	if tokenIndex == -1 && !hasAuth {
		return req
	}

	clone := req.Clone(req.Context())

	if hasAuth {
		clone.Header.Set("Authorization", "[REDACTED]")
	}

	if tokenIndex != -1 {
		urlCopy := *req.URL
		redacted := make([]string, len(parts))
		copy(redacted, parts)
		redacted[tokenIndex] = "[REDACTED]"
		urlCopy.Path = "/" + strings.Join(redacted, "/")
		clone.URL = &urlCopy
	}

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
// validateAPIPath is the only line of defense against path injection via Snowflake values.
// Each path segment that looks like a digit-only ID is validated against Snowflake.Validate().
// Non-numeric segments are accepted as-is (route literals).
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

func (c *RestClient) WithBotAuthorization() func(req *http.Request) {
	return func(req *http.Request) {
		req.Header.Set("Authorization", "Bot "+c.token)
	}
}

func WithUserAuthorization(token string) func(req *http.Request) {
	return func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

// WithAuditLogReason attaches an audit-log reason to the request.
// The reason is URL-encoded so non-ASCII characters are transmitted safely.
// Reasons longer than 512 characters (Discord's limit) are silently truncated.
func WithAuditLogReason(reason string) func(req *http.Request) {
	if reason == "" {
		return func(_ *http.Request) {}
	}
	runes := []rune(reason)
	if len(runes) > 512 {
		runes = runes[:512]
	}
	encoded := url.PathEscape(string(runes))
	return func(req *http.Request) {
		req.Header.Set("X-Audit-Log-Reason", encoded)
	}
}

// generateRequest builds an authenticated HTTP request and embeds the relative
// route path into the request context for the rate limiter to consume.
func (c *RestClient) generateRequest(ctx context.Context, method, path string, body io.Reader, options ...func(req *http.Request)) (*http.Request, error) {
	if err := validateAPIPath(path); err != nil {
		return nil, err
	}
	// Store the relative path (e.g. "/channels/111/messages") so that the
	// rate limiter can normalize it into a bucket key without reparsing the URL.
	ctx = context.WithValue(ctx, routePathKey{}, path)

	req, err := http.NewRequestWithContext(ctx, method, c.buildURL()+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("GoDiscordWrapper (%s@%s)", discord.RepositoryURL, discord.RepositoryVersion))

	for _, option := range options {
		option(req)
	}

	return req, nil
}

// doRequest executes req and decodes a successful response into v (when v != nil).
// The returned *http.Response has its Body already closed; callers must not
// read from it. Inspect headers and status codes via the returned value, but
// do not call resp.Body.Read or resp.Body.Close again.

// successResponseCodeData is based on map[statusCode]returnRequestBody
// if statusCode is 204 (No Content), no request body will be returned, so you do not have to set it to false

// Waiting for Go 1.27 to implement this in the RestClient;
// https://github.com/golang/go/issues/77273

func doRequestWithoutResponse(c *RestClient, req *http.Request, enforceStatusCodes ...int) error {
	if len(enforceStatusCodes) == 0 {
		_, err := doRequest[NoReturnData](c, req, map[int]bool{
			http.StatusOK:        false,
			http.StatusCreated:   false,
			http.StatusAccepted:  false,
			http.StatusNoContent: false,
		})

		return err
	}

	statusCodes := map[int]bool{}

	for _, code := range enforceStatusCodes {
		statusCodes[code] = false
	}

	_, err := doRequest[NoReturnData](c, req, statusCodes)

	return err
}

// doRawBytes runs a request through the full rate-limit/retry pipeline and
// returns the raw response body on success (HTTP 200). It is intended for
// endpoints that return non-JSON bodies (e.g. PNG images).
func doRawBytes(c *RestClient, req *http.Request) ([]byte, error) {
	if req == nil {
		return nil, errors.New("request must not be nil")
	}

	routePath, _ := req.Context().Value(routePathKey{}).(string)

	bodyBytes, err := captureRequestBody(req)
	if err != nil {
		return nil, err
	}

	maxAttempts := c.retryOptions.MaxRetries + 1

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := c.waitForMinInterval(req.Context()); err != nil {
			return nil, err
		}
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
			select {
			case <-time.After(delay):
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
			continue
		}

		buf, readErr := io.ReadAll(io.LimitReader(resp.Body, c.maxResponseBodySize))
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("read response body: %w", readErr)
		}
		resp.Body = io.NopCloser(bytes.NewReader(buf))

		if c.rateLimiter != nil && routePath != "" {
			c.rateLimiter.update(req.Method, routePath, resp)
		}

		c.emitEvent(RestEvent{Type: RestEventResponse, Request: safeReq, Response: snapshotResponse(resp, buf), Attempt: attempt})

		if resp.StatusCode == http.StatusOK {
			return buf, nil
		}

		if attempt < maxAttempts && c.shouldRetry(resp.StatusCode) {
			retryAfter := parseRetryAfter(resp)
			delay := c.retryDelay(attempt, retryAfter)
			select {
			case <-time.After(delay):
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
			continue
		}

		return nil, decodeGatewayErrorFromBytes(buf, resp)
	}

	return nil, errors.New("request failed after retries")
}

// doRequestSlice is like doRequest but dereferences the pointer so callers
// receive a plain slice instead of a pointer-to-slice.
func doRequestSlice[T any](c *RestClient, req *http.Request, codes map[int]bool) ([]*T, error) {
	result, err := doRequest[[]*T](c, req, codes)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return *result, nil
}

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

		// Buffer the response body so that event handlers and the decoder both see
		// the full body without competing for the live network reader.
		bodyBuf, readErr := io.ReadAll(io.LimitReader(resp.Body, c.maxResponseBodySize))
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("read response body: %w", readErr)
		}
		resp.Body = io.NopCloser(bytes.NewReader(bodyBuf))

		// Update rate limit state from response headers on every response,
		// including 429s, so bucket state stays accurate for future requests.
		if c.rateLimiter != nil && routePath != "" {
			c.rateLimiter.update(req.Method, routePath, resp)
		}

		// Event handlers receive a snapshot with a fresh body reader so they cannot
		// interfere with the decoder below even if they read or close the body.
		c.emitEvent(RestEvent{Type: RestEventResponse, Request: safeReq, Response: snapshotResponse(resp, bodyBuf), Attempt: attempt})

		if entry, ok := successResponseCodeData[resp.StatusCode]; ok {
			var returnType T

			if resp.StatusCode == http.StatusNoContent || reflect.TypeOf(returnType) == reflect.TypeOf(NoReturnData{}) || !entry {
				return nil, nil
			}

			if decErr := json.NewDecoder(bytes.NewReader(bodyBuf)).Decode(&returnType); decErr != nil {
				return nil, fmt.Errorf("decode response body: %w", decErr)
			}

			return &returnType, nil
		}

		if attempt < maxAttempts && c.shouldRetry(resp.StatusCode) {
			retryAfter := parseRetryAfter(resp)
			delay := c.retryDelay(attempt, retryAfter)
			if retryAfter > 0 {
				c.emitEvent(RestEvent{Type: RestEventRateLimit, Request: safeReq, Response: snapshotResponse(resp, bodyBuf), Attempt: attempt, RetryAfter: retryAfter})
			}

			c.emitEvent(RestEvent{Type: RestEventRetry, Request: safeReq, Response: snapshotResponse(resp, bodyBuf), Attempt: attempt, RetryAfter: delay})
			select {
			case <-time.After(delay):
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
			continue
		}

		respErr := decodeGatewayErrorFromBytes(bodyBuf, resp)
		c.emitEvent(RestEvent{Type: RestEventError, Request: safeReq, Response: snapshotResponse(resp, bodyBuf), Attempt: attempt, Err: respErr})
		return nil, respErr
	}

	return nil, errors.New("request failed after retries")
}

func (c *RestClient) shouldRetry(statusCode int) bool {
	if statusCode == http.StatusTooManyRequests {
		return c.retryOptions.RetryOnRateLimit
	}

	if statusCode >= http.StatusInternalServerError {
		return c.retryOptions.RetryOnServerErrors
	}

	return false
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
	now := time.Now()
	if c.nextRequestAt.Before(now) {
		c.nextRequestAt = now
	}
	wait := c.nextRequestAt.Sub(now)
	c.nextRequestAt = c.nextRequestAt.Add(c.minRequestInterval)
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

func decodeGatewayErrorFromBytes(buf []byte, resp *http.Response) error {
	var body discord.GatewayError
	_ = json.Unmarshal(buf, &body)
	return &Error{
		HTTPStatus: resp.StatusCode,
		Code:       discord.GatewayErrorCode(body.Code),
		Message:    body.Message,
		Errors:     body.Errors,
	}
}

// snapshotResponse returns a shallow copy of resp with a fresh body reader
// backed by buf. The caller retains ownership of the original resp.
func snapshotResponse(resp *http.Response, buf []byte) *http.Response {
	snap := *resp
	snap.Body = io.NopCloser(bytes.NewReader(buf))
	return &snap
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
