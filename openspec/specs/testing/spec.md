## Purpose

Defines the testing strategy: handler tests using httptest, coverage reporting, and shared test helpers.

## Requirements

### Requirement: Handler tests with httptest
HTTP handlers SHALL be tested using Go's `net/http/httptest` package with the full middleware stack.

#### Scenario: Handler test uses full middleware stack
- **WHEN** a developer writes a handler test
- **THEN** they SHALL use `NewTestServer(t, router)` from `internal/testutil/` which wraps the handler in an `httptest.Server` with all middleware applied

#### Scenario: Auth middleware tests cover success and failure
- **WHEN** a developer runs `go test ./...`
- **THEN** middleware tests SHALL cover the `Authenticate` middleware for valid token access and 401 on missing/expired/invalid tokens

### Requirement: Test coverage reporting
The CI pipeline and Makefile SHALL support generating test coverage reports.

#### Scenario: Coverage report generated
- **WHEN** a developer runs `make test-coverage`
- **THEN** the command SHALL run all unit tests and output an HTML coverage report to `coverage/index.html`

### Requirement: Test helpers
The repository SHALL provide shared test helpers for common setup patterns.

#### Scenario: HTTP test helper exists
- **WHEN** a developer writes a handler test
- **THEN** they SHALL find a helper in `internal/testutil/` that creates a test HTTP server with the full middleware stack
