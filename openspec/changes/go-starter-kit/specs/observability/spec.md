## ADDED Requirements

### Requirement: Structured logging with slog
The application SHALL use Go stdlib `log/slog` for all structured logging with format determined by environment.

#### Scenario: JSON logging in production
- **WHEN** `APP_ENV=production`
- **THEN** all log output SHALL be JSON-formatted with fields: `time`, `level`, `msg`, and any additional key-value pairs

#### Scenario: Text logging in development
- **WHEN** `APP_ENV=development`
- **THEN** all log output SHALL be human-readable text format

#### Scenario: Request logging
- **WHEN** any HTTP request completes
- **THEN** the logger middleware SHALL emit a log entry with `method`, `path`, `status`, `duration`, and `request_id` fields
