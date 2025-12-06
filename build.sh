#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    source .env
fi

# Variables
FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp
MAIL_BINARY=mailApp
MODE=${MODE:-development}

# Set compose file based on MODE
if [ "$MODE" = "development" ]; then
    COMPOSE_FILE=docker-compose-local.yml
    echo "[MODE] Development mode selected - using $COMPOSE_FILE"
else
    COMPOSE_FILE=docker-compose.yml
    echo "[MODE] Production mode selected - using $COMPOSE_FILE"
fi

# Change to project directory
cd "$(dirname "$0")/project" || exit 1

## up: starts all containers in the background without forcing build
up() {
    echo "Starting Docker images using $COMPOSE_FILE..."
    docker-compose -f $COMPOSE_FILE up -d
    echo "Docker images started!"
}

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build() {
    build_broker
    build_logger
    build_auth
    build_mail
    echo "Stopping docker images (if running...)"
    docker-compose -f $COMPOSE_FILE down
    echo "Building (when required) and starting docker images using $COMPOSE_FILE..."
    docker-compose -f $COMPOSE_FILE up --build -d
    echo "Docker images built and started!"
}

## down: stop docker compose
down() {
    echo "Stopping docker compose..."
    docker-compose -f $COMPOSE_FILE down
    echo "Done!"
}

## build_broker: builds the broker binary as a linux executable
build_broker() {
    echo "Building broker binary..."
    cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o $BROKER_BINARY ./cmd/api
    cd - > /dev/null
    echo "Done!"
}

## build_auth: builds the auth binary as a linux executable
build_auth() {
    echo "Building auth binary..."
    cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o $AUTH_BINARY ./cmd/api
    cd - > /dev/null
    echo "Done!"
}

## build_logger: builds the logger binary as a linux executable
build_logger() {
    echo "Building logger binary..."
    cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o $LOGGER_BINARY ./cmd/api
    cd - > /dev/null
    echo "Done!"
}

## build_mail: builds the mail binary as a linux executable
build_mail() {
    echo "Building mail binary..."
    cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o $MAIL_BINARY ./cmd/api
    cd - > /dev/null
    echo "Done!"
}

## build_front: builds the front end binary
build_front() {
    echo "Building front end binary..."
    cd ../front-end && env CGO_ENABLED=0 go build -o $FRONT_END_BINARY ./cmd/web
    cd - > /dev/null
    echo "Done!"
}

## start: starts the front end
start() {
    build_front
    echo "Starting front end"
    cd ../front-end && ./$FRONT_END_BINARY &
}

## stop: stop the front end
stop() {
    echo "Stopping front end..."
    pkill -SIGTERM -f "./$FRONT_END_BINARY" || true
    echo "Stopped front end!"
}

# Show usage
usage() {
    echo "Usage: $0 [command] [MODE=development|production]"
    echo ""
    echo "Commands:"
    echo "  up          - Start containers without rebuilding"
    echo "  up_build    - Build and start containers"
    echo "  down        - Stop containers"
    echo "  build_broker - Build broker binary"
    echo "  build_auth  - Build auth binary"
    echo "  build_logger  - Build logger binary"
    echo "  build_front - Build front end binary"
    echo "  start       - Start front end"
    echo "  stop        - Stop front end"
    echo ""
    echo "Examples:"
    echo "  $0 up_build                    # Development mode (default)"
    echo "  $0 up_build MODE=production    # Production mode"
    echo "  MODE=production $0 up_build    # Production mode (alternative)"
}

# Parse MODE from arguments (e.g., MODE=production)
# for arg in "$@"; do
#     case $arg in
#         MODE=*)
#             MODE="${arg#*=}"
#             # Re-evaluate compose file
#             if [ "$MODE" = "development" ]; then
#                 COMPOSE_FILE=docker-compose-local.yml
#                 echo "[MODE] x Development mode selected - using $COMPOSE_FILE"
#             else
#                 COMPOSE_FILE=docker-compose.yml
#                 echo "[MODE] x Production mode selected - using $COMPOSE_FILE"
#             fi
#             ;;
#     esac
# done

# Run command
case $1 in
    up)
        up
        ;;
    up_build)
        up_build
        ;;
    down)
        down
        ;;
    build_broker)
        build_broker
        ;;
    build_logger)
        build_logger
        ;;
    build_auth)
        build_auth
        ;;
    build_front)
        build_front
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    *)
        usage
        ;;
esac