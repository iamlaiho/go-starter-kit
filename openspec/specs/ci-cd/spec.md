## Purpose

Defines GitHub Actions CI/CD workflows: PR checks (lint, test, build, vuln scan) and release pipeline (Docker image push to GHCR on semver tag).

## Requirements

### Requirement: PR check workflow
The repository SHALL include a GitHub Actions workflow that runs on every pull request targeting `main`.

#### Scenario: Lint check runs
- **WHEN** a pull request is opened or updated
- **THEN** `golangci-lint` SHALL run against all changed and affected packages

#### Scenario: Unit tests run
- **WHEN** a pull request is opened or updated
- **THEN** `go test ./...` SHALL run and the workflow SHALL fail if any test fails

#### Scenario: Build check runs
- **WHEN** a pull request is opened or updated
- **THEN** `go build ./...` SHALL run to verify the project compiles

#### Scenario: Vulnerability scan runs
- **WHEN** a pull request is opened or updated
- **THEN** `govulncheck ./...` SHALL run and the workflow SHALL fail if any known CVEs are found in dependencies

#### Scenario: Workflow fails fast
- **WHEN** any step in the PR check workflow fails
- **THEN** subsequent steps SHALL be skipped and the workflow SHALL be marked as failed

### Requirement: Release workflow
The repository SHALL include a GitHub Actions workflow triggered on semver tag pushes that builds and pushes a Docker image to GHCR.

#### Scenario: Release triggered by semver tag
- **WHEN** a tag matching `v*` is pushed to the repository
- **THEN** the release workflow SHALL trigger and build the Docker image

#### Scenario: Image pushed to GHCR
- **WHEN** the Docker image builds successfully
- **THEN** the workflow SHALL push the image to `ghcr.io/${{ github.repository }}` tagged with the semver tag and git SHA

#### Scenario: GITHUB_TOKEN used for auth
- **WHEN** the release workflow pushes an image
- **THEN** it SHALL authenticate using the automatic `GITHUB_TOKEN` with no additional secrets required

### Requirement: Golangci-lint configuration
The repository SHALL include a `.golangci.yml` configuration file enabling a curated set of linters.

#### Scenario: Linter config present
- **WHEN** `golangci-lint` runs
- **THEN** it SHALL read `.golangci.yml` from the repository root and apply the configured linters and exclusions

#### Scenario: Minimum linters enabled
- **WHEN** golangci-lint runs
- **THEN** it SHALL at minimum run `errcheck`, `govet`, `staticcheck`, and `unused`
