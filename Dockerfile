# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X github.com/iamlaiho/go-starter-kit/internal/handler.Version=${VERSION} -w -s" \
    -o server ./cmd/server

# Final stage
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/server /server

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/server"]
