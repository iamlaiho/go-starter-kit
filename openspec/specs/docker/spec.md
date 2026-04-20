## Purpose

Defines the Docker setup: multi-stage production Dockerfile, docker-compose for local development, and Makefile targets for common Docker operations.

## Requirements

### Requirement: Multi-stage Dockerfile
The repository SHALL include a multi-stage `Dockerfile` producing a minimal production image.

#### Scenario: Builder stage compiles binary
- **WHEN** the Docker build runs the builder stage
- **THEN** it SHALL use a `golang:1.22-alpine` base image, copy source, and produce a statically linked binary

#### Scenario: Runtime stage is minimal
- **WHEN** the Docker build runs the runtime stage
- **THEN** it SHALL use a `gcr.io/distroless/static` base image containing only the binary

#### Scenario: Image runs as non-root
- **WHEN** the container starts
- **THEN** the process SHALL run as a non-root user defined in the Dockerfile

### Requirement: docker-compose for local development
The repository SHALL include a `docker-compose.yml` for running the application locally.

#### Scenario: App starts with environment from .env
- **WHEN** a developer runs `docker compose up`
- **THEN** the application container SHALL start and read environment variables from the `.env` file in the project root

#### Scenario: Health check configured in docker-compose
- **WHEN** the container is running
- **THEN** Docker SHALL poll `GET /health` every 30 seconds with a 5-second timeout and 3 retries before marking the container unhealthy

### Requirement: Makefile Docker targets
The Makefile SHALL include targets for common Docker operations.

#### Scenario: Build target exists
- **WHEN** a developer runs `make docker-build`
- **THEN** the Makefile SHALL build the Docker image tagged with the current git SHA

#### Scenario: Run target exists
- **WHEN** a developer runs `make docker-run`
- **THEN** the Makefile SHALL start the stack using docker-compose
