# Boatrace Open API MCP Server

A Model Context Protocol (MCP) server providing access to official Boatrace race data through three core tools.

## Overview

This MCP server enables Claude and other AI models to access Boatrace race information through the official Boatrace Open API. Three tools provide access to:

- **programs**: Race program data (出走表) - racer lineups, race conditions, odds
- **results**: Race results (結果) - winners, payouts, race statistics
- **previews**: Race previews (直前情報) - updated odds, predictions, analysis

## Prerequisites

- Go 1.21 or later
- Access to Boatrace Open API (https://boatraceopenapi.github.io/)

## Installation

### From Source

```bash
git clone https://github.com/boatrace/open-api-mcp-server.git
cd boatrace-open-api-mcp-server
go mod download
make build
```

The binary will be created at `bin/boatrace-mcp`.

## Usage

### Starting the Server

```bash
./bin/boatrace-mcp
```

The server will start listening for MCP calls.

### Using the Tools

#### programs - Fetch Race Programs

```
programs(year="2025", date="20251222")
```

Returns complete race program for the specified date including:
- Race lineup with racer numbers and names
- Race conditions (distance, weather)
- Betting odds and favorites
- Race metadata

#### results - Fetch Race Results

```
results(year="2025", date="20251220")
```

Returns race results for completed races including:
- Winning racer information
- Payouts (win, place, show)
- Race statistics and timing
- Payout details

#### previews - Fetch Race Previews

```
previews(year="2025", date="20251222")
```

Returns pre-race analysis and information including:
- Updated odds
- Expert predictions
- Race analysis
- Late-breaking news

## IDE & CLI Integration

### Claude Code

To use this MCP server with Claude Code:

1. **Build the server**:
   ```bash
   make build
   ```

2. **Configure Claude Code** - Add to your `claude_code.json` or MCP configuration:
   ```json
   {
     "mcpServers": {
       "boatrace": {
         "command": "/path/to/boatrace-open-api-mcp-server/bin/boatrace-mcp",
         "args": []
       }
     }
   }
   ```

3. **Use in Claude Code** - The tools will be automatically available:
   - `programs(year="2025", date="20251222")` - Get race programs
   - `results(year="2025", date="20251220")` - Get race results
   - `previews(year="2025", date="20251222")` - Get race previews

### GitHub Copilot

To use this MCP server with GitHub Copilot (when MCP support is available):

1. **Build and locate the binary**:
   ```bash
   make build
   # Binary location: ./bin/boatrace-mcp
   ```

2. **Configure in `.copilot/config.json`** (or your Copilot configuration):
   ```json
   {
     "servers": [
       {
         "name": "boatrace",
         "command": "./bin/boatrace-mcp"
       }
     ]
   }
   ```

3. **Usage in Copilot chat** - Ask Copilot questions about Boatrace data:
   - "What are the race programs for 2025-12-22?"
   - "Show me the results for December 22, 2025"
   - "Get previews for today's races"

### Gemini CLI

To use this MCP server with Gemini CLI (via gemini-cli):

1. **Build the server**:
   ```bash
   make build
   ```

2. **Start the MCP server in the background**:
   ```bash
   ./bin/boatrace-mcp &
   ```

3. **Configure Gemini CLI** - Set environment variable or config:
   ```bash
   export MCP_SERVER_BOATRACE="localhost:9000"
   # or add to your gemini-cli config
   ```

4. **Use with Gemini CLI**:
   ```bash
   gemini "Get the race programs for 2025-12-22"
   # The boatrace tools will be available to Gemini
   ```

Example conversation with Gemini CLI:
```
$ gemini chat
> Get the race programs for 2025-12-22
[Gemini calls the programs tool via MCP]
> What were the results from that day?
[Gemini calls the results tool]
> Show me previews for tomorrow's races
[Gemini calls the previews tool]
```

## Development

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

### Cleaning

```bash
make clean
```

## Project Structure

```
.
├── cmd/boatrace-mcp/          # CLI application entry point
├── internal/
│   ├── models/                # Data models (RaceProgram, RaceResult, RacePreview)
│   ├── api/                   # HTTP client and API endpoints
│   ├── validation/            # Input validation
│   ├── cache/                 # Session-based caching
│   ├── tools/                 # MCP tool handlers
│   └── logger/                # Logging utilities
├── tests/
│   ├── unit/                  # Unit tests
│   ├── integration/           # Integration tests
│   └── fixtures/              # Test fixtures
├── Makefile                   # Build targets
├── .golangci.yml              # Linter configuration
└── go.mod                     # Go module definition
```

## Configuration

The server uses these defaults:

- **HTTP Timeout**: 10 seconds
- **Retry Count**: Maximum 2 retries for transient failures
- **Cache**: In-memory session cache (cleared on server shutdown)
- **Date Range**: YYYY (1900-2100), YYYYMMDD (valid calendar dates)

## Error Handling

All tools return detailed error messages:

- **ValidationError**: Invalid date format or out-of-range values
- **APIError**: Upstream API errors (5xx, connection issues)
- **Timeout**: Request exceeded 10-second timeout
- **ParseError**: Invalid JSON response from upstream API

## Testing

### Unit Tests

Test validation, parsing, and caching logic:

```bash
go test -v ./tests/unit/...
```

### Integration Tests

Test live API calls:

```bash
go test -v -run Integration ./tests/integration/...
```

### Coverage

Generate coverage report:

```bash
go test -cover ./...
```

## API Compliance

- **Principle I: API Contract Fidelity** - Tools are thin wrappers over Boatrace API
- **Principle II: Explicit Parameter Validation** - All inputs validated before API calls
- **Principle III: Structured Error Handling** - All errors caught and logged
- **Principle IV: Transport-Agnostic** - Works with any MCP transport
- **Principle V: Type Safety** - All responses unmarshalled to typed structs

## Contributing

See CONTRIBUTING.md for development guidelines.

## License

[Add appropriate license here]

## Support

For issues or questions:
- GitHub Issues: https://github.com/boatrace/open-api-mcp-server/issues
- Documentation: See `specs/001-mcp-tools/` for detailed specifications
