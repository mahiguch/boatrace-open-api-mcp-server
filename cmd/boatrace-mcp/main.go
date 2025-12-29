package main

import (
	"fmt"
	"os"

	"github.com/boatrace/open-api-mcp-server/internal/api"
	"github.com/boatrace/open-api-mcp-server/internal/logger"
	"github.com/boatrace/open-api-mcp-server/internal/tools"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Initialize HTTP client with caching
	httpClient := api.NewHTTPClient()

	// Initialize tool handler
	toolHandler := tools.NewToolHandler(httpClient, log)

	// Create tool-specific handlers
	programsHandler := tools.NewProgramsHandler(toolHandler)
	resultsHandler := tools.NewResultsHandler(toolHandler)
	previewsHandler := tools.NewPreviewsHandler(toolHandler)

	// Log startup
	log.Info("boatrace_mcp_server_starting", map[string]interface{}{
		"version": "1.0.0",
	})

	// Example tool invocation (for testing)
	// This would be called by the MCP SDK in a real implementation
	if len(os.Args) > 1 {
		if os.Args[1] == "test-programs" {
			resp := programsHandler.Execute("2025", "20251222")
			data, _ := tools.MarshalResponse(resp)
			fmt.Println(string(data))
			return
		}
		if os.Args[1] == "test-results" {
			resp := resultsHandler.Execute("2025", "20251222")
			data, _ := tools.MarshalResponse(resp)
			fmt.Println(string(data))
			return
		}
		if os.Args[1] == "test-previews" {
			resp := previewsHandler.Execute("2025", "20251222")
			data, _ := tools.MarshalResponse(resp)
			fmt.Println(string(data))
			return
		}
	}

	// In a real MCP server, this would initialize the MCP SDK
	// and register these tool handlers for incoming MCP calls
	fmt.Printf("Server initialized with %d tools\n", 3)
	fmt.Println("Tools: programs, results, previews")

	// Keep the server running
	select {}
}
