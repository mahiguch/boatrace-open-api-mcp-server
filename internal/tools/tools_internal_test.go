package tools

import (
	"testing"

	"github.com/boatrace/open-api-mcp-server/internal/api"
	"github.com/boatrace/open-api-mcp-server/internal/logger"
)

// TestToolHandlerPackageInternal tests tool handler initialization
func TestToolHandlerPackageInternal(t *testing.T) {
	client := api.NewHTTPClient()
	log := logger.New()

	handler := NewToolHandler(client, log)

	if handler == nil {
		t.Fatalf("NewToolHandler() returned nil")
	}

	if handler.client == nil {
		t.Errorf("ToolHandler client is nil")
	}

	if handler.logger == nil {
		t.Errorf("ToolHandler logger is nil")
	}
}

// TestProgramsHandlerInternal tests programs handler
func TestProgramsHandlerInternal(t *testing.T) {
	client := api.NewHTTPClient()
	log := logger.New()
	handler := NewToolHandler(client, log)

	programsHandler := NewProgramsHandler(handler)

	if programsHandler == nil {
		t.Fatalf("NewProgramsHandler() returned nil")
	}

	// Test Execute with invalid date
	resp := programsHandler.Execute("25", "20251222")

	if resp == nil {
		t.Fatalf("Execute() returned nil")
	}

	if resp.Success {
		t.Errorf("Expected error for invalid year, got success")
	}

	if resp.Error == nil {
		t.Errorf("Expected error info, got nil")
	}
}

// TestResultsHandlerInternal tests results handler
func TestResultsHandlerInternal(t *testing.T) {
	client := api.NewHTTPClient()
	log := logger.New()
	handler := NewToolHandler(client, log)

	resultsHandler := NewResultsHandler(handler)

	if resultsHandler == nil {
		t.Fatalf("NewResultsHandler() returned nil")
	}

	// Test Execute with invalid date
	resp := resultsHandler.Execute("2025", "20251301")

	if resp == nil {
		t.Fatalf("Execute() returned nil")
	}

	if resp.Success {
		t.Errorf("Expected error for invalid date, got success")
	}
}

// TestPreviewsHandlerInternal tests previews handler
func TestPreviewsHandlerInternal(t *testing.T) {
	client := api.NewHTTPClient()
	log := logger.New()
	handler := NewToolHandler(client, log)

	previewsHandler := NewPreviewsHandler(handler)

	if previewsHandler == nil {
		t.Fatalf("NewPreviewsHandler() returned nil")
	}

	// Test Execute with mismatched year
	resp := previewsHandler.Execute("2024", "20251222")

	if resp == nil {
		t.Fatalf("Execute() returned nil")
	}

	if resp.Success {
		t.Errorf("Expected error for mismatched year, got success")
	}
}

// TestResponseStructure tests Response structure
func TestResponseStructure(t *testing.T) {
	// Test success response
	resp := &Response{
		Success: true,
		Data:    map[string]string{"key": "value"},
	}

	if !resp.Success {
		t.Errorf("Response Success mismatch")
	}

	// Test error response
	errResp := &Response{
		Success: false,
		Error: &ErrorInfo{
			Type:    "ValidationError",
			Message: "Invalid input",
			Code:    "INVALID_INPUT",
		},
	}

	if errResp.Success {
		t.Errorf("Error response should have Success=false")
	}

	if errResp.Error == nil {
		t.Errorf("Error response should have Error info")
	}
}

// TestMarshalResponse tests response marshalling
func TestMarshalResponse(t *testing.T) {
	resp := &Response{
		Success: true,
		Data:    map[string]string{"test": "data"},
	}

	data, err := MarshalResponse(resp)

	if err != nil {
		t.Errorf("MarshalResponse() error: %v", err)
	}

	if len(data) == 0 {
		t.Errorf("MarshalResponse() returned empty data")
	}

	if string(data[0:1]) != "{" {
		t.Errorf("MarshalResponse() should return JSON")
	}
}
