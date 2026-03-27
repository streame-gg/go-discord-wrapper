package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/streame-gg/go-discord-wrapper/types/common"
	"github.com/streame-gg/go-discord-wrapper/util"
)

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

type RetryOptions struct {
	MaxRetries          int
	BaseBackoff         time.Duration
	MaxBackoff          time.Duration
	RetryOnRateLimit    bool
	RetryOnServerErrors bool
}

func defaultRetryOptions() RetryOptions {
	return RetryOptions{
		MaxRetries:          3,
		BaseBackoff:         500 * time.Millisecond,
		MaxBackoff:          5 * time.Second,
		RetryOnRateLimit:    true,
		RetryOnServerErrors: true,
	}
}

type RestClient struct {
	BaseURL string
	token   string
	Version common.APIVersion

	httpClient *http.Client

	retryOptions RetryOptions

	minRequestInterval time.Duration
	rateLimitMu        sync.Mutex
	nextRequestAt      time.Time

	eventEmitter *util.EventEmitter[RestEventType, RestEventHandler]
}

type RestClientOption func(*RestClient)

func WithBaseURL(baseURL string) RestClientOption {
	return func(c *RestClient) {
		c.BaseURL = baseURL
	}
}

func WithApiVersion(version common.APIVersion) RestClientOption {
	return func(c *RestClient) {
		c.Version = version
	}
}

func WithHttpClient(client *http.Client) RestClientOption {
	return func(c *RestClient) {
		c.httpClient = client
	}
}

func WithRetryOptions(options RetryOptions) RestClientOption {
	return func(c *RestClient) {
		c.retryOptions = options
	}
}

func WithMinRequestInterval(interval time.Duration) RestClientOption {
	return func(c *RestClient) {
		c.minRequestInterval = interval
	}
}

func NewRestClient(token string, options ...RestClientOption) *RestClient {
	c := &RestClient{
		BaseURL:      "https://discord.com/api",
		token:        token,
		Version:      common.APIVersion10,
		httpClient:   http.DefaultClient,
		retryOptions: defaultRetryOptions(),
		eventEmitter: util.NewEventEmitter[RestEventType, RestEventHandler](),
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *RestClient) buildURL() string {
	return c.BaseURL + "/v" + c.Version.ToString()
}

func (c *RestClient) OnEvent(eventType RestEventType, handler RestEventHandler) {
	if c.eventEmitter == nil {
		c.eventEmitter = util.NewEventEmitter[RestEventType, RestEventHandler]()
	}

	c.eventEmitter.On(eventType, handler)
}

func (c *RestClient) emitEvent(event RestEvent) {
	if c.eventEmitter == nil {
		return
	}

	for _, handler := range c.eventEmitter.Handlers(event.Type) {
		handler(c, event)
	}
}

func (c *RestClient) generateRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.buildURL()+path, body)
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

	bodyBytes, err := captureRequestBody(req)
	if err != nil {
		return nil, err
	}

	maxAttempts := c.retryOptions.MaxRetries + 1

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := c.waitForClientRateLimit(); err != nil {
			return nil, err
		}

		attemptReq, err := cloneRequest(req, bodyBytes)
		if err != nil {
			return nil, err
		}

		c.emitEvent(RestEvent{Type: RestEventRequest, Request: attemptReq, Attempt: attempt})

		resp, reqErr := c.httpClient.Do(attemptReq)
		if reqErr != nil {
			c.emitEvent(RestEvent{Type: RestEventError, Request: attemptReq, Attempt: attempt, Err: reqErr})
			if attempt == maxAttempts {
				return nil, reqErr
			}

			delay := c.retryDelay(attempt, 0)
			c.emitEvent(RestEvent{Type: RestEventRetry, Request: attemptReq, Attempt: attempt, RetryAfter: delay, Err: reqErr})
			time.Sleep(delay)
			continue
		}

		c.emitEvent(RestEvent{Type: RestEventResponse, Request: attemptReq, Response: resp, Attempt: attempt})

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
				c.emitEvent(RestEvent{Type: RestEventRateLimit, Request: attemptReq, Response: resp, Attempt: attempt, RetryAfter: retryAfter})
			}

			c.emitEvent(RestEvent{Type: RestEventRetry, Request: attemptReq, Response: resp, Attempt: attempt, RetryAfter: delay})
			_ = resp.Body.Close()
			time.Sleep(delay)
			continue
		}

		respErr := decodeGatewayError(resp)
		_ = resp.Body.Close()
		c.emitEvent(RestEvent{Type: RestEventError, Request: attemptReq, Response: resp, Attempt: attempt, Err: respErr})
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

func (c *RestClient) waitForClientRateLimit() error {
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
		time.Sleep(wait)
	}

	return nil
}

func parseRetryAfter(resp *http.Response) time.Duration {
	for _, header := range []string{"Retry-After", "X-RateLimit-Reset-After"} {
		value := strings.TrimSpace(resp.Header.Get(header))
		if value == "" {
			continue
		}

		seconds, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}

		if seconds <= 0 {
			continue
		}

		return time.Duration(seconds * float64(time.Second))
	}

	return 0
}

func decodeGatewayError(resp *http.Response) error {
	var respErr common.GatewayError
	if err := json.NewDecoder(resp.Body).Decode(&respErr); err == nil && (respErr.Message != "" || respErr.Code != 0) {
		return respErr
	}

	return fmt.Errorf("request failed with status %s", resp.Status)
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
