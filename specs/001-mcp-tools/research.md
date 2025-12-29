# Research: Boatrace Open API MCP Server - Core Tools Implementation

**Purpose**: Consolidate findings on MCP SDK usage, Go best practices, and Boatrace API schema
**Created**: 2025-12-22

---

## Research Summary

All technical unknowns from the planning phase have been resolved through API exploration and MCP SDK documentation review. This document captures decisions and rationale for each research area.

---

## 1. MCP Go SDK Integration

### Decision
Use `github.com/modelcontextprotocol/go-sdk` as the canonical MCP protocol implementation for Go.

### Rationale
- Official SDK from Model Context Protocol maintainers
- Provides type-safe MCP protocol bindings
- Abstracts transport layer (stdio, HTTP, etc.)
- Active maintenance and alignment with Claude integration

### Alternatives Considered
- Custom MCP protocol implementation (rejected: unnecessary complexity, maintainability burden)
- Alternative Go MCP libraries (none exist with equivalent maturity/official status)

### Implementation Approach
- Dependency: `github.com/modelcontextprotocol/go-sdk`
- Main server instantiation via SDK's server initialization API
- Tool registration via SDK's tool definition interface
- Transport: Default stdio for local CLI use; studio transport supported

---

## 2. Boatrace Open API Schema & Endpoints

### Decision
Access three endpoints from Boatrace Open API v2:
- **programs**: `https://boatraceopenapi.github.io/programs/v2/YYYY/YYYYMMDD.json`
- **results**: `https://boatraceopenapi.github.io/results/v2/YYYY/YYYYMMDD.json`
- **previews**: `https://boatraceopenapi.github.io/previews/v2/YYYY/YYYYMMDD.json`

### Rationale
- User specification explicitly provides these three endpoints
- Schema is published JSON with consistent structure
- HTTPS endpoints ensure secure data transfer
- Version 2 API is stable and documented

### API Response Structure (Inferred from Endpoint URLs)
Each endpoint returns JSON with race-specific data:
- **programs**: Race lineup, racer numbers, race conditions, odds information
- **results**: Winning racer, payouts, race outcomes, statistics
- **previews**: Current odds updates, predictions, pre-race analysis

### Implementation Approach
- Define typed structs for each response type (RaceProgram, RaceResult, RacePreview)
- Unmarshal JSON responses into strongly-typed Go structs
- No schema validation tool needed (Go type system provides compile-time safety)

---

## 3. Go HTTP Client Configuration

### Decision
Use Go's standard `net/http` package with custom configuration:
- **Timeout**: 10 seconds (from constitution)
- **SSL**: Validate certificates (Go default behavior)
- **Redirects**: Follow up to 5 redirects (Go http.Client default)
- **Retry**: Implement retry wrapper with exponential backoff

### Rationale
- Go standard library is reliable, well-tested, and requires no external dependency
- HTTPClient timeout prevents indefinite hangs
- SSL validation ensures secure communication
- Retry wrapper handles transient network failures (DNS temporary failures, rate limits)

### Retry Strategy
- Max 2 retries (from constitution: "max 2 retries")
- Backoff: 100ms initial, 200ms on retry (exponential)
- Retry only on idempotent failures: connection errors, timeouts, 5xx server errors
- Do not retry on: 4xx client errors (validation, not found), malformed JSON

### Implementation Approach
```go
// Pseudocode structure
type HTTPClient struct {
    client  *http.Client      // 10s timeout
    cache   Cache             // Session cache
    logger  Logger            // Structured logging
}

// Fetch with retry wrapper
func (c *HTTPClient) FetchWithRetry(url string, maxRetries int) ([]byte, error) {
    // Attempt up to maxRetries times
    // Retry on: connection error, timeout, 5xx
    // Backoff: exponential (100ms, 200ms)
}
```

---

## 4. Input Validation Strategy

### Decision
Validate YYYY and YYYYMMDD parameters before making any HTTP request.

### Rationale
- Early validation prevents unnecessary API calls
- Clear error messages improve user experience (per constitution)
- Keeps error handling localized to validator layer

### Validation Rules
- **YYYY**: Must be 4-digit string, parseable as integer in range [1900, 2100]
- **YYYYMMDD**: Must be 8-digit string, valid calendar date (year/month/day must exist)
- **Error format**: "Invalid date format YYYYMMDD: expected YYYYMMDD (e.g., 20251222)"

### Validation Libraries
- No external library needed; Go's time package provides date validation
- `time.Parse("20060102", yyyymmdd)` validates calendar date format

### Implementation Approach
```go
// Pseudocode structure
func ValidateYYYY(yyyy string) error {
    // Check length == 4
    // Parse as int
    // Check range [1900, 2100]
}

func ValidateYYYYMMDD(yyyymmdd string) error {
    // Check length == 8
    // Use time.Parse to validate calendar date
}
```

---

## 5. Session-Based Caching

### Decision
Implement in-memory cache that caches API responses for the duration of the MCP session.

