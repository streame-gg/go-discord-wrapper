package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/options"
	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// newTestClient creates a RestClient pointed at ts with rate limiting disabled.
func newTestClient(ts *httptest.Server) *RestClient {
	rc, err := NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
		options.WithRetry(options.RetryOptions{
			MaxRetries:          0,
			RetryOnRateLimit:    false,
			RetryOnServerErrors: false,
		}),
	)
	if err != nil {
		panic(err)
	}
	return rc
}

func TestGetChannel(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantID     common.Snowflake
		wantName   string
		wantErr    error
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			body:       `{"id":"123456789","name":"general","type":0}`,
			wantID:     "123456789",
			wantName:   "general",
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			body:       `{"code":10003,"message":"Unknown Channel"}`,
			wantErr:    ErrNotFound,
		},
		{
			name:       "forbidden",
			statusCode: http.StatusForbidden,
			body:       `{"code":50013,"message":"Missing Permissions"}`,
			wantErr:    ErrForbidden,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			body:       `{"code":0,"message":"401: Unauthorized"}`,
			wantErr:    ErrUnauthorized,
		},
		{
			name:       "rate limited",
			statusCode: http.StatusTooManyRequests,
			body:       `{"code":0,"message":"You are being rate limited."}`,
			wantErr:    ErrRateLimited,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer ts.Close()

			client := newTestClient(ts)
			ch, err := client.GetChannel(context.Background(), "123456789")

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tc.wantErr), "expected errors.Is(%v) to be true, got: %v", tc.wantErr, err)
				assert.Nil(t, ch)
			} else {
				require.NoError(t, err)
				require.NotNil(t, ch)
				assert.Equal(t, tc.wantID, ch.ID)
				assert.Equal(t, tc.wantName, ch.Name)
			}
		})
	}
}

func TestCreateMessage(t *testing.T) {
	t.Run("success with correct body", func(t *testing.T) {
		var receivedBody []byte
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			var err error
			receivedBody = make([]byte, r.ContentLength)
			_, err = r.Body.Read(receivedBody)
			_ = err

			resp := common.Message{
				ID:        "987654321",
				ChannelID: "111222333",
				Content:   "hello world",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(resp)
		}))
		defer ts.Close()

		client := newTestClient(ts)
		params := CreateMessageParams{Content: "hello world"}
		msg, err := client.CreateMessage(context.Background(), "111222333", params)

		require.NoError(t, err)
		require.NotNil(t, msg)
		assert.Equal(t, common.Snowflake("987654321"), msg.ID)
		assert.Equal(t, "hello world", msg.Content)

		var decoded map[string]interface{}
		require.NoError(t, json.Unmarshal(receivedBody, &decoded))
		assert.Equal(t, "hello world", decoded["content"])
	})

	t.Run("non-2xx returns APIError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"code":50013,"message":"Missing Permissions"}`))
		}))
		defer ts.Close()

		client := newTestClient(ts)
		msg, err := client.CreateMessage(context.Background(), "111222333", CreateMessageParams{Content: "hello"})

		require.Error(t, err)
		assert.Nil(t, msg)

		var apiErr *Error
		assert.True(t, errors.As(err, &apiErr))
		assert.Equal(t, http.StatusForbidden, apiErr.HTTPStatus)
	})
}

func TestGetGuild(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := common.Guild{
			ID:   "777888999",
			Name: "Test Server",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := newTestClient(ts)
	guild, err := client.GetGuild(context.Background(), "777888999", false)

	require.NoError(t, err)
	require.NotNil(t, guild)
	assert.Equal(t, common.Snowflake("777888999"), guild.ID)
	assert.Equal(t, "Test Server", guild.Name)
}

func TestContextCancellation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"1","name":"general","type":0}`))
	}))
	defer ts.Close()

	client := newTestClient(ts)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.GetChannel(ctx, "1")
	require.Error(t, err)
	assert.True(t,
		errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded),
		"expected context cancellation error, got: %v", err,
	)
}

func TestRetryOn429(t *testing.T) {
	var callCount int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			w.Header().Set("Retry-After", "0")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"code":0,"message":"rate limited"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"1","name":"general","type":0}`))
	}))
	defer ts.Close()

	client, err := NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
		options.WithRetry(options.RetryOptions{
			MaxRetries:          2,
			BaseBackoff:         0,
			RetryOnRateLimit:    true,
			RetryOnServerErrors: false,
		}),
	)
	require.NoError(t, err)

	ch, err := client.GetChannel(context.Background(), "1")
	require.NoError(t, err)
	require.NotNil(t, ch)
	assert.Equal(t, int32(2), atomic.LoadInt32(&callCount), "expected exactly 2 requests")
}

func TestNoRetryWhenDisabled(t *testing.T) {
	var callCount int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":0,"message":"internal server error"}`))
	}))
	defer ts.Close()

	client, err := NewRestClient("test-token",
		options.WithBaseURL(ts.URL),
		options.WithRateLimiting(options.RateLimiterOptions{Disabled: true}),
		options.WithRetry(options.RetryOptions{
			MaxRetries:          0,
			RetryOnRateLimit:    false,
			RetryOnServerErrors: false,
		}),
	)
	require.NoError(t, err)

	_, err = client.GetChannel(context.Background(), "1")
	require.Error(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&callCount), "expected exactly 1 request")
}

func TestDecodeAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":10003,"message":"Unknown Channel"}`))
	}))
	defer ts.Close()

	client := newTestClient(ts)
	_, err := client.GetChannel(context.Background(), "1")

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound), "expected ErrNotFound")

	var apiErr *Error
	require.True(t, errors.As(err, &apiErr), "expected *APIError")
	assert.Equal(t, 10003, int(apiErr.Code))
	assert.Equal(t, "Unknown Channel", apiErr.Message)
	assert.Equal(t, http.StatusNotFound, apiErr.HTTPStatus)
}
