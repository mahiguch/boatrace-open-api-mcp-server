package tools

import (
	"encoding/json"

	"github.com/boatrace/open-api-mcp-server/internal"
)

// ProgramsHandler handles the programs MCP tool
type ProgramsHandler struct {
	*ToolHandler
}

// NewProgramsHandler creates a new programs tool handler
func NewProgramsHandler(handler *ToolHandler) *ProgramsHandler {
	return &ProgramsHandler{ToolHandler: handler}
}

// Response represents the response structure for tools
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Execute handles the programs tool execution
func (h *ProgramsHandler) Execute(year, date string) *Response {
	// Log tool invocation
	h.logger.Info("tool_invoked", map[string]interface{}{
		"tool": "programs",
		"year": year,
		"date": date,
	})

	// Fetch programs data
	program, err := h.client.FetchPrograms(year, date)
	if err != nil {
		h.logger.Error("tool_execution_failed", err, map[string]interface{}{
			"tool": "programs",
		})
		return errorResponse(err)
	}

	// Log successful fetch
	h.logger.Info("api_call_success", map[string]interface{}{
		"tool": "programs",
		"date": date,
	})

	return &Response{
		Success: true,
		Data:    program,
	}
}

// results tool handler

// ResultsHandler handles the results MCP tool
type ResultsHandler struct {
	*ToolHandler
}

// NewResultsHandler creates a new results tool handler
func NewResultsHandler(handler *ToolHandler) *ResultsHandler {
	return &ResultsHandler{ToolHandler: handler}
}

// Execute handles the results tool execution
func (h *ResultsHandler) Execute(year, date string) *Response {
	// Log tool invocation
	h.logger.Info("tool_invoked", map[string]interface{}{
		"tool": "results",
		"year": year,
		"date": date,
	})

	// Fetch results data
	result, err := h.client.FetchResults(year, date)
	if err != nil {
		h.logger.Error("tool_execution_failed", err, map[string]interface{}{
			"tool": "results",
		})
		return errorResponse(err)
	}

	// Log successful fetch
	h.logger.Info("api_call_success", map[string]interface{}{
		"tool": "results",
		"date": date,
	})

	return &Response{
		Success: true,
		Data:    result,
	}
}

// previews tool handler

// PreviewsHandler handles the previews MCP tool
type PreviewsHandler struct {
	*ToolHandler
}

// NewPreviewsHandler creates a new previews tool handler
func NewPreviewsHandler(handler *ToolHandler) *PreviewsHandler {
	return &PreviewsHandler{ToolHandler: handler}
}

// Execute handles the previews tool execution
func (h *PreviewsHandler) Execute(year, date string) *Response {
	// Log tool invocation
	h.logger.Info("tool_invoked", map[string]interface{}{
		"tool": "previews",
		"year": year,
		"date": date,
	})

	// Fetch previews data
	preview, err := h.client.FetchPreviews(year, date)
	if err != nil {
		h.logger.Error("tool_execution_failed", err, map[string]interface{}{
			"tool": "previews",
		})
		return errorResponse(err)
	}

	// Log successful fetch
	h.logger.Info("api_call_success", map[string]interface{}{
		"tool": "previews",
		"date": date,
	})

	return &Response{
		Success: true,
		Data:    preview,
	}
}

// errorResponse converts an error to a response
func errorResponse(err error) *Response {
	if appErr, ok := err.(*internal.AppError); ok {
		return &Response{
			Success: false,
			Error: &ErrorInfo{
				Type:    string(appErr.Type),
				Message: appErr.Message,
				Code:    appErr.Code,
			},
		}
	}

	return &Response{
		Success: false,
		Error: &ErrorInfo{
			Type:    "UnknownError",
			Message: err.Error(),
			Code:    "UNKNOWN",
		},
	}
}

// MarshalResponse converts a response to JSON
func MarshalResponse(resp *Response) ([]byte, error) {
	return json.MarshalIndent(resp, "", "  ")
}
