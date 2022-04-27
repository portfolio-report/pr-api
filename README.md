# Portfolio Report API

[![Go version](https://img.shields.io/github/go-mod/go-version/portfolio-report/pr-api?style=flat-square)](https://go.dev/)
[![Build status](https://img.shields.io/github/workflow/status/portfolio-report/pr-api/CI?style=flat-square)](https://github.com/portfolio-report/pr-api/actions/workflows/ci.yml)
[![Go report card](https://goreportcard.com/badge/github.com/portfolio-report/pr-api?style=flat-square)](https://goreportcard.com/report/github.com/portfolio-report/pr-api)
[![Test coverage](https://img.shields.io/codecov/c/github/portfolio-report/pr-api?style=flat-square)](https://app.codecov.io/gh/portfolio-report/pr-api)

This repo holds the API of [Portfolio Report](https://www.portfolio-report.net/). It provides [Portfolio Performance](https://www.portfolio-performance.info/) with master data and historical quotes of securities. In addition it provides beta functionality to synchronize portfolios across devices and make them accessible via different UIs.

Please use the [discussion forum](https://forum.portfolio-performance.info/) to get help and discuss improvements.

## Getting started for development

### Preparation
- Install [Go](https://go.dev/)
- Get access to a [PostgreSQL](https://www.postgresql.org/) database
  - [Install it locally](https://www.postgresql.org/download/)
  - Run [docker image](https://hub.docker.com/_/postgres)
  - Use a [cloud service](https://www.postgresql.org/support/professional_hosting/)
- Clone this repo

### Run
```bash
# Set `DATABASE_URL` environment variable or `.env`
$ DATABASE_URL="postgresql://user:password@host:5432/database?sslmode=disable"

# Start
$ go run .
```

The backend provides a SwaggerUI for the REST API on `/doc` and a GraphQL playground on `/graphql`.

## Development hints

### Change GraphQL schema

After changing the graphql schema, make sure to run:
```bash
$ go generate ./...
```

### Execute tests

```bash
$ go test ./...
```

## Configuration

Configuration is done via environment variables or in the `.env` file.

### Mandatory parameters

Application will not start if mandatory parameters are missing.

```ini
# PostgreSQL database URL
DATABASE_URL="postgresql://user:password@host:5432/database?sslmode=disable"
```

### Recommended parameters
Missing recommended parameters will not prevent the application to start, but can lead to limited functionality.

```ini
# E-mail address used as recipient in contact endpoint
CONTACT_RECIPIENT_EMAIL="me@example.com"

# Mail server URL
MAILER_TRANSPORT="smtp://username:password@smtp.example.com:587/"

# Token to download GeoIP database from www.ip2location.com
IP2LOCATION_TOKEN="..."
```

### Optional parameters
Optional parameters will use default value if not set.

```ini
# Mode of gin gonic (debug, test or release)
GIN_MODE="release"

# Maximum number of open database connections
DATABASE_MAX_OPEN_CONN=25

# Maximum number of idle database connections
DATABASE_MAX_IDLE_CONN=25

# Maximum lifetime (in seconds) of database connections
DATABASE_CONN_MAX_LIFE=300

# Allowed period of inactivity for sessions in seconds
SESSION_TIMEOUT=900

# Maximum number of search results
SECURITIES_SEARCH_MAX_RESULTS=10
```
