---
description: "Task list for Boatrace Open API MCP Server implementation"
---

# Tasks: Boatrace Open API MCP Server - Core Tools Implementation

**Input**: Design documents from `/specs/001-mcp-tools/`
**Prerequisites**: plan.md (required), spec.md (required), data-model.md, contracts/, research.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story. Each user story can be developed, tested, and deployed independently.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Root: `src/` for main code
- Tests: `tests/` at repository root
- Internal packages: `internal/` subdirectories
- Paths shown below assume single Go project structure

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic Go/MCP structure

- [ ] T001 Create Go module and initialize `go.mod` with `module github.com/yourusername/boatrace-open-api-mcp-server`
- [ ] T002 [P] Initialize project directory structure per plan.md: `cmd/`, `internal/`, `tests/`, and root files
- [ ] T003 [P] Add Go dependencies in `go.mod`: github.com/modelcontextprotocol/go-sdk and standard library imports
- [ ] T004 [P] Create root `main.go` MCP server entry point that initializes the Go MCP SDK
- [ ] T005 [P] Configure `.golangci.yml` linting rules for golangci-lint validation
- [ ] T006 [P] Create `Makefile` with build, test, lint, and clean targets
- [ ] T007 Create `README.md` with project overview, build instructions, and tool descriptions
- [ ] T008 Initialize git repository structure (`.gitignore`, initial commit)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T009 Create validation layer in `internal/validation/dates.go` with `ValidateYYYY()` and `ValidateYYYYMMDD()` functions validating date format and range (YYYY: 1900-2100, YYYYMMDD: valid calendar dates)
- [ ] T010 [P] Create data models in `internal/models/`:
  - [ ] T010a Create `internal/models/program.go` with RaceProgram, RaceInfo, RacerInfo, BettingInfo structs
  - [ ] T010b Create `internal/models/result.go` with RaceResult, RaceOutcome, RacerResult, PayoutInfo structs
  - [ ] T010c Create `internal/models/preview.go` with RacePreview, PreviewInfo, CurrentOddsInfo, PredictionRecord structs
- [ ] T011 Create HTTP client with configuration in `internal/api/client.go` with:
  - 10-second timeout configuration
  - SSL certificate validation
  - Redirect following (max 5)
  - Structured logging capability
- [ ] T012 Implement retry logic in `internal/api/client.go` with exponential backoff (100ms, 200ms) for transient failures (max 2 retries)
- [ ] T013 Create session-based caching layer in `internal/cache/cache.go` with in-memory map cache, key format `{tool}:{yyyymmdd}`, and session lifecycle management
- [ ] T014 [P] Implement error handling and logging:
  - [ ] T014a Create error types in `internal/errors.go` (ValidationError, APIError, TimeoutError, ParseError)
  - [ ] T014b Create structured logger wrapper in `internal/logger.go` using Go standard log package
- [ ] T015 Create MCP tool registration framework in `internal/tools/tools.go` for registering the three MCP tools with SDK

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Fetch Race Programs (Priority: P1) 🎯 MVP

**Goal**: Implement the `programs` MCP tool that fetches race program data for any valid date

**Independent Test**: Can be fully tested by calling `programs` with valid YYYY and YYYYMMDD, verifying structured RaceProgram data is returned with racer details

### Tests for User Story 1 (OPTIONAL - write tests FIRST, ensure they FAIL before implementation)

- [ ] T016 [P] [US1] Unit test for date validation in `tests/unit/validation_test.go` covering valid YYYY/YYYYMMDD, invalid formats, and boundary cases
- [ ] T017 [P] [US1] Unit test for programs API parsing in `tests/unit/models_test.go` testing RaceProgram JSON unmarshalling
- [ ] T018 [P] [US1] Mock HTTP test in `tests/unit/api_test.go` for programs endpoint with mock server (success, timeout, 5xx error, malformed JSON)
- [ ] T019 [P] [US1] Unit test for caching behavior in `tests/unit/cache_test.go` testing cache hit/miss, cache key generation
- [ ] T020 [US1] Integration test in `tests/integration/programs_test.go` calling live Boatrace Open API with valid date, verifying response structure

### Implementation for User Story 1

- [ ] T021 [P] [US1] Implement programs API endpoint in `internal/api/programs.go` with function `FetchPrograms(year, date string) (*RaceProgram, error)` that:
  - Validates year and date using validation layer
  - Checks cache before API call
  - Calls `GET https://boatraceopenapi.github.io/programs/v2/{year}/{date}.json` with retry logic
  - Unmarshals response into RaceProgram struct
  - Stores in cache with key `programs:{yyyymmdd}`
  - Returns RaceProgram or error
