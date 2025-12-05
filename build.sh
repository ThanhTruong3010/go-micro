#!/bin/bash

FRONT_END=frontApp
BROKER_APP=brokerApp
AUTH_APP=authApp

up() {
    echo "Starting Docker images..."
    cd ../project && docker compose up -d
    echo "Docker images started!"
}

down() {
    echo "Stopping Docker images..."
    cd ../project && docker compose down
    echo "Done!"
}

build_broker() {
    echo "Buidling $BROKER_APP..."
    cd ../broker-service && env CGO_ENABLED=0 go build -o $BROKER_APP ./cmd/api
    echo "Done!"
}

build_auth() {
    echo "Buidling $AUTH_APP..."
    cd ../authentication-service && env CGO_ENABLED=0 go build -o $AUTH_APP ./cmd/api
    echo "Done!"
}

build_front() {
    echo "Buidling $FRONT_END..."
    cd ../frontend && env CGO_ENABLED=0 go build -o $FRONT_END ./cmd/web
    echo "Done!"
}