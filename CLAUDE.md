# Claude Code — Project Guide

## Module

```
github.com/iamlaiho/go-starter-kit
```

## Commands

| Task | Command |
|------|---------|
| Run dev server | `make dev` |
| Build binary | `make build` |
| Run tests | `make test` |
| Coverage report | `make test-coverage` |
| Lint | `make lint` |
| Vuln scan | `make vuln` |
| Tidy deps | `make tidy` |
| Install git hooks | `make hooks-install` |
| Docker build | `make docker-build` |

## Layout

```
cmd/server/       — entry point (main.go)
internal/
  config/         — viper-based config loader
  handler/        — HTTP handlers + shared helpers (errors, validate)
  middleware/     — HTTP middleware stack
  observability/  — slog logger factory
  service/        — business logic (auth token generation)
  testutil/       — test helpers
scripts/          — git hook scripts
.github/          — CI/CD workflows + Dependabot
.claude/          — Claude Code hooks
```

## Key Conventions

- **No database layer** — this is a stateless API starter; add persistence yourself.
- **Auth** — JWT middleware only. Call `service.GenerateTokenPair` after your own credential check and return tokens to the client.
- **Config** — all config via environment variables (see `.env.example`). `JWT_SECRET` is required.
- **Logging** — use `slog.Default()` everywhere; JSON in production, text in development.
- **Error responses** — always use `handler.WriteError(w, status, message, requestID)`.
- **Validation** — use `handler.Validate(v)` before processing request bodies.
- **Commit style** — Conventional Commits enforced by lefthook pre-commit hook.

## Safety Hooks

`.claude/hooks/guard-risky-cmds.sh` hard-blocks destructive shell commands (`rm -rf`, force-push, hard-reset, etc.). Exit code 2 means no override.

## After cloning

```bash
cp .env.example .env   # fill in JWT_SECRET
go mod tidy
make hooks-install
make test
make dev
```
