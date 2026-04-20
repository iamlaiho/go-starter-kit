## ADDED Requirements

### Requirement: Chi router with middleware stack
The HTTP server SHALL use `go-chi/chi` as the router with a default middleware stack applied to all routes.

#### Scenario: Default middleware applied
- **WHEN** any HTTP request is received
- **THEN** the server SHALL apply timeout, body-limit, security-headers, request-id, structured logger, recoverer, and CORS middleware in that order

#### Scenario: Panic recovery
- **WHEN** a handler panics
- **THEN** the recoverer middleware SHALL catch the panic, log it with the request-id, and return HTTP 500 without crashing the server

#### Scenario: Request ID propagation
- **WHEN** a request arrives without an `X-Request-ID` header
- **THEN** the middleware SHALL generate a UUID and attach it to the request context and response header

### Requirement: Structured JSON error responses
All HTTP error responses SHALL be returned as JSON with a consistent structure.

#### Scenario: Not found error
- **WHEN** a route is not registered for the requested path
- **THEN** the server SHALL return HTTP 404 with body `{"error": "not found", "request_id": "<id>"}`

#### Scenario: Validation error
- **WHEN** a handler receives invalid input
- **THEN** the server SHALL return HTTP 400 with body `{"error": "<message>", "request_id": "<id>"}`

#### Scenario: Internal server error
- **WHEN** an unexpected error occurs in a handler
- **THEN** the server SHALL return HTTP 500 with body `{"error": "internal server error", "request_id": "<id>"}` without leaking internal details

### Requirement: Graceful shutdown
The server SHALL support graceful shutdown on SIGINT and SIGTERM signals.

#### Scenario: In-flight requests complete
- **WHEN** the server receives SIGINT or SIGTERM
- **THEN** it SHALL stop accepting new connections and allow in-flight requests up to 30 seconds to complete before exiting

### Requirement: Request timeout middleware
The server SHALL enforce a per-request timeout to prevent slow clients from holding goroutines indefinitely.

#### Scenario: Request exceeds timeout
- **WHEN** a request handler does not complete within 30 seconds
- **THEN** the server SHALL cancel the request context and return HTTP 503

#### Scenario: Timeout applied to all routes
- **WHEN** the server starts
- **THEN** the timeout middleware SHALL be applied globally in the middleware stack before handlers execute

### Requirement: Request body size limit
The server SHALL reject request bodies exceeding a configured maximum size.

#### Scenario: Oversized body rejected
- **WHEN** a request body exceeds 1MB (default)
- **THEN** the server SHALL return HTTP 413 with `{"error": "request body too large", "request_id": "<id>"}`

#### Scenario: Size limit configurable
- **WHEN** `MAX_REQUEST_BODY_BYTES` is set in configuration
- **THEN** the server SHALL apply that value as the body size limit

### Requirement: Input validation
The server SHALL provide a shared validation helper using `go-playground/validator` that handlers use to validate decoded request bodies.

#### Scenario: Valid request body accepted
- **WHEN** a handler decodes and validates a request body that passes all validation rules
- **THEN** the handler SHALL proceed to business logic

#### Scenario: Invalid request body rejected
- **WHEN** a handler validates a request body that fails one or more rules
- **THEN** the server SHALL return HTTP 400 with `{"error": "<field>: <message>", "request_id": "<id>"}` for each failing field

#### Scenario: Validation struct tags used
- **WHEN** a developer defines a request struct
- **THEN** they SHALL use `validate:"required"`, `validate:"email"`, `validate:"min=N"` tags and call the shared `Validate(v any) error` helper

### Requirement: Health check endpoint
The server SHALL expose a `/health` endpoint for deployment readiness and liveness probes.

#### Scenario: Health endpoint returns 200
- **WHEN** a GET request is made to `/health`
- **THEN** the server SHALL return HTTP 200 with body `{"status": "ok", "version": "<git-sha>"}`

#### Scenario: Health endpoint is unauthenticated
- **WHEN** a GET request is made to `/health` without an auth token
- **THEN** the server SHALL respond without requiring authentication

### Requirement: Security headers middleware
The server SHALL apply HTTP security headers to all responses.

#### Scenario: Security headers present on every response
- **WHEN** any HTTP response is sent
- **THEN** the server SHALL include `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`, and `Content-Security-Policy: default-src 'self'`

