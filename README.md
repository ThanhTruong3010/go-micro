# Go Microservices

A microservices architecture project built with Go, featuring a broker service, authentication service, logger service, and front-end application.

## Architecture

```
┌─────────────────┐     ┌─────────────────────┐     ┌──────────────┐     ┌──────────────┐
│   Front-End     │────▶│   Broker Service    │────▶│Logger Service│────▶│   MongoDB    │
│   (Port 80)     │     │    (Port 8080)      │    │  (Port 8082) │     │ (Port 27017) │
└────────┬────────┘     └──────────┬──────────┘     └──────────────┘     └──────────────┘
         │                         │
         │              ┌──────────┴──────┐
         └─────────────▶│  Auth Service   │
                        │  (Port 8081)    │
                        └────────┬────────┘
                                 │
                                 ▼
                        ┌──────────────┐
                        │  PostgreSQL  │
                        │  (Port 5432) │
                        └──────────────┘
```

## Project Structure

```
go-micro/
├── authentication-service/    # User authentication microservice
│   ├── cmd/api/              # Application entry point
│   └── .docker/              # Docker configurations
├── broker-service/           # API gateway/broker microservice
│   ├── cmd/api/              # Application entry point
│   └── .docker/              # Docker configurations
├── logger-service/           # Logging microservice
│   ├── cmd/api/              # Application entry point
│   └── .docker/              # Docker configurations
├── front-end/                # Web front-end application
│   └── cmd/web/              # Application entry point
├── project/                  # Docker compose configurations
│   ├── docker-compose.yml           # Production compose file
│   ├── docker-compose-local.yml     # Development compose file
│   └── Makefile                     # Make commands
├── .env                      # Environment variables
└── build.sh                  # Build script
```

## Services

| Service                | Port  | Description                                                   |
| ---------------------- | ----- | ------------------------------------------------------------- |
| Broker Service         | 8080  | API gateway that routes requests to appropriate microservices |
| Authentication Service | 8081  | Handles user authentication                                   |
| Logger Service         | 8082  | Handles logging to MongoDB                                    |
| PostgreSQL             | 5432  | Database for user data                                        |
| MongoDB                | 27017 | Database for logs                                             |
| Front-End              | 80    | Web interface                                                 |

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional)

## Quick Start

### Using build.sh (Recommended)

```bash
# Make the script executable (first time only)
chmod +x build.sh

# Build and start all services (development mode)
./build.sh up_build

# Start services without rebuilding
./build.sh up

# Stop all services
./build.sh down
```

### Using Make

```bash
cd project

# Build and start all services
make up_build

# Start services without rebuilding
make up

# Stop all services
make down
```

## Environment Modes

The project supports two modes controlled by the `MODE` variable in `.env`:

| Mode          | Compose File               | Description               |
| ------------- | -------------------------- | ------------------------- |
| `development` | `docker-compose-local.yml` | Uses pre-built binaries   |
| `production`  | `docker-compose.yml`       | Multi-stage Docker builds |

### Switching Modes

Edit `.env` file:

```bash
# Development mode (default)
MODE=development

# Production mode
MODE=production
```

Or pass as argument:

```bash
./build.sh up_build MODE=production
```

## Available Commands

### build.sh Commands

| Command        | Description                         |
| -------------- | ----------------------------------- |
| `up`           | Start containers without rebuilding |
| `up_build`     | Build and start all containers      |
| `down`         | Stop all containers                 |
| `build_broker` | Build broker binary                 |
| `build_auth`   | Build authentication binary         |
| `build_front`  | Build front-end binary              |
| `start`        | Start front-end application         |
| `stop`         | Stop front-end application          |

### Make Commands

| Command             | Description                         |
| ------------------- | ----------------------------------- |
| `make up`           | Start containers without rebuilding |
| `make up_build`     | Build and start all containers      |
| `make down`         | Stop all containers                 |
| `make build_broker` | Build broker binary                 |
| `make build_auth`   | Build authentication binary         |
| `make build_front`  | Build front-end binary              |
| `make start`        | Start front-end application         |
| `make stop`         | Stop front-end application          |

## API Endpoints

### Broker Service (Port 8080)

| Method | Endpoint  | Description                          |
| ------ | --------- | ------------------------------------ |
| POST   | `/`       | Health check - returns broker status |
| POST   | `/handle` | Handle submission requests           |
| GET    | `/ping`   | Heartbeat endpoint                   |

### Authentication Service (Port 8081)

| Method | Endpoint        | Description                   |
| ------ | --------------- | ----------------------------- |
| POST   | `/authenticate` | Authenticate user credentials |

### Logger Service (Port 8082)

| Method | Endpoint | Description     |
| ------ | -------- | --------------- |
| POST   | `/log`   | Write log entry |

## Example Requests

### Authenticate User

```bash
curl -X POST http://localhost:8080/handle \
  -H "Content-Type: application/json" \
  -d '{
    "action": "auth",
    "auth": {
      "email": "admin@example.com",
      "password": "verysecret"
    }
  }'
```

## Databases

### PostgreSQL (User Data)

| Property | Value                                     |
| -------- | ----------------------------------------- |
| Host     | `localhost` (or `postgres` within Docker) |
| Port     | `5432`                                    |
| User     | `postgres`                                |
| Password | `password`                                |
| Database | `users`                                   |

### MongoDB (Logs)

| Property | Value                                  |
| -------- | -------------------------------------- |
| Host     | `localhost` (or `mongo` within Docker) |
| Port     | `27017`                                |
| User     | `admin`                                |
| Password | `password`                             |
| Database | `logs`                                 |

## Development

### Building Individual Services

```bash
# Build broker service
cd broker-service
go build -o brokerApp ./cmd/api

# Build authentication service
cd authentication-service
go build -o authApp ./cmd/api

# Build logger service
cd logger-service
go build -o loggerApp ./cmd/api

# Build front-end
cd front-end
go build -o frontApp ./cmd/web
```

### Running Locally (without Docker)

1. Start PostgreSQL
2. Set environment variables
3. Run each service:

```bash
# Terminal 1 - Auth Service
cd authentication-service && ./authApp

# Terminal 2 - Broker Service
cd broker-service && ./brokerApp

# Terminal 3 - Logger Service
cd logger-service && ./loggerApp

# Terminal 4 - Front-end
cd front-end && ./frontApp
```

## License

Thanh Truong
