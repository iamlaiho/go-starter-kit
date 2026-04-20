## ADDED Requirements

### Requirement: Standard directory structure
The repository SHALL follow golang-standards/project-layout with `cmd/`, `internal/`, `api/`, and `scripts/` top-level directories.

#### Scenario: Entry point exists
- **WHEN** a developer clones the repository
- **THEN** the server entry point SHALL be at `cmd/server/main.go`

#### Scenario: Internal packages are protected
- **WHEN** an external module attempts to import a package under `internal/`
- **THEN** the Go toolchain SHALL reject the import at compile time


### Requirement: Application layer structure
The `internal/` directory SHALL be organised into `handler/`, `service/`, `middleware/`, `config/`, `observability/`, and `testutil/` packages.

#### Scenario: Handler package exists
- **WHEN** a developer adds a new HTTP handler
- **THEN** it SHALL be placed under `internal/handler/<resource>/`

#### Scenario: Service layer exists
- **WHEN** a developer adds business logic
- **THEN** it SHALL be placed under `internal/service/`

