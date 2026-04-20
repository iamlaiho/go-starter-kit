## ADDED Requirements

### Requirement: JWT token generation helper
The auth package SHALL expose a `GenerateTokenPair(userID string)` helper that returns signed JWT access and refresh tokens using `golang-jwt/jwt/v5` with HS256.

#### Scenario: Token pair generated
- **WHEN** a caller invokes `GenerateTokenPair(userID)`
- **THEN** it SHALL return a signed access token with 15-minute expiry and a refresh token with 7-day expiry

#### Scenario: Token contains standard claims
- **WHEN** a token is generated
- **THEN** it SHALL include `sub` (user ID), `exp` (expiry), `iat` (issued at), and `jti` (unique token ID) claims

### Requirement: JWT authentication middleware
The auth middleware SHALL validate Bearer tokens on protected routes and reject invalid or expired tokens.

#### Scenario: Valid token grants access
- **WHEN** a request includes a valid `Authorization: Bearer <token>` header
- **THEN** the middleware SHALL extract the user ID from claims, attach it to the request context, and call the next handler

#### Scenario: Missing token rejected
- **WHEN** a request to a protected route has no `Authorization` header
- **THEN** the middleware SHALL return HTTP 401 with `{"error": "unauthorized"}`

#### Scenario: Expired token rejected
- **WHEN** a request includes an expired token
- **THEN** the middleware SHALL return HTTP 401 with `{"error": "token expired"}`

#### Scenario: Invalid signature rejected
- **WHEN** a request includes a token with an invalid signature
- **THEN** the middleware SHALL return HTTP 401 with `{"error": "invalid token"}`

### Requirement: Protected route group
The router SHALL expose a protected route group where all routes require a valid Bearer token.

#### Scenario: Protected routes require authentication
- **WHEN** the router is configured
- **THEN** the `Authenticate` middleware SHALL be applied to a sub-router that teams register their protected routes on

### Requirement: JWT secret configuration
The JWT signing secret SHALL be loaded from configuration and never hardcoded.

#### Scenario: Missing secret at startup
- **WHEN** `JWT_SECRET` is not set in configuration
- **THEN** the application SHALL exit with code 1 before binding any port
