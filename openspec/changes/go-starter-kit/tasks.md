## 1. Repository Bootstrap

- [ ] 1.1 Run `go mod init` with the module path and set Go version to 1.22 in `go.mod`
- [ ] 1.2 Create top-level directories: `cmd/server/`, `internal/`, `api/`, `scripts/`
- [ ] 1.3 Create `internal/` subdirectories: `handler/`, `service/`, `middleware/`, `config/`, `observability/`, `testutil/`
- [ ] 1.4 Add `.gitignore` covering: compiled binaries (`/bin/`, `*.exe`, `*.test`), secrets (`.env`, `*.pem`, `*.key`), test output (`coverage/`, `*.coverprofile`), editor (`.vscode/`), OS (`.DS_Store`, `Thumbs.db`), and scratch dirs (`tmp/`); do NOT ignore `.env.example`; no vendor directory
- [ ] 1.5 Add `Makefile` with targets: `run`, `build`, `test`, `test-coverage`, `tidy`, `vuln`, `lint`, `docker-build`, `docker-run`, `hooks-install`, `dev`; `build` and `run` targets SHALL inject git SHA via `-ldflags "-X main.version=$(git rev-parse --short HEAD)"`
- [ ] 1.6 Add `.air.toml` configuration for Air live reload (watches `cmd/`, `internal/`; excludes `tmp/`, test files); add `make dev` target that runs `air`; document `go install github.com/air-verse/air@latest` as a prerequisite in README
- [ ] 1.7 Write `README.md` covering: prerequisites (Go 1.22+, Docker), quickstart (`make dev`), all `make` targets, environment variables overview (link to `.env.example`), and running tests
- [ ] 1.8 Add `.editorconfig` at repository root: tabs for `.go` files, 2-space indent for `.yml`/`.yaml`/`.json`/`.md`, LF line endings, UTF-8, trim trailing whitespace, insert final newline

## 2. Configuration

- [ ] 2.1 Add dependencies: `github.com/spf13/viper`
- [ ] 2.2 Create `internal/config/config.go` with typed `Config` struct covering `APP_ENV` (default: `development`), `PORT` (default: `8080`), `JWT_SECRET` (required), `MAX_REQUEST_BODY_BYTES` (default: `1048576`), `CORS_ALLOWED_ORIGINS` (default: `*` in development, required in production)
- [ ] 2.3 Implement `config.Load()` using viper to bind env vars and `.env` file
- [ ] 2.4 Add startup validation — exit code 1 with field name if required fields are missing
- [ ] 2.5 Create `.env.example` listing every config key with description and placeholder value

## 3. HTTP Server

- [ ] 3.1 Add dependencies: `github.com/go-chi/chi/v5`, `github.com/go-playground/validator/v10`
- [ ] 3.2 Create `cmd/server/main.go` as the entry point — init config, router, server
- [ ] 3.3 Create `internal/middleware/` — request-id, structured logger, recoverer, CORS middleware
- [ ] 3.4 Create `internal/middleware/timeout.go` — wraps `chi/middleware.Timeout(30s)`; apply in global middleware stack; return HTTP 503 on expiry
- [ ] 3.5 Create `internal/middleware/bodylimit.go` — wraps `http.MaxBytesReader` with configured `MAX_REQUEST_BODY_BYTES` (default 1MB); return HTTP 413 on oversized body
- [ ] 3.6 Create `internal/middleware/secureheaders.go` — sets `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`, `Content-Security-Policy: default-src 'self'` on every response
- [ ] 3.7 Create `internal/handler/errors.go` — `WriteError(w, status, message, requestID)` helper for consistent JSON error responses
- [ ] 3.8 Create `internal/handler/validate.go` — `Validate(v any) error` helper using `go-playground/validator`; returns formatted 400-ready error string per failing field
- [ ] 3.9 Add `GET /health` handler returning `{"status": "ok", "version": "<git-sha>"}` — unauthenticated, registered outside the auth middleware group; version is a `var version string` in `main.go` injected at build time via `ldflags`; defaults to `"dev"` if not set
- [ ] 3.10 Implement graceful shutdown in `cmd/server/main.go` — listen for SIGINT/SIGTERM, 30s timeout

## 4. Observability

- [ ] 4.1 Create `internal/observability/logger.go` — `NewLogger(env)` returning slog with JSON handler (production) or text handler (development)
- [ ] 4.2 Wire logger initialisation into `cmd/server/main.go`

