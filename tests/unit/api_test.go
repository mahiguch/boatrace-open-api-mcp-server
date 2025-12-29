package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boatrace/open-api-mcp-server/internal/api"
)

// TestFetchWithCacheHit tests cache hit scenario
func TestFetchWithCacheHit(t *testing.T) {
	client := api.NewHTTPClient()
	cacheKey := "programs:20251222"
	expectedData := []byte(`{"date":"20251222","races":[]}`)

	// Pre-populate cache
	client.GetCache().Set(cacheKey, expectedData)

	// Fetch should return cached data without making HTTP request
	data, err := client.FetchWithCache("http://example.com/api", cacheKey)
	if err != nil {
		t.Fatalf("FetchWithCache failed: %v", err)
	}

	if string(data) != string(expectedData) {
		t.Errorf("expected cached data, got %s", string(data))
	}
}

// TestFetchWithRetrySuccess tests successful fetch with retry logic
func TestFetchWithRetrySuccess(t *testing.T) {
	// Create mock server that returns success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"date":"20251222","races":[]}`))
	}))
	defer server.Close()

	client := api.NewHTTPClient()
	data, err := client.FetchWithRetry(server.URL)

	if err != nil {
		t.Fatalf("FetchWithRetry failed: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("expected data, got empty")
	}
}

// TestFetchWithRetry404 tests 404 handling (no retry)
func TestFetchWithRetry404(t *testing.T) {
	// Create mock server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := api.NewHTTPClient()
	_, err := client.FetchWithRetry(server.URL)

	if err == nil {
		t.Errorf("expected error for 404, got nil")
	}
}

// TestFetchWithRetry500 tests 500 handling (with retry)
func TestFetchWithRetry500(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := api.NewHTTPClient()
	_, err := client.FetchWithRetry(server.URL)

	if err == nil {
		t.Errorf("expected error for 500, got nil")
	}

	// Should have been called 3 times (initial + 2 retries)
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

// TestUnmarshalJSONValid tests valid JSON unmarshalling
func TestUnmarshalJSONValid(t *testing.T) {
	jsonData := []byte(`{"test": "value"}`)
	var result map[string]string

	err := api.UnmarshalJSON(jsonData, &result)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if result["test"] != "value" {
		t.Errorf("expected 'value', got %s", result["test"])
	}
}

// TestUnmarshalJSONInvalid tests invalid JSON handling
func TestUnmarshalJSONInvalid(t *testing.T) {
	jsonData := []byte(`{invalid json}`)
	var result map[string]string

	err := api.UnmarshalJSON(jsonData, &result)
	if err == nil {
		t.Errorf("expected error for invalid JSON, got nil")
	}
}

// TestCacheKeyGeneration tests cache key format
func TestCacheKeyGeneration(t *testing.T) {
	client := api.NewHTTPClient()
	cache := client.GetCache()

	key := "programs:20251222"
	data := []byte("test data")

	cache.Set(key, data)

	retrieved, ok := cache.Get(key)
	if !ok {
		t.Errorf("cache miss for key %s", key)
	}

	if string(retrieved) != "test data" {
		t.Errorf("cache data mismatch")
	}
}