- [ ] T022 [US1] Implement programs MCP tool handler in `internal/tools/programs.go` that:
  - Receives `year` and `date` string parameters from MCP call
  - Calls FetchPrograms() from api package
  - Returns RaceProgram as structured response
  - Returns detailed error messages on failure
- [ ] T023 [US1] Register programs tool in `cmd/boatrace-mcp/main.go` with MCP SDK using tool handler from T022
- [ ] T024 [US1] Add logging in programs tool (T022) for: tool invocation, API call details, cache hit/miss, errors with context

**Checkpoint**: User Story 1 should be fully functional and independently testable. Test by calling `programs(year="2025", date="20251222")` and verify complete RaceProgram is returned.

---

## Phase 4: User Story 2 - Fetch Race Results (Priority: P2)

**Goal**: Implement the `results` MCP tool that fetches race results including winners and payouts for completed races

**Independent Test**: Can be fully tested by calling `results` with valid YYYY and YYYYMMDD, verifying structured RaceResult data with winners and payouts is returned

### Tests for User Story 2 (OPTIONAL - write tests FIRST, ensure they FAIL before implementation)

- [ ] T025 [P] [US2] Unit test for results API parsing in `tests/unit/models_test.go` testing RaceResult JSON unmarshalling with payout structures
- [ ] T026 [P] [US2] Mock HTTP test in `tests/unit/api_test.go` for results endpoint with mock server covering success, error scenarios
- [ ] T027 [US2] Integration test in `tests/integration/results_test.go` calling live Boatrace Open API with valid past date, verifying results structure

### Implementation for User Story 2

- [ ] T028 [P] [US2] Implement results API endpoint in `internal/api/results.go` with function `FetchResults(year, date string) (*RaceResult, error)` that:
  - Validates year and date using validation layer
  - Checks cache before API call
  - Calls `GET https://boatraceopenapi.github.io/results/v2/{year}/{date}.json` with retry logic
  - Unmarshals response into RaceResult struct
  - Stores in cache with key `results:{yyyymmdd}`
  - Returns RaceResult or error
- [ ] T029 [US2] Implement results MCP tool handler in `internal/tools/results.go` that:
  - Receives `year` and `date` string parameters from MCP call
  - Calls FetchResults() from api package
  - Returns RaceResult as structured response
  - Returns detailed error messages on failure
- [ ] T030 [US2] Register results tool in `cmd/boatrace-mcp/main.go` with MCP SDK using tool handler from T029
- [ ] T031 [US2] Add logging in results tool (T029) for: tool invocation, API call details, cache hit/miss, errors with context

**Checkpoint**: User Story 2 should be fully functional and independently testable. Test by calling `results(year="2025", date="20251220")` and verify complete RaceResult with winners/payouts is returned. Stories 1 & 2 should both work independently.

---

## Phase 5: User Story 3 - Fetch Race Previews (Priority: P3)

**Goal**: Implement the `previews` MCP tool that fetches race preview data including updated odds, predictions, and analysis

**Independent Test**: Can be fully tested by calling `previews` with valid YYYY and YYYYMMDD, verifying race preview information including odds and predictions is returned

### Tests for User Story 3 (OPTIONAL - write tests FIRST, ensure they FAIL before implementation)

- [ ] T032 [P] [US3] Unit test for previews API parsing in `tests/unit/models_test.go` testing RacePreview JSON unmarshalling with odds structures
- [ ] T033 [P] [US3] Mock HTTP test in `tests/unit/api_test.go` for previews endpoint with mock server covering success, error scenarios
- [ ] T034 [US3] Integration test in `tests/integration/previews_test.go` calling live Boatrace Open API with valid date, verifying preview structure

### Implementation for User Story 3

- [ ] T035 [P] [US3] Implement previews API endpoint in `internal/api/previews.go` with function `FetchPreviews(year, date string) (*RacePreview, error)` that:
  - Validates year and date using validation layer
  - Checks cache before API call
  - Calls `GET https://boatraceopenapi.github.io/previews/v2/{year}/{date}.json` with retry logic
  - Unmarshals response into RacePreview struct
  - Stores in cache with key `previews:{yyyymmdd}`
  - Returns RacePreview or error
- [ ] T036 [US3] Implement previews MCP tool handler in `internal/tools/previews.go` that:
  - Receives `year` and `date` string parameters from MCP call
  - Calls FetchPreviews() from api package
  - Returns RacePreview as structured response
  - Returns detailed error messages on failure
