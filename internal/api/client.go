package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/boatrace/open-api-mcp-server/internal"
	"github.com/boatrace/open-api-mcp-server/internal/cache"
)

const (
	baseURL         = "https://boatraceopenapi.github.io"
	httpTimeout     = 10 * time.Second
	maxRetries      = 2
	initialBackoff  = 100 * time.Millisecond
	retryBackoff    = 200 * time.Millisecond
)

// HTTPClient wraps the standard HTTP client with retry logic, caching, and timeouts
type HTTPClient struct {
	client *http.Client
	cache  *cache.SessionCache
}

// NewHTTPClient creates a new HTTP client with timeout and caching
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: httpTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
		cache: cache.NewSessionCache(),
	}
}

// GetCache returns the session cache
func (c *HTTPClient) GetCache() *cache.SessionCache {
	return c.cache
}

// FetchWithCache fetches data from the API with caching and retry logic
func (c *HTTPClient) FetchWithCache(url, cacheKey string) ([]byte, error) {
	// Check cache first
	if cached, ok := c.cache.Get(cacheKey); ok {
		return cached, nil
	}

	// Fetch from API with retry logic
	data, err := c.FetchWithRetry(url)
	if err != nil {
		return nil, err
	}

	// Store in cache
	c.cache.Set(cacheKey, data)

	return data, nil
}

// FetchWithRetry fetches data from the API with automatic retry on transient failures
func (c *HTTPClient) FetchWithRetry(url string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		data, err := c.fetch(url)
		if err == nil {
			return data, nil
		}

		lastErr = err

		// Don't retry on validation or parse errors
		if _, ok := err.(*internal.AppError); ok {
			if appErr, ok := err.(*internal.AppError); ok {
				if appErr.Type == internal.ValidationError || appErr.Type == internal.ParseError {
					return nil, err
				}
			}
		}

		// If not the last attempt, wait before retrying
		if attempt < maxRetries {
			backoff := initialBackoff + time.Duration(attempt)*retryBackoff
			time.Sleep(backoff)
		}
	}

	return nil, lastErr
}

// fetch performs a single HTTP GET request
func (c *HTTPClient) fetch(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, internal.NewAPIError(0, fmt.Sprintf("failed to create request: %v", err))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if err == http.ErrHandlerTimeout {
			return nil, internal.NewTimeoutError(fmt.Sprintf("request exceeded %v timeout", httpTimeout))
		}
		return nil, internal.NewAPIError(0, fmt.Sprintf("request failed: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, internal.NewAPIError(resp.StatusCode, fmt.Sprintf("failed to read response: %v", err))
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, internal.NewNotFoundError("no races found for this date")
		}
		if resp.StatusCode >= 500 {
			return nil, internal.NewServerError(fmt.Sprintf("server returned status %d", resp.StatusCode))
		}
		return nil, internal.NewAPIError(resp.StatusCode, fmt.Sprintf("HTTP %d", resp.StatusCode))
	}

	return body, nil
}

// UnmarshalJSON is a helper to unmarshal JSON with error handling
func UnmarshalJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return internal.NewParseError(fmt.Sprintf("failed to parse JSON: %v", err))
	}
	return nil
}
