package unit

import (
	"testing"

	"github.com/boatrace/open-api-mcp-server/internal/cache"
)

// TestNewSessionCache tests cache initialization
func TestNewSessionCache(t *testing.T) {
	c := cache.NewSessionCache()
	if c == nil {
		t.Errorf("NewSessionCache() returned nil")
	}
	if c.Size() != 0 {
		t.Errorf("NewSessionCache() expected size 0, got %d", c.Size())
	}
}

// TestCacheSetAndGet tests cache set and get operations
func TestCacheSetAndGet(t *testing.T) {
	c := cache.NewSessionCache()
	key := "programs:20251222"
	data := []byte(`{"date":"20251222","races":[]}`)

	// Set value
	c.Set(key, data)

	// Get value
	cached, ok := c.Get(key)
	if !ok {
		t.Errorf("cache.Get() returned false, expected true")
	}
	if string(cached) != string(data) {
		t.Errorf("cache.Get() returned %s, expected %s", string(cached), string(data))
	}
}

// TestCacheMiss tests cache miss scenario
func TestCacheMiss(t *testing.T) {
	c := cache.NewSessionCache()
	key := "programs:20251222"

	cached, ok := c.Get(key)
	if ok {
		t.Errorf("cache.Get() returned true, expected false for missing key")
	}
	if cached != nil {
		t.Errorf("cache.Get() returned non-nil data, expected nil")
	}
}

// TestGenerateKey tests cache key generation
func TestGenerateKey(t *testing.T) {
	tests := []struct {
		tool     string
		date     string
		expected string
	}{
		{
			tool:     "programs",
			date:     "20251222",
			expected: "programs:20251222",
		},
		{
			tool:     "results",
			date:     "20251220",
			expected: "results:20251220",
		},
		{
			tool:     "previews",
			date:     "20251222",
			expected: "previews:20251222",
		},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			key := cache.GenerateKey(tt.tool, tt.date)
			if key != tt.expected {
				t.Errorf("GenerateKey(%s, %s) returned %s, expected %s", tt.tool, tt.date, key, tt.expected)
			}
		})
	}
}

// TestCacheSize tests cache size tracking
func TestCacheSize(t *testing.T) {
	c := cache.NewSessionCache()

	if c.Size() != 0 {
		t.Errorf("expected size 0, got %d", c.Size())
	}

	c.Set("key1", []byte("data1"))
	if c.Size() != 1 {
		t.Errorf("expected size 1, got %d", c.Size())
	}

	c.Set("key2", []byte("data2"))
	if c.Size() != 2 {
		t.Errorf("expected size 2, got %d", c.Size())
	}

	c.Set("key1", []byte("updated_data1"))
	if c.Size() != 2 {
		t.Errorf("expected size 2 (update, not new), got %d", c.Size())
	}
}

// TestCacheClear tests cache clearing
func TestCacheClear(t *testing.T) {
	c := cache.NewSessionCache()

	c.Set("key1", []byte("data1"))
	c.Set("key2", []byte("data2"))

	if c.Size() != 2 {
		t.Errorf("expected size 2, got %d", c.Size())
	}

	c.Clear()

	if c.Size() != 0 {
		t.Errorf("expected size 0 after clear, got %d", c.Size())
	}

	_, ok := c.Get("key1")
	if ok {
		t.Errorf("cache should be empty after clear, but found key1")
	}
}

// TestConcurrentAccess tests concurrent cache access (basic test)
func TestConcurrentAccess(t *testing.T) {
	c := cache.NewSessionCache()

	// Simple concurrent test
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			c.Set("concurrent", []byte("data"))
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			c.Get("concurrent")
		}
		done <- true
	}()

	<-done
	<-done

	// Should not panic or error
	if c.Size() > 0 {
		t.Logf("concurrent access completed, cache size: %d", c.Size())
	}
}
