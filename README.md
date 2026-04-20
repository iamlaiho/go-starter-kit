# go-starter-kit

Production-ready Go API server scaffold. Fork and build on top.

## Prerequisites

- Go 1.22+
- Docker
- [Air](https://github.com/air-verse/air) for live reload: `go install github.com/air-verse/air@latest`
- [golangci-lint](https://golangci-lint.run) for linting: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) for vulnerability scanning: `go install golang.org/x/vuln/cmd/govulncheck@latest`

## Quickstart

```bash
cp .env.example .env
# edit .env and set JWT_SECRET
make dev
```

## Make Targets

| Target | Description |
|---|---|
| `make run` | Run server with version injected |
| `make build` | Build binary to `bin/server` |
| `make dev` | Run with Air live reload |
| `make test` | Run all tests |
| `make test-coverage` | Run tests and generate HTML coverage report |
| `make tidy` | Run `go mod tidy` |
| `make vuln` | Run `govulncheck` |
| `make lint` | Run `golangci-lint` |
| `make docker-build` | Build Docker image tagged with git SHA |
| `make docker-run` | Start stack with `docker compose up` |
| `make hooks-install` | Install Lefthook git hooks |

## Environment Variables

See `.env.example` for all configuration keys with descriptions.

## Running Tests

```bash
make test
make test-coverage   # generates coverage/index.html
```

## Git Hooks

Run `make hooks-install` after cloning to activate commit message validation and pre-commit linting.