- [ ] T037 [US3] Register previews tool in `cmd/boatrace-mcp/main.go` with MCP SDK using tool handler from T036
- [ ] T038 [US3] Add logging in previews tool (T036) for: tool invocation, API call details, cache hit/miss, errors with context

**Checkpoint**: All user stories should now be independently functional. Test all three tools:
- `programs(year="2025", date="20251222")` returns RaceProgram
- `results(year="2025", date="20251220")` returns RaceResult
- `previews(year="2025", date="20251222")` returns RacePreview

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T039 [P] Run `golangci-lint` on all code in `internal/` and `cmd/` directories, fix all warnings
- [ ] T040 [P] Run full test suite: `go test ./...` in root directory, ensure all tests pass
- [ ] T041 [P] Run code coverage: `go test -cover ./...`, verify ≥90% coverage for validation, cache, models packages
- [ ] T042 [P] Add doc comments to all exported functions in `internal/` packages (every exported function must have doc comment)
- [ ] T043 Add integration test script `tests/integration/run_live_api_tests.sh` to run live API tests against current Boatrace API
- [ ] T044 Create `CONTRIBUTING.md` with development guidelines, testing procedures, and contribution workflow
- [ ] T045 Update `README.md` with complete installation, build, usage examples for all three tools
- [ ] T046 Create example usage in `examples/` directory with sample Go code calling each tool

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
  - ⚠️ **CRITICAL**: Must complete before any user story work begins
  - Validation, models, HTTP client, caching all needed by every tool
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can proceed in parallel (different files, no cross-story dependencies)
  - Or sequentially in priority order (P1 → P2 → P3)
  - Each story independently testable and deployable
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

| Story | Depends On | Can Start After | Notes |
|-------|-----------|-----------------|-------|
| **US1 (Programs)** | Phase 2 Foundational | Foundation complete | MVP - core functionality |
| **US2 (Results)** | Phase 2 Foundational | Foundation complete | Independent of US1 |
| **US3 (Previews)** | Phase 2 Foundational | Foundation complete | Independent of US1, US2 |

**Parallel Opportunities**:
- Phase 1: All setup tasks marked [P] can run in parallel
- Phase 2: All foundational tasks marked [P] can run in parallel (within phase)
- Phase 3+: All three user stories can be worked on in parallel by different developers
  - US1 development can proceed simultaneously with US2 and US3
  - US2 development can proceed simultaneously with US3
  - No cross-story dependencies; each story uses shared infrastructure from Phase 2

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- API endpoint implementation before MCP tool handler
- Tool handler before tool registration
- Core implementation before logging/monitoring
- Story complete before moving to next priority

### Parallel Example: User Story 1

```
Developer A:
  Task: T016 - Write unit tests for validation
  Task: T017 - Write unit tests for models
  Task: T018 - Write mock HTTP tests
  Task: T019 - Write cache tests

Developer B:
  Task: T021 - Implement programs API endpoint
  (T020 waits until implementation complete)

Developer C:
  Task: T022 - Implement programs MCP tool handler
  (depends on T021 complete)
  Task: T023 - Register tool (depends on T022 complete)
  Task: T024 - Add logging

All tests MUST pass (T016-T020) before T021-T024
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup ✅
2. Complete Phase 2: Foundational ✅ (BLOCKING)
3. Complete Phase 3: User Story 1 (programs tool) ✅
4. **STOP and VALIDATE**: Test User Story 1 independently
   ```bash
   go run ./cmd/boatrace-mcp &
   # Call: programs(year="2025", date="20251222")
   # Verify: RaceProgram returned with complete data
   ```
5. Deploy/demo MVP if ready

### Incremental Delivery (Recommended)

1. Complete Setup → Foundational → Foundation ready ✅
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add Polish & Cross-Cutting Concerns
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. **Developer 1**: Phase 1 (Setup) + Phase 2 (Foundational) core tasks
2. Once Foundational ready:
   - **Developer A**: User Story 1 (Programs)
   - **Developer B**: User Story 2 (Results)
   - **Developer C**: User Story 3 (Previews)
3. Stories complete and integrate independently
4. **All**: Polish & cross-cutting concerns together

---

## Notes

- [P] tasks = different files, no dependencies (can run parallel)
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- ✅ Verify tests fail BEFORE implementing each feature
- 📍 Stop at each checkpoint to validate story independently
- ✅ Commit after each logical task group
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Estimated total: ~1000-2000 lines of Go code across all phases
- Test-first approach: Write failing tests, then implement to make them pass
