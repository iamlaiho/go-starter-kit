## Context

This is a greenfield Go starter kit repository. There is no existing code — all architectural decisions are made from scratch. The goal is to give teams a production-ready foundation they can fork and immediately build on without rethinking boilerplate.

Constraints:
- Go 1.22+ (stdlib `slog`, enhanced `net/http` ServeMux)
- Docker required for local development
- GitHub Actions for CI/CD
- Opinionated but minimal — every included library must earn its place

## Goals / Non-Goals

**Goals:**
- Provide a complete, runnable Go project scaffold out of the box
- Establish consistent patterns: HTTP handler → service
- Cover the full "day one" surface: config, auth, logging, Docker, CI
- Keep the dependency footprint small and auditable
- Make it easy to delete capabilities that a team doesn't need
- Automate dependency updates with grouped Dependabot + auto-merge
- Hard-block dangerous commands via Claude Code hooks with no override mechanism

**Non-Goals:**
- GraphQL or gRPC support (HTTP/REST only)
- Database connectivity (teams bring their own driver)
- Message queue / event streaming integration
- Front-end or SSR templating
- Kubernetes manifests or Helm charts (out of scope for v1)

## Decisions

### D1: Router — Chi over Gin/Echo/Fiber

**Decision:** Use `go-chi/chi` as the HTTP router.

**Rationale:** Chi is 100% `net/http`-compatible, has zero dependencies beyond stdlib, and supports composable middleware and sub-routers cleanly. Gin and Echo have their own `Context` types that bleed into handler signatures, making handlers harder to test and port. Fiber uses `fasthttp`, incompatible with `net/http` middleware.

**Alternative considered:** stdlib `net/http` ServeMux (Go 1.22+) — adequate for simple routing but lacks middleware chaining and named-parameter extraction that most real apps need quickly.

---

### D2: Configuration — Viper + typed struct

**Decision:** Load config via `spf13/viper`, bind to a strongly-typed `Config` struct, validate at startup.

**Rationale:** Viper handles `.env` files, environment variables, and config file fallback in one library. Binding to a typed struct centralises the config schema, enables IDE completion, and surfaces missing values at boot time rather than at call sites.

**Alternative considered:** `joho/godotenv` + `os.Getenv` — simpler but no schema, no validation, scattered access across the codebase.

---

### D3: Auth — JWT middleware + generation helper only

**Decision:** Ship `Authenticate` middleware and `GenerateTokenPair(userID string)` helper using `golang-jwt/jwt/v5`. No login or refresh handlers.

**Rationale:** Without a database, a login handler has no credential source — shipping one is misleading. A refresh endpoint without a token blocklist gives false safety. The kit's job is to protect routes; teams know their own credential source and call `GenerateTokenPair` after verifying it. Access token TTL: 15 min. Refresh token TTL: 7 days (for teams that add their own refresh logic).

---

### D4: Logging — stdlib `slog`

**Decision:** Use Go 1.21+ stdlib `log/slog`. JSON handler in production, text handler in development.

**Rationale:** Zero dependency. Native structured key-value logging. Teams can swap to zerolog/zap if higher throughput logging is needed — the interface is compatible.

---

### D5: Observability — structured logging only

**Decision:** Use Go stdlib `log/slog` for structured logging. No metrics or tracing included.

**Rationale:** Prometheus and OTEL add significant transitive dependencies and impose an observability stack choice on teams. Teams have diverse monitoring setups (Datadog, New Relic, CloudWatch) that conflict with a bundled Prometheus exporter. Logging via `slog` is zero-dependency, universally useful, and integrates with any downstream platform.

---

### D6: Project layout — golang-standards/project-layout

**Decision:** `cmd/`, `internal/`, `api/`, `scripts/`. Single entry point at `cmd/server/main.go`.

**Rationale:** Industry-standard, widely understood. `internal/` prevents accidental external import of application internals.

---

### D7: Dependabot — grouped updates with auto-merge

**Decision:** `.github/dependabot.yml` with weekly grouped updates and auto-merge enabled for patch-level bumps that pass CI.

Groups:
- Go modules: `http` (chi, middleware), `auth` (golang-jwt), `dev` (golangci-lint, govulncheck)
- GitHub Actions: all actions pinned to SHA + tag, updated as a single group

