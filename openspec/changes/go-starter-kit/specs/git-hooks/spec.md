## ADDED Requirements

### Requirement: Lefthook hook manager
The repository SHALL use `lefthook` as the Git hook manager, configured via `lefthook.yml` at the repository root.

#### Scenario: Lefthook installed via go install
- **WHEN** a developer runs `make hooks-install`
- **THEN** the Makefile target SHALL install lefthook via `go install` and run `lefthook install` to register hooks in `.git/hooks/`

#### Scenario: Hooks active after install
- **WHEN** `lefthook install` has been run
- **THEN** Git SHALL invoke lefthook on `commit-msg` and `pre-commit` events automatically

### Requirement: Conventional Commits enforcement
The `commit-msg` hook SHALL reject commit messages that do not conform to the Conventional Commits format.

#### Scenario: Valid commit message accepted
- **WHEN** a developer commits with a message matching `^(feat|fix|chore|docs|refactor|test|ci|build|perf|style|revert)(\(.+\))?: .+`
- **THEN** the hook SHALL exit 0 and the commit SHALL proceed

#### Scenario: Invalid commit message rejected
- **WHEN** a developer commits with a message that does not match the pattern (e.g., `"updated stuff"`)
- **THEN** `scripts/check-commit-msg.sh` SHALL exit 1 with a message showing the required format and valid types

#### Scenario: Scope is optional
- **WHEN** a developer commits with a message like `feat: add login endpoint` (no scope)
- **THEN** the hook SHALL accept the message

#### Scenario: Scope is validated when present
- **WHEN** a developer commits with a message like `feat(auth): add login endpoint`
- **THEN** the hook SHALL accept the message

### Requirement: Pre-commit quality checks
The `pre-commit` hook SHALL run linting on staged Go packages before each commit.

#### Scenario: Lint runs on pre-commit
- **WHEN** a developer runs `git commit`
- **THEN** `golangci-lint run` SHALL execute and the commit SHALL be blocked if any lint errors are found

#### Scenario: Pre-commit can be skipped
- **WHEN** a developer runs `git commit` with `LEFTHOOK=0` set
- **THEN** all hooks SHALL be skipped (for emergency commits)

### Requirement: Commit message validator script
The repository SHALL include `scripts/check-commit-msg.sh` implementing the Conventional Commits validation without external dependencies.

#### Scenario: Script reads commit message file
- **WHEN** the `commit-msg` hook fires
- **THEN** the script SHALL read the commit message from the file path passed as `$1` (Git's standard commit-msg argument)

#### Scenario: Script outputs helpful error
- **WHEN** the commit message is invalid
- **THEN** the script SHALL print the invalid message, list valid types, and show an example valid message
