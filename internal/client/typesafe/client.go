package typesafe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const (
	defaultBaseURL     = "https://api.typesafe.ai"
	systemOneEndpoint  = "/v1/systemone"
	defaultHTTPTimeout = 30 * time.Second
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
}

// APIError represents an error response from the TypeSafe API
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("typesafe api error: status %d, body: %s", e.StatusCode, e.Body)
}

// NewClient creates a new TypeSafe API client
func NewClient(apiKey string, logger *slog.Logger) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
		logger: logger.With("component", "typesafe_client"),
	}
}

// Option is a functional option for Client
type Option func(*Client)

// WithBaseURL overrides the default base URL
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient overrides the HTTP client (useful for testing)
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// NewClientWithOptions creates a new TypeSafe API client with custom options
func NewClientWithOptions(apiKey string, logger *slog.Logger, opts ...Option) *Client {
	c := NewClient(apiKey, logger)
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Call invokes the TypeSafe SystemOne API with exponential backoff retry
func (c *Client) Call(ctx context.Context, req *Request) (*Response, error) {
	// Create backoff strategy
	strategy := backoff.NewExponentialBackOff()
	strategy.InitialInterval = 500 * time.Millisecond
	strategy.MaxInterval = 30 * time.Second
	strategy.MaxElapsedTime = 2 * time.Minute
	strategy.Multiplier = 2.0
	strategy.RandomizationFactor = 0.5

	// Wrap with context
	b := backoff.WithContext(strategy, ctx)

	var attempt int
	var result *Response

	// Execute with retries
	err := backoff.RetryNotify(func() error {
		attempt++
		resp, callErr := c.callOnce(ctx, req, attempt)
		if callErr == nil {
			result = resp
		}
		return callErr
	}, b, func(err error, duration time.Duration) {
		if err == nil {
			return
		}

		// Log retry attempt
		c.logger.Warn("api call failed, retrying",
			"attempt", attempt,
			"error", err.Error(),
			"retry_after", duration.String(),
		)
	})

	if err != nil {
		c.logger.Error("api call failed after all retries",
			"attempt", attempt,
			"error", err.Error(),
		)
		return nil, err
	}

	return result, nil
}

// callOnce executes a single API call without retries
func (c *Client) callOnce(ctx context.Context, req *Request, attempt int) (*Response, error) {
	// Marshal request
	body, err := json.Marshal(req)
	if err != nil {
		return nil, backoff.Permanent(fmt.Errorf("marshaling request: %w", err))
	}

	// Create request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+systemOneEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, backoff.Permanent(fmt.Errorf("creating request: %w", err))
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		// Network errors are retryable
		c.logger.Debug("network error on api call", "attempt", attempt, "error", err.Error())
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Log the response
	c.logger.Debug("api call completed", "attempt", attempt, "status", resp.StatusCode)

	// Check status code
	if resp.StatusCode >= 500 {
		// 5xx errors are retryable
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		// 429 is retryable
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	if resp.StatusCode >= 400 {
		// Other 4xx errors are not retryable
		return nil, backoff.Permanent(&APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		})
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Unexpected status codes
		return nil, backoff.Permanent(&APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		})
	}

	// 2xx success - decode response
	var result Response
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, backoff.Permanent(fmt.Errorf("decoding response: %w", err))
	}

	c.logger.Info("api call succeeded", "attempt", attempt)
	return &result, nil
}
