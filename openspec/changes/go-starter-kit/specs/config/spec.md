## ADDED Requirements

### Requirement: Typed configuration struct
All application configuration SHALL be loaded into a single strongly-typed `Config` struct at startup using `spf13/viper`.

#### Scenario: Config loaded from environment
- **WHEN** the application starts with environment variables set
- **THEN** viper SHALL bind those variables to the `Config` struct fields automatically

#### Scenario: Config loaded from .env file
- **WHEN** a `.env` file exists in the working directory
- **THEN** viper SHALL load it and merge its values with environment variables, with environment variables taking precedence

#### Scenario: Config struct is accessible
- **WHEN** any package needs a configuration value
- **THEN** it SHALL receive the `Config` struct via dependency injection, not call viper directly

### Requirement: Startup validation
The application SHALL validate required configuration fields at startup and exit with a descriptive error if any are missing or invalid.

#### Scenario: Missing required field
- **WHEN** a required field such as `JWT_SECRET` is not set
- **THEN** the application SHALL log the missing field name and exit with code 1 before binding any port

#### Scenario: Optional fields use safe defaults
- **WHEN** optional fields are not set
- **THEN** the application SHALL apply defaults: `PORT=8080`, `APP_ENV=development`, `MAX_REQUEST_BODY_BYTES=1048576`

#### Scenario: Invalid value
- **WHEN** a configuration field has an invalid value (e.g., non-integer port)
- **THEN** the application SHALL log the field name, the invalid value, and exit with code 1

### Requirement: Environment-specific defaults
The configuration SHALL support a `APP_ENV` field with values `development`, `test`, and `production`, and apply safe defaults per environment.

#### Scenario: Development defaults
- **WHEN** `APP_ENV=development`
- **THEN** the logger SHALL use text format and debug level, and CORS SHALL allow all origins

#### Scenario: Production defaults
- **WHEN** `APP_ENV=production`
- **THEN** the logger SHALL use JSON format and info level, and CORS SHALL restrict to configured allowed origins

### Requirement: Example env file
The repository SHALL include a `.env.example` file listing all supported configuration keys with descriptions and safe placeholder values.

#### Scenario: Example file present
- **WHEN** a developer clones the repository
- **THEN** `.env.example` SHALL exist at the root and contain every key that the `Config` struct reads
