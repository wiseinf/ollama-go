package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client represents an Ollama API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	opts       *ClientOptions
	logger     Logger
	limiter    RateLimiter
}

// ClientOption is a function that modifies the client
type ClientOption func(*ClientOptions)

// NewClient creates a new Ollama client
func NewClient(options ...ClientOption) *Client {
	opts := defaultOptions()
	for _, opt := range options {
		opt(opts)
	}
	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: opts.Timeout,
		}
	}
	return &Client{
		baseURL:    opts.BaseURL,
		opts:       opts,
		httpClient: httpClient,
		logger:     opts.Logger,
		limiter:    newRateLimiter(opts.RateLimit),
	}
}

// sendRequest is a helper function to send requests with retries and rate limiting
func (c *Client) sendRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	// Apply rate limiting
	if err := c.limiter.Wait(); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	var resp *http.Response
	var err error

	// Retry logic
	for attempt := 0; attempt <= c.opts.MaxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Debug("Retrying request (attempt %d/%d)", attempt, c.opts.MaxRetries)
			// Calculate backoff time
			waitTime := c.opts.RetryWaitTime * time.Duration(1<<uint(attempt-1))
			if waitTime > c.opts.RetryMaxWaitTime {
				waitTime = c.opts.RetryMaxWaitTime
			}
			time.Sleep(waitTime)
		}

		var buf bytes.Buffer
		if body != nil {
			if err := json.NewEncoder(&buf).Encode(body); err != nil {
				return nil, fmt.Errorf("failed to encode request body: %w", err)
			}
		}
		var req *http.Request
		req, err = http.NewRequestWithContext(ctx, method, c.baseURL+path, &buf)
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		c.logger.Debug("Sending request: %s %s", method, path)
		resp, err = c.httpClient.Do(req)
		if err != nil {
			c.logger.Error("Request failed: %v", err)
			continue
		}
		c.logger.Debug("Receiving response: %s %s", method, path)
		if resp.StatusCode == http.StatusTooManyRequests ||
			(resp.StatusCode >= 500 && resp.StatusCode < 600) {
			resp.Body.Close()
			continue
		}

		break
	}

	if err != nil {
		return nil, fmt.Errorf("all retries failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("http status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("api error: %s", errResp.Error)
	}

	return resp, nil
}
