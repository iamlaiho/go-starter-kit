VERSION := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/iamlaiho/go-starter-kit/internal/handler.Version=$(VERSION)"
BINARY := bin/server

.PHONY: run build test test-coverage tidy vuln lint docker-build docker-run hooks-install dev

run:
	go run $(LDFLAGS) ./cmd/server

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/server

test:
	go test ./...

test-coverage:
	mkdir -p coverage
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/index.html

tidy:
	go mod tidy

vuln:
	govulncheck ./...

lint:
	golangci-lint run

docker-build:
	docker build -t go-starter-kit:$(VERSION) .

docker-run:
	docker compose up

hooks-install:
	go install github.com/evilmartians/lefthook@latest && lefthook install

dev:
	air
