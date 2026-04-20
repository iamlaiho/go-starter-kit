## Purpose

Defines Dependabot configuration for automated, grouped dependency updates with auto-merge for patch-level bumps that pass CI.

## Requirements

### Requirement: Grouped Go module updates
Dependabot SHALL be configured to open grouped pull requests for Go module dependencies on a weekly schedule.

#### Scenario: Dependencies grouped by category
- **WHEN** Dependabot runs its weekly check
- **THEN** it SHALL open at most one PR per group: `http` (chi, middleware), `auth` (golang-jwt), and `dev` (golangci-lint, govulncheck)

#### Scenario: Minor and patch updates only
- **WHEN** Dependabot detects an available update
- **THEN** it SHALL open a PR for patch and minor version bumps, and leave major version bumps to be handled manually

### Requirement: Grouped GitHub Actions updates
Dependabot SHALL be configured to open a single grouped pull request for all GitHub Actions dependency updates on a weekly schedule.

#### Scenario: Actions pinned to SHA
- **WHEN** Dependabot opens a GitHub Actions update PR
- **THEN** all action references SHALL be updated to the new SHA and tag simultaneously in one PR

### Requirement: Auto-merge for patch updates
Dependabot pull requests for patch-level bumps SHALL be automatically merged when CI passes.

#### Scenario: Patch PR auto-merged
- **WHEN** Dependabot opens a PR for a patch version bump and all CI checks pass
- **THEN** the PR SHALL be automatically merged without requiring human review

#### Scenario: Minor PR requires review
- **WHEN** Dependabot opens a PR for a minor version bump
- **THEN** the PR SHALL NOT be auto-merged and SHALL require a human reviewer to approve

#### Scenario: Auto-merge requires CI
- **WHEN** any CI check fails on a Dependabot patch PR
- **THEN** auto-merge SHALL NOT proceed and the PR SHALL remain open for manual inspection
