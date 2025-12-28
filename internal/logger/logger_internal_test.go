package logger

import (
	"testing"
)

// TestLoggerPackageInternal tests logger functionality
func TestLoggerPackageInternal(t *testing.T) {
	logger := New()

	if logger == nil {
		t.Fatalf("New() returned nil")
	}

	// Test Info logging
	logger.Info("test message", map[string]interface{}{
		"key": "value",
		"num": 42,
	})

	// Test Error logging
	testErr := &struct {
		error
		msg string
	}{
		msg: "test error",
	}
	_ = testErr

	// Create a simple error for testing
	var errVal error = nil
	if errVal == nil {
		errVal = &testError{"test error"}
	}

	logger.Error("error message", errVal, map[string]interface{}{
		"context": "test",
	})

	// No panic means success
	t.Log("Logger functionality verified")
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