**Rationale:** Grouped updates produce one PR per category per week instead of one PR per package — manageable without being noisy. Auto-merge on patch bumps (CI required) eliminates manual review for low-risk updates while keeping humans in the loop for minor/major bumps.

**Alternative considered:** Individual PRs per package — too noisy on an active repo with many transitive deps.

---

### D8: Claude Code safety — hard-block hooks, no override

**Decision:** Ship `.claude/settings.json` registering a `PreToolUse` hook on the Bash tool, and `.claude/hooks/guard-risky-cmds.sh` that hard-blocks dangerous commands. No override mechanism (`CLAUDE_ALLOW_RISKY` or similar) — blocked commands require the user to run them manually.

**Hook categories and classification:**

| Category | Examples | Action |
|---|---|---|
| Credentials / SSH | `cat ~/.ssh/id_*`, `gpg --export-secret-keys`, `printenv` full dump | Block |
| Pipe-to-shell | `curl \| bash`, `wget \| sh` | Block |
| Sudo | Any `sudo` invocation | Block |
| /etc writes | Any write to `/etc/` | Block |
| Git destructive | `push --force` to main/master, `reset --hard`, `clean -f` | Block |
| Docker destructive | `system prune`, `volume rm`, `rm -f` containers | Block |
| SSH tunnels / netcat | `ssh -L/-R/-D`, `nc -l` | Block |

**Rationale:** Hard-blocking with no override puts the human firmly in control of destructive operations. An override flag (`CLAUDE_ALLOW_RISKY=1`) would be bypassed accidentally or out of habit, defeating the purpose. When Claude blocks, it explains what was attempted and why — the user runs it manually if they intend it.

**Alternative considered:** Prompt-and-confirm model — rejected because it still allows Claude to proceed on user confirmation, which can be given carelessly.

### D9: Git hooks — Lefthook + shell-based commit-msg validator

**Decision:** Use `evilmartians/lefthook` for Git hook management. The `commit-msg` hook calls `scripts/check-commit-msg.sh` which validates Conventional Commits format via shell regex. No Node.js or commitlint required.

**Conventional Commits pattern enforced:**
```
^(feat|fix|chore|docs|refactor|test|ci|build|perf|style|revert)(\(.+\))?: .+
```

**Rationale:** Lefthook is a single Go binary — installable via `go install` or Homebrew, consistent with the Go toolchain already required. A shell regex check covers the most important commitlint rule (type prefix) without pulling in the Node.js ecosystem. Teams needing scope enforcement or max-length rules can extend the shell script.

**Alternative considered:** Husky + commitlint — requires Node.js as a dev dependency in a pure Go project, adding toolchain friction for no significant benefit given the project's scope.

**Lefthook hooks configured:**
- `commit-msg`: runs `scripts/check-commit-msg.sh` to validate message format
- `pre-commit`: runs `golangci-lint` on staged packages

### D10: CI/CD — GHCR, tag-based release, govulncheck

**Decision:** PR checks run lint + test + build + `govulncheck`. Release triggers on `v*` tags (not every push to `main`) and pushes to GHCR using `GITHUB_TOKEN`.

**Rationale:**
- **GHCR over Docker Hub**: `GITHUB_TOKEN` is automatic — no secrets to configure before first use. Teams that fork get a working release pipeline immediately.
- **Tag-based release**: Pushing an image on every merge to `main` produces unintentional releases. A semver tag is a deliberate act; it's the right release trigger for a kit teams are iterating on.
- **govulncheck**: `golang.org/x/vuln` is a zero-dependency stdlib tool. Catching known CVEs in the dependency tree at PR time costs nothing and models good practice.

---

## Risks / Trade-offs

| Risk | Mitigation |
|---|---|
| Chi unfamiliar to Gin/Echo teams | API is similar; README includes notes. Chi docs are comprehensive. |
| JWT has no built-in revocation | Short-lived access (15 min) + refresh (7 days). Blocklist table is a documented extension point. |
| Viper adds ~1MB to binary | Acceptable trade-off. Teams can swap to `env` package if binary size is critical. |
| Lefthook pre-commit runs tests on every commit | Can slow down commits on large repos; teams can tune to lint-only if needed. |
| Dependabot patch auto-merge could ship a breaking patch | CI must pass (tests + lint) before auto-merge. Teams can disable auto-merge per group. |
| Hook script blocks legitimate commands during development | Categories are narrow and explicit. False positives are expected to be rare; false negatives are the priority concern. |
