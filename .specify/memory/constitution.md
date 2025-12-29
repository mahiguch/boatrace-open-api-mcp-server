# Boatrace Open API MCP Server Constitution

<!--
  Sync Impact Report:
  - Version change: 0.0.0 → 1.0.0 (initial ratification)
  - Added 5 core principles aligned with MCP SDK development
  - Added section for API Integration Standards
  - Added section for Quality Assurance & Testing
  - All templates updated: ✅ spec-template, ✅ plan-template, ✅ tasks-template
  - No breaking changes; project is newly initialized
-->

## Core Principles

### I. API Contract Fidelity

Every MCP tool MUST accurately represent the underlying Boatrace Open API schema. Tools are thin, semantically transparent wrappers over HTTP endpoints—no data transformation or custom business logic.

**Rationale**: Users depend on MCP tools to reliably access official Boatrace data without surprises. Contracts are the user-facing API; deviation erodes trust and predictability.

### II. Explicit Parameter Validation

All MCP tool inputs (YYYY, YYYYMMDD date strings) MUST be validated before calling upstream APIs. Invalid dates are rejected with clear error messages, not passed upstream.

**Rationale**: Early validation prevents cryptic upstream errors, improves user experience, and keeps error handling localized. No silent failures.

### III. Structured Error Handling

All errors (network, validation, malformed responses) MUST be caught, logged with context, and returned to the user with actionable messages. No panics or unhandled exceptions in production code.

**Rationale**: MCP servers run in Claude's context; silent crashes degrade the user experience. Explicit error handling ensures debuggability and resilience.

### IV. Transport-Agnostic Implementation

The Go MCP server MUST be implementable over any MCP transport (stdio, HTTP, etc.) without code changes. Core tool logic is decoupled from transport concerns.

**Rationale**: Maximizes reusability and allows future migration (e.g., from studio transport to HTTP) without rewrites.

### V. Type Safety Over Convenience

Go code MUST use concrete, typed data structures (avoid `interface{}` or equivalent dynamic typing). All Boatrace API responses MUST be unmarshalled into well-defined structs.

**Rationale**: Compile-time type safety catches bugs early, improves code clarity, and prevents runtime type assertion panics. Go's strength is static typing—we lean into it.

## API Integration Standards

- All HTTP calls to Boatrace Open API MUST include appropriate timeouts (default: 10s)
- HTTP client MUST validate SSL certificates and follow redirects (max 5)
- All responses MUST be unmarshalled into typed Go structs (no raw JSON passes through)
- Failed external API calls MUST be retried up to 2 times before returning error to user
- API responses MUST be cached in memory for the lifetime of the MCP session to avoid redundant requests

## Quality Assurance & Testing

- **Unit tests**: All validation and parsing logic MUST have ≥90% code coverage
- **Integration tests**: All three MCP tools (programs, results, previews) MUST be tested against live API endpoints
- **Error cases**: Tests MUST cover network failures, invalid dates, malformed responses, and timeout scenarios
- **Linting**: Code MUST pass `golangci-lint` with no warnings
- **Documentation**: Every exported function MUST have a doc comment

## Governance

- This constitution supersedes all other development practices and is binding on all PRs and contributions
- Any deviation (e.g., untyped data, missing validation, untested features) MUST be justified in the PR description before review
- Principle amendments require a new version and explicit approval before implementation
- All pull requests MUST verify compliance with principles I–V before merge

**Version**: 1.0.0 | **Ratified**: 2025-12-22 | **Last Amended**: 2025-12-22