## 5. Auth

- [ ] 5.1 Add dependencies: `github.com/golang-jwt/jwt/v5`
- [ ] 5.2 Create `internal/service/auth.go` — `GenerateTokenPair(userID string)` returning access (15m) and refresh (7d) JWT tokens with standard claims (`sub`, `exp`, `iat`, `jti`)
- [ ] 5.3 Create `internal/middleware/auth.go` — `Authenticate` middleware that validates Bearer token, attaches user ID to context; returns 401 for missing/expired/invalid tokens
- [ ] 5.4 Register an empty protected route group behind `Authenticate` middleware in `cmd/server/main.go`; teams register their protected routes here

## 6. Testing Infrastructure

- [ ] 6.1 Create `internal/testutil/server.go` — `NewTestServer(t, router)` helper returning an `httptest.Server` with the full middleware stack
- [ ] 6.2 Write tests for the `Authenticate` middleware: valid token grants access, missing/expired/invalid token returns 401
- [ ] 6.3 Configure `make test-coverage` to output HTML report to `coverage/index.html`

## 7. Docker

- [ ] 7.1 Create multi-stage `Dockerfile` — builder stage (`golang:1.22-alpine`), runtime stage (`distroless/static`), non-root user, copy binary
- [ ] 7.2 Create `docker-compose.yml` — app service with env from `.env` file; include a `healthcheck` directive calling `GET /health` every 30s with a 5s timeout and 3 retries
- [ ] 7.3 Add `.dockerignore` to exclude `.git`, `.env`, `coverage/`, test files from build context

## 8. CI/CD

- [ ] 8.1 Create `.github/workflows/pr-check.yml` — triggers on PR to `main`; caches Go modules; runs lint (`golangci-lint`), unit tests (`go test ./...`), build (`go build ./...`), and vulnerability scan (`govulncheck ./...`)
- [ ] 8.2 Create `.github/workflows/release.yml` — triggers on push of tags matching `v*`; builds Docker image and pushes to GHCR (`ghcr.io/${{ github.repository }}`) tagged with the semver tag and git SHA; authenticates using `GITHUB_TOKEN` (no secrets required)
- [ ] 8.3 Create `.golangci.yml` enabling `errcheck`, `govet`, `staticcheck`, `unused` with project-appropriate exclusions

## 9. Dependabot

- [ ] 9.1 Create `.github/dependabot.yml` with Go module ecosystem on weekly schedule, grouped by: `http` (chi, middleware), `auth` (golang-jwt), `dev` (golangci-lint, govulncheck)
- [ ] 9.2 Add GitHub Actions ecosystem to `dependabot.yml` on weekly schedule as a single `actions` group
- [ ] 9.3 Configure auto-merge for patch-level Dependabot PRs via `.github/workflows/dependabot-automerge.yml` — merge only when all CI checks pass, skip for minor/major bumps

## 10. Git Hooks

- [ ] 10.1 Add `lefthook` install instructions to Makefile (`hooks-install` target: `go install github.com/evilmartians/lefthook@latest && lefthook install`)
- [ ] 10.2 Create `lefthook.yml` at repository root configuring `commit-msg` (calls `scripts/check-commit-msg.sh {1}`) and `pre-commit` (runs `golangci-lint run`)
- [ ] 10.3 Create `scripts/check-commit-msg.sh` — reads commit message from `$1`, validates against Conventional Commits regex, exits 1 with helpful error message listing valid types and an example if invalid
- [ ] 10.4 Make `scripts/check-commit-msg.sh` executable (`chmod +x`)
- [ ] 10.5 Document hook setup in README (`make hooks-install` required after cloning)

## 11. Claude Code Safety Settings

- [ ] 11.1 Create `.claude/settings.json` registering `guard-risky-cmds.sh` as a `PreToolUse` hook for the Bash tool
- [ ] 11.2 Create `.claude/hooks/guard-risky-cmds.sh` — hard-block patterns: SSH key reads, pipe-to-shell, sudo, `/etc/` writes, force-push/reset-hard/clean-f, Docker prune/volume-rm, SSH tunnels, netcat listener; exit 2 with descriptive message for each
- [ ] 11.3 Set `.claude/hooks/guard-risky-cmds.sh` as executable (`chmod +x`)
- [ ] 11.4 Create `CLAUDE.md` at repository root documenting all blocked categories, rationale, and instructions for running blocked commands manually
