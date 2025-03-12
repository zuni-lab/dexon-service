# Dexon service

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Configuration](#configuration)
5. [Testing](#testing)
6. [Deployment](#deployment)
7. [Contributing](#contributing)
8. [License](#license)

## Overview

Structure of the project:

```
┣📦 .github                      # Github actions
┃
┣📦 cmd                          # Command line
┃ ┗ 📂 api
┃    ┗ 📂 server                 # Setup routes, middlewares, services,...
┃    ┣ 📂 test                   # Integration tests
┃    ┗ 📜 main.go                # Api entry point
┣ 📦 config                        # Configuration
┃    ┣ 📜 env.go                 # Environment setup
┃    ┗ 📜 logger.go              # Logger setup: zerolog, openobserve, stdout
┣ 📦 constants                    # Constants: errors, mail, ...
┣ 📦 docs                         # Documents
┣ 📦 internal                     # Internal packages
┃    ┣ 📂 chat                    # Chat module (AI bot)
┃    ┣ 📂 orders                  # Order module (limit/stop/twap logic)
┃    ┣ 📂 pools                   # Pool module  (in development)
┣ 📦 pkg                          # Public packages
┃    ┣ 📂 db
┃    ┃   ┣ 📂 migratation        # Database migrations
┃    ┃   ┗ 📂 query              # Database queries
┃    ┣ 📜 *.sql.go               # SQLC generated go file
┃    ┣ 📜 *.go                   # Transactions or additional logic
┃    ┣ 📂 evm                    # EVM based services
┃    ┃ ┣ 📜 *.contract.go        # SOLC generated contract interfaces
┃    ┃ ┗ 📜 real_time_manager.go # Listen swap event to handle matching
┃    ┃ ┗ 📜 tx.go                # Call transaction
┃    ┣ 📂 openai                 # Open AI service, support chat
┃    ┣ 📂 openobserve            # Observability: traces, logs...
┃    ┣ 📂 swap                   # Swap handler (interfaces)
┃    ┗ 📂 utils                  # Utilities
┃
┣ 📜 .air.toml                    # Air configuration
┣ 📜 .env.example                 # Env example
┣ 📜 .gitignore                   # Git ignore
┣ 📜 app.compose.yml              # App docker compose
┣ 📜 compose.yml                  # Docker compose
┣ 📜 Dockerfile                    # Dockerfile
┣ 📜 go.mod                       # Go modules
┣ 📜 go.sum                       # Go modules
┣ 📜 Makefile                      # Makefile
┗ 📜 README.md                    # Readme
```

## Installation

1. Clone the repository
2. Install dependencies

- `go mod tidy`
- Install Air for hot reload: `go install github.com/cosmtrek/air@latest`

3. Set up configuration (if any)

- Copy `.env.example` to `.env` and update the values

## Usage

- Run bootstrapping: `make compose` to start the services
  1. This will start the Openobserve service that runs on port 5080
     - You can access the Openobserve dashboard at `http://localhost:5080`
     - Login with the default email and password in `compose.yml`, or update the values
     - Access the [Ingestion API - Trace](http://localhost:5080/web/ingestion/custom/traces/opentelemetry)
     - Copy the `Authorization` header value and update the `OPENOBSERVE_CREDENTIAL` in the `.env` file
     - Access the [Trace Tab](http://localhost:5080/web/traces) to view the traces
  2. This will start the Postgres service that runs on port 5432
- Run the application: `make run`
- Run the application with hot reload: `make watch`
- Run application with Docker: `make start`
- Stop the docker container: `make stop`
- Shutdown and clean up: `make shutdown`
- Restart the application: `make restart`

## Configuration

- Update the configuration in the `.env` file
- Update the logger configuration in `config/logger.go`
- Update the environment setup in `config/env.go`
- Update the database connection in `pkg/db/init.go`
- Update the database models in `pkg/db/migration`
- Add new routes in `cmd/api/server/routes.go`
- Add new services in `internal/<module>/services`
- Add new handlers in `internal/<module>/handlers`
- Add new middlewares in `internal/middleware`
- Add new constants in `constants`
- If you want to add new packages, add them in `pkg` such as `pkg/<package>`: jwt, cookie, mail, cache, oauth, s3, ...

## Testing

- Run tests: `make test`
- Run tests with coverage: `make coverage`

## Deployment

- With existing Dockerfile, you can deploy the application to any cloud provider
- Update the Dockerfile if needed

## Contributing

- Fork the repository
- Create a new branch
- Make your changes
- Create a pull request

## License
