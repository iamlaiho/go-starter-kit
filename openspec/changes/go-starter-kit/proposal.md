## Why

Starting a Go project from scratch requires assembling a large number of boilerplate decisions — project layout, HTTP routing, config loading, auth, observability, and CI/CD — all consistently and correctly. This starter kit eliminates that setup cost and establishes opinionated, production-ready defaults so teams can begin writing business logic on day one.

## What Changes

- Introduce a new repository scaffold with a standard Go project layout (`cmd/`, `internal/`, `api/`, `scripts/`)
- Add an HTTP server layer with routing, middleware chaining, and structured error responses
- Add configuration management supporting environment variables and `.env` files
- Add JWT-based authentication middleware
- Add structured logging with `slog`
- Add Dockerfile and `docker-compose.yml` for local development and production
- Add test infrastructure with helpers for unit tests
- Add CI/CD pipeline via GitHub Actions (lint, test, build)
- Add a `Makefile` with common dev tasks (`make run`, `make test`, `make lint`, etc.)
- Add `/health` endpoint for deployment readiness probes
- Add Air live reload for hot-reloading during development
- Add README with prerequisites, make targets, and local dev guide
- Add security headers middleware (X-Frame-Options, X-Content-Type-Options, etc.)
- Add request input validation using `go-playground/validator` with structured error responses
- Add request timeout and body size limit middleware for production safety

## Capabilities

### New Capabilities
- `project-layout`: Standard Go project structure following `golang-standards/project-layout` with `cmd/`, `internal/`, and `api/` directories
- `http-server`: HTTP server using the standard `net/http` with a lightweight router (Chi), middleware stack (logging, recovery, CORS, request-id), and structured JSON error responses
- `config`: Configuration loading from environment variables and `.env` files using `viper`, with a typed config struct and validation at startup
- `auth`: JWT middleware for protecting routes and a `GenerateTokenPair(userID)` helper; teams supply their own login handler and credential verification
- `observability`: Structured logging with `slog` (Go stdlib), JSON in production, text in development, with request logging middleware
- `docker`: Multi-stage `Dockerfile` (builder + distroless runtime), `docker-compose.yml` for local development
- `testing`: Test helpers using `httptest` for handler tests and coverage reporting
- `ci-cd`: GitHub Actions workflows for PR checks (lint with `golangci-lint`, unit tests, build) and release (Docker image push to registry)
- `dependabot`: Grouped Dependabot configuration for Go modules (by category: http, auth, dev) and GitHub Actions, with auto-merge enabled for patch-level bumps that pass CI
- `claude-settings`: Claude Code project settings (`.claude/settings.json`) with hard-block hooks for risky commands — SSH/credential exposure, pipe-to-shell, sudo, force-push, and Docker prune; no override mechanism
- `git-hooks`: Lefthook for Git hook management with a shell-based `commit-msg` validator enforcing Conventional Commits format; no Node.js required

### Modified Capabilities
<!-- None — this is a greenfield starter kit -->

## Impact

- **New repository**: Entire codebase is new; no existing code is modified
- **Dependencies introduced**: `chi`, `viper`, `golang-jwt/jwt`, `go-playground/validator`, `golangci-lint` (dev)
- **GitHub config**: `.github/dependabot.yml` for automated dependency PRs with grouped updates and auto-merge
- **Claude Code config**: `.claude/settings.json` + `.claude/hooks/guard-risky-cmds.sh` for hard-block enforcement of dangerous commands
- **Git hooks**: `lefthook.yml` + `scripts/check-commit-msg.sh` enforcing Conventional Commits; no Node.js dependency
- **Go version**: 1.22+ (uses stdlib `slog`, `net/http` enhancements)
- **Infrastructure**: Docker required for local dev
