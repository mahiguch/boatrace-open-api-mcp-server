# Feature Specification: Boatrace Open API MCP Server - Core Tools Implementation

**Feature Branch**: `001-mcp-tools`
**Created**: 2025-12-22
**Status**: Draft
**Input**: User description: "Boatrace Open API MCP Server の初期実装 - programs、results、previews ツールの実装"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Fetch Race Programs (Priority: P1)

Claude users need to retrieve race program data (出走表) to understand upcoming races, including racer information, race conditions, and odds predictions. The `programs` tool provides access to detailed program information for any date.

**Why this priority**: Programs data is the foundational information needed to understand any race. Without this, other information (results, previews) has no context. This is the MVP tool.

**Independent Test**: Can be fully tested by calling `programs` with a valid YYYY and YYYYMMDD, verifying that structured race program data is returned with racer details and race information.

**Acceptance Scenarios**:

1. **Given** a user has access to the MCP server, **When** they call the `programs` tool with year "2025" and date "20251222", **Then** the tool returns a complete race program for that date with racer lineups and race metadata
2. **Given** the programs tool has returned data for a date, **When** the user inspects the response, **Then** all races and racer information is accurately represented with no missing fields
3. **Given** a user calls programs with a valid date, **When** the system makes the HTTP request, **Then** the response is cached in memory to avoid redundant API calls during the session

---

### User Story 2 - Fetch Race Results (Priority: P2)

Claude users need to retrieve race results (結果) to see outcomes, winner information, payouts, and race statistics. The `results` tool provides access to completed race data for any date.

**Why this priority**: Results are critical for analysis and understanding past performance, but depend on races being defined first (programs). This can be independently tested once programs tool is working.

**Independent Test**: Can be fully tested by calling `results` with a valid YYYY and YYYYMMDD, verifying that structured race results data is returned with winners, payouts, and race outcomes.

**Acceptance Scenarios**:

1. **Given** a user has access to the MCP server, **When** they call the `results` tool with year "2025" and date "20251222", **Then** the tool returns complete race results for that date including winner information and payouts
2. **Given** the results tool has returned data, **When** the user inspects the response, **Then** all race outcomes are accurately represented with winning racer numbers and payouts
3. **Given** a user calls results with a valid date, **When** the system makes the HTTP request, **Then** the response is cached to prevent redundant calls

---

### User Story 3 - Fetch Race Previews (Priority: P3)

Claude users need to retrieve race preview data (直前情報) including latest odds, predictions, and last-minute race information. The `previews` tool provides access to pre-race analysis and information.

**Why this priority**: Preview information is valuable for current analysis but less critical than programs (which define races) or results (which show outcomes). This feature enhances the offering but isn't essential for basic functionality.

**Independent Test**: Can be fully tested by calling `previews` with a valid YYYY and YYYYMMDD, verifying that race preview information including odds and predictions is returned.

**Acceptance Scenarios**:

1. **Given** a user has access to the MCP server, **When** they call the `previews` tool with year "2025" and date "20251222", **Then** the tool returns preview information for that date including odds and pre-race predictions
2. **Given** the previews tool has returned data, **When** the user inspects the response, **Then** all odds, predictions, and race preview data is present and properly formatted
3. **Given** a user calls previews with a valid date, **When** the system makes the HTTP request, **Then** the response is cached for the session duration

---

### Edge Cases

- What happens when a date has no races scheduled? (System returns empty array with appropriate status)
- What happens when an invalid date format (e.g., "20259999" or "abcd") is provided? (System rejects with validation error before calling upstream API)
- What happens when the Boatrace API is temporarily unavailable? (System returns error with clear message after retry attempts)
- What happens when a date is too far in the past or future? (System validates date bounds and returns appropriate error)
- What happens when the MCP session ends? (Cached data is cleared from memory; next session starts fresh)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a `programs` MCP tool that accepts YYYY (year) and YYYYMMDD (date) as parameters and returns race program data from Boatrace Open API
- **FR-002**: System MUST provide a `results` MCP tool that accepts YYYY (year) and YYYYMMDD (date) as parameters and returns race results data from Boatrace Open API
- **FR-003**: System MUST provide a `previews` MCP tool that accepts YYYY (year) and YYYYMMDD (date) as parameters and returns race preview data from Boatrace Open API
- **FR-004**: System MUST validate all date inputs (YYYY format: 1900-2100, YYYYMMDD format: valid calendar dates) before calling upstream API
- **FR-005**: System MUST return validation errors with clear messages when date parameters are invalid or out of bounds
- **FR-006**: System MUST implement HTTP caching to avoid redundant API calls for the same date within a single MCP session
- **FR-007**: System MUST handle network errors, timeouts, and malformed responses from Boatrace Open API with clear error messages returned to user
- **FR-008**: System MUST implement retry logic (max 2 retries) for transient failures before returning error
- **FR-009**: System MUST represent all Boatrace API responses using strongly-typed data structures (no raw JSON or generic objects)
- **FR-010**: System MUST include HTTP request timeouts (default 10 seconds) for all external API calls
- **FR-011**: System MUST validate SSL certificates for all HTTPS connections to Boatrace Open API
- **FR-012**: System MUST log all tool invocations, API calls, and errors with appropriate context for debugging

### Key Entities

- **RaceProgram**: Represents a single race day's program data including race number, racer lineups, race conditions, betting odds, and race metadata
- **RaceResult**: Represents a completed race outcome including winning racer number, payouts, winning odds, race statistics, and time results
- **RacePreview**: Represents pre-race information for upcoming races including odds updates, predictions, race analysis, and late-breaking information
- **DateRange**: Represents valid date bounds (YYYY year: 1900-2100; YYYYMMDD: valid calendar dates in ISO 8601 format)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All three MCP tools (programs, results, previews) are functional and can be called with valid dates, returning complete data within 15 seconds (including network latency and caching)
- **SC-002**: Invalid date inputs are rejected with clear error messages before any upstream API call is made (100% of invalid inputs caught)
- **SC-003**: Repeated calls to the same tool with the same date return cached results within 100ms (cache hit time)
- **SC-004**: Network failures and API errors are handled gracefully with retry logic, succeeding on at least 2 of 3 attempts for transient failures
- **SC-005**: All tool inputs and outputs match the Boatrace Open API schema with no data loss or transformation
- **SC-006**: Documentation clearly describes each tool's parameters, return values, error scenarios, and usage examples
- **SC-007**: All three tools can be independently tested and deployed without dependencies on each other

## Assumptions

- Boatrace Open API endpoints remain stable and accessible at the URLs specified
- The API response format matches the documented JSON schema without breaking changes during development
- MCP studio transport is the primary deployment environment (transport-agnostic design allows future flexibility)
- Date parameters follow standard format conventions: YYYY (4-digit year), YYYYMMDD (8-digit date)
- System operates in a stateless manner where each session maintains its own cache
- Users have basic familiarity with date formats and will provide valid inputs with occasional errors
