package cache

import (
	"testing"
)

// TestCachePackageInternal tests internal cache operations
func TestCachePackageInternal(t *testing.T) {
	cache := NewSessionCache()

	// Test basic operations
	key := "test:key"
	data := []byte("test data")

	cache.Set(key, data)

	retrieved, ok := cache.Get(key)
	if !ok {
		t.Errorf("cache.Get() returned false, expected true")
	}

	if string(retrieved) != string(data) {
		t.Errorf("cache data mismatch")
	}

	// Test size
	if cache.Size() != 1 {
		t.Errorf("expected size 1, got %d", cache.Size())
	}

	// Test clear
	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("expected size 0 after clear, got %d", cache.Size())
	}
}

// TestGenerateKeyInternal tests key generation
func TestGenerateKeyInternal(t *testing.T) {
	key := GenerateKey("programs", "20251222")
	expected := "programs:20251222"

	if key != expected {
		t.Errorf("GenerateKey mismatch: expected %s, got %s", expected, key)
	}
}
