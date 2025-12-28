package tools

import (
	"github.com/boatrace/open-api-mcp-server/internal/api"
	"github.com/boatrace/open-api-mcp-server/internal/logger"
)

// ToolHandler provides MCP tool implementations
type ToolHandler struct {
	client *api.HTTPClient
	logger *logger.Logger
}

// NewToolHandler creates a new tool handler
func NewToolHandler(client *api.HTTPClient, logger *logger.Logger) *ToolHandler {
	return &ToolHandler{
		client: client,
		logger: logger,
	}
}

// GetClient returns the HTTP client
func (h *ToolHandler) GetClient() *api.HTTPClient {
	return h.client
}

// GetLogger returns the logger
func (h *ToolHandler) GetLogger() *logger.Logger {
	return h.logger
}