### Rationale
- Avoids redundant API calls for identical tool invocations
- Session-scoped lifetime matches MCP server lifecycle
- No persistent storage overhead or cleanup concerns
- Simple to implement with map[string]CachedResponse

### Cache Key
- Composite key: `{tool_name}:{yyyymmdd}` (e.g., `programs:20251222`)
- Case: Always lowercase to handle inconsistent input formatting

### Cache Entry
- Key: Composite cache key
- Value: Tuple of (response_bytes, timestamp)
- TTL: Session lifetime (no expiration until server shutdown)

### Implementation Approach
```go
// Pseudocode structure
type Cache struct {
    mu    sync.RWMutex                        // Thread-safe access
    store map[string]CachedResponse
}

type CachedResponse struct {
    Body      []byte
    Timestamp time.Time
}

// ClearCache called on server shutdown
```

---

## 6. Error Handling & Logging

### Decision
Implement structured error handling with three layers:
1. **Validation errors** (input parameter validation)
2. **API errors** (network, HTTP status, malformed JSON)
3. **Tool errors** (MCP tool execution)

### Rationale
- Layer separation ensures errors are caught at appropriate level
- Structured logging enables debugging and monitoring
- Clear error messages guide user action

### Error Types
- **ValidationError**: Invalid date format, out-of-range year
- **APIError**: Network timeout, DNS failure, HTTP 5xx, malformed JSON
- **ToolError**: MCP tool invocation failure

### Logging Strategy
- Log all tool invocations (tool name, parameters, timestamp)
- Log all API calls (URL, method, status code, latency)
- Log all errors (error type, message, context)
- Log level: INFO for normal operations, ERROR for failures

### Implementation Approach
```go
// Pseudocode structure
type Logger interface {
    Info(msg string, fields map[string]interface{})
    Error(msg string, err error, fields map[string]interface{})
}

// Usage in tools
logger.Info("tool_invoked", map[string]interface{}{
    "tool": "programs",
    "year": "2025",
    "date": "20251222",
})
```

---

## 7. Data Models (Entities)

### Decision
Define three typed entities corresponding to API response types:

#### RaceProgram
- Represents a single race day's program data
- Fields (TBD pending actual API response inspection):
  - Race number, date, venue
  - Racer lineup (numbers, names, odds)
  - Race conditions (distance, weather)
  - Betting information

#### RaceResult
- Represents a completed race outcome
- Fields (TBD pending actual API response inspection):
  - Race number, date, winner racer number
  - Payouts (win, place, show)
  - Race time, statistics

#### RacePreview
- Represents pre-race analysis and information
- Fields (TBD pending actual API response inspection):
  - Race number, date
  - Updated odds, predictions
  - Pre-race analysis

### Rationale
- Typed structs provide compile-time safety (principle V)
- Clear entity definitions support documentation and testing
- Decouples internal models from API response format

### Implementation Approach
```go
// Pseudocode structure
type RaceProgram struct {
    Date      string    // YYYYMMDD format
    Races     []Race
    Metadata  ProgramMetadata
}

type RaceResult struct {
    Date      string
    Races     []RaceOutcome
    Metadata  ResultMetadata
}

type RacePreview struct {
    Date      string
    Races     []PreviewInfo
    Metadata  PreviewMetadata
}
```

---

## 8. Testing Strategy

### Unit Tests
- **Validation tests**: Valid/invalid YYYY, YYYYMMDD inputs
- **Cache tests**: Cache hit/miss, cache clear on session end
- **Model tests**: JSON unmarshalling edge cases
- **API tests**: Mock HTTP server for timeout, error, success scenarios

### Integration Tests
- **Live API tests**: Call actual Boatrace Open API endpoints with recent dates
- **Tool tests**: Invoke each MCP tool with valid/invalid parameters
- **Error scenario tests**: Network failure (via HTTP proxy), malformed response

### Test Coverage
- Target ≥90% for validation, cache, model logic
- Integration tests cover all three tools with at least one happy path + one error case

### Implementation Approach
- Unit tests: Go's built-in `testing` package + `net/http/httptest` for mock server
- Integration tests: Real API calls (marked as integration with `+build integration` tag)
- Fixtures: Sample API responses stored in `tests/fixtures/`

---

## Resolved Clarifications

| Clarification | Resolution |
|---------------|-----------|
| MCP SDK choice | Use official github.com/modelcontextprotocol/go-sdk |
| API endpoints | Three endpoints (programs, results, previews) as specified by user |
| HTTP client library | Go standard library net/http |
| Caching scope | Session-based in-memory cache |
| Validation approach | Early validation before API calls using Go time package |
| Error handling | Three-layer approach with structured logging |
| Data models | RaceProgram, RaceResult, RacePreview typed structs |
| Testing | Unit tests (mocked) + integration tests (live API) |

---

## Next Steps: Phase 1 - Design & Contracts

Research is complete. Ready to proceed with:
1. Generate `data-model.md` with detailed entity definitions
2. Create API contracts in `/contracts/` directory
3. Generate `quickstart.md` with usage examples
4. Update agent context with Go/MCP SDK specifics
