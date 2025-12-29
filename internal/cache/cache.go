package cache

import (
	"fmt"
	"sync"
)

// CachedResponse represents a cached API response
type CachedResponse struct {
	Data []byte
}

// SessionCache represents an in-memory cache for the current MCP session
type SessionCache struct {
	mu    sync.RWMutex
	store map[string]CachedResponse
}

// NewSessionCache creates a new session cache
func NewSessionCache() *SessionCache {
	return &SessionCache{
		store: make(map[string]CachedResponse),
	}
}

// Get retrieves a cached response by key
// Returns the cached data and true if found, nil and false if not found
func (c *SessionCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.store[key]
	if !ok {
		return nil, false
	}

	return cached.Data, true
}

// Set stores a response in the cache
func (c *SessionCache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = CachedResponse{Data: data}
}

// GenerateKey generates a cache key from tool name and date
// Format: {tool}:{yyyymmdd}
func GenerateKey(tool, date string) string {
	return fmt.Sprintf("%s:%s", tool, date)
}

// Clear removes all cached entries from the cache
func (c *SessionCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store = make(map[string]CachedResponse)
}

// Size returns the number of cached entries
func (c *SessionCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.store)
}
