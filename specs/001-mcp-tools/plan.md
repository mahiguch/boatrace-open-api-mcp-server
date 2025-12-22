# Implementation Plan: Boatrace Open API MCP Server - Core Tools Implementation

**Branch**: `001-mcp-tools` | **Date**: 2025-12-22 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-mcp-tools/spec.md`

## Summary

The Boatrace Open API MCP Server provides three MCP tools (`programs`, `results`, `previews`) that enable Claude users to access official Boatrace race data. Each tool accepts YYYY (year) and YYYYMMDD (date) parameters and returns structured race information via HTTP calls to the Boatrace Open API. The implementation prioritizes API contract fidelity, explicit input validation, structured error handling, type safety, and caching to ensure reliable, performant access to official data.

**Technical approach**: Go-based MCP server using the Go MCP SDK, implementing HTTP client with validation, error handling, retry logic, and in-memory session caching. Each tool is implemented as an independent MCP tool with clear separation of concerns.

## Technical Context

**Language/Version**: Go 1.21 or later (standard MCP SDK requirement)
**Primary Dependencies**:
- github.com/modelcontextprotocol/go-sdk (MCP protocol implementation)
- Go standard library: net/http, encoding/json, time (HTTP client, JSON unmarshalling, timeouts)
- Optional: logrus or Go standard log for structured logging

**Storage**: N/A (session-based in-memory caching only; no persistent storage)
**Testing**: Go testing (testing package with http.MockServer for unit tests; live API calls for integration tests)
**Target Platform**: Linux/macOS CLI server (deployable to any platform with Go runtime)
**Project Type**: Single CLI project (Go server executable)
**Performance Goals**:
- Tool invocation response time: <15 seconds (including network round trip)
- Cache hit response time: <100ms
- No specific throughput requirement (single-user MCP session)

**Constraints**:
- HTTP request timeout: default 10 seconds
- Cache lifecycle: session duration only (no persistent cache)
- Date validation: YYYY (1900-2100), YYYYMMDD (valid calendar dates)
- Retry limit: max 2 retries for transient failures

**Scale/Scope**: Single MCP server binary, ~1000-2000 lines of Go code

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Core Principles Alignment ✅

| Principle | Requirement | Status | Notes |
|-----------|------------|--------|-------|
| **I. API Contract Fidelity** | Tools are thin wrappers over Boatrace API; no custom business logic | ✅ PASS | Tools marshal API responses into typed structs; no transformation |
| **II. Explicit Parameter Validation** | All YYYY/YYYYMMDD inputs validated before upstream calls | ✅ PASS | Validation layer validates dates (1900-2100 range, valid calendar dates) |
| **III. Structured Error Handling** | All errors caught, logged, returned with actionable messages | ✅ PASS | Error handling for network, validation, timeout, malformed response |
| **IV. Transport-Agnostic Implementation** | Core logic independent of transport; works with any MCP transport | ✅ PASS | Tool implementation decoupled from MCP transport (stdio/HTTP/studio) |
| **V. Type Safety Over Convenience** | No `interface{}` or dynamic typing; concrete typed structs | ✅ PASS | All API responses unmarshalled into well-defined Go structs |

### API Integration Standards ✅

| Standard | Requirement | Status | Notes |
|----------|------------|--------|-------|
| HTTP Timeouts | 10-second default on all external calls | ✅ PASS | Implemented via http.Client timeout configuration |
| SSL Validation | Certificates validated; redirects followed (max 5) | ✅ PASS | Go http.Client validates by default; redirect limit configurable |
| Typed Responses | No raw JSON; all responses unmarshalled to structs | ✅ PASS | Custom structs for RaceProgram, RaceResult, RacePreview entities |
| Retry Logic | Up to 2 retries for failed API calls | ✅ PASS | Retry wrapper on HTTP client for transient failures |
| Session Caching | API responses cached for session lifetime | ✅ PASS | In-memory cache map with date as key, cleared on session end |

### Quality Assurance & Testing ✅

| Requirement | Status | Notes |
|-----------|--------|-------|
| Unit test coverage ≥90% | ✅ PASS | All validation, parsing, caching logic covered by unit tests |
| Integration tests (live API) | ✅ PASS | Integration tests call live Boatrace Open API endpoints |
| Error case coverage | ✅ PASS | Tests for network failures, invalid dates, malformed responses, timeouts |
| Linting (golangci-lint) | ✅ PASS | All code must pass golangci-lint with no warnings |
| Doc comments | ✅ PASS | Every exported function has doc comment |

### Gate Result: ✅ PASS

All constitution principles and standards are met by the technical design. No violations or exceptions needed. Plan approved for Phase 0 research.

## Project Structure

### Documentation (this feature)

```text
specs/001-mcp-tools/
├── plan.md                   # This file (implementation plan)
├── spec.md                   # Feature specification
├── research.md               # Phase 0: Research findings (TBD)
├── data-model.md             # Phase 1: Data model entities (TBD)
├── quickstart.md             # Phase 1: Quickstart guide (TBD)
├── contracts/                # Phase 1: API contracts (TBD)
│   ├── programs-tool.md      # programs tool contract
│   ├── results-tool.md       # results tool contract
│   └── previews-tool.md      # previews tool contract
├── checklists/
│   └── requirements.md       # Specification quality checklist
└── tasks.md                  # Phase 2: Task breakdown (via /speckit.tasks)
```

### Source Code (repository root)

**Structure Decision**: Single CLI project with Go module structure

```text
.
├── go.mod                    # Go module definition
├── go.sum                    # Go dependencies lock file
├── main.go                   # MCP server entry point
├── cmd/
│   └── boatrace-mcp/        # CLI application
│       └── main.go          # Server initialization
├── internal/
│   ├── models/              # Data models (RaceProgram, RaceResult, RacePreview)
│   │   ├── program.go       # RaceProgram entity definition
│   │   ├── result.go        # RaceResult entity definition
│   │   └── preview.go       # RacePreview entity definition
│   ├── api/                 # Boatrace API client
│   │   ├── client.go        # HTTP client with retry/caching/timeouts
│   │   ├── programs.go      # Programs endpoint implementation
│   │   ├── results.go       # Results endpoint implementation
│   │   └── previews.go      # Previews endpoint implementation
│   ├── validation/          # Input validation
│   │   └── dates.go         # Date format and range validation
│   ├── cache/               # Session cache
│   │   └── cache.go         # In-memory cache with session lifecycle
│   └── tools/               # MCP tool handlers
│       ├── programs.go      # programs MCP tool handler
│       ├── results.go       # results MCP tool handler
│       └── previews.go      # previews MCP tool handler
├── tests/
│   ├── unit/                # Unit tests
│   │   ├── validation_test.go
│   │   ├── cache_test.go
│   │   ├── models_test.go
│   │   └── api_test.go
│   ├── integration/         # Integration tests (live API)
│   │   ├── programs_test.go
│   │   ├── results_test.go
│   │   └── previews_test.go
│   └── fixtures/            # Test fixtures and mock data
│       └── sample_responses.go
├── .golangci.yml            # golangci-lint configuration
├── Makefile                 # Build targets
└── README.md                # Project documentation
```

## Complexity Tracking

> **No Constitution Check violations identified. No complexity tracking required.**

All design decisions align with project constitution. Implementation is straightforward with clear separation of concerns (validation → HTTP client → MCP tools).
