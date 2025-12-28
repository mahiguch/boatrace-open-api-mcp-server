package api

import (
	"testing"
)

// TestNewHTTPClient tests client initialization
func TestNewHTTPClient(t *testing.T) {
	client := NewHTTPClient()
	if client == nil {
		t.Fatalf("NewHTTPClient() returned nil")
	}
	if client.cache == nil {
		t.Errorf("Expected cache to be initialized")
	}
}

// TestGetCache tests cache accessor
func TestGetCache(t *testing.T) {
	client := NewHTTPClient()
	cache := client.GetCache()
	if cache == nil {
		t.Errorf("GetCache() returned nil")
	}
}
