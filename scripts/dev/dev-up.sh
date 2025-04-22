#!/bin/bash

set -e

ENV_PATH="./docker/dev/.env"
COMPOSE_FILE="./docker/dev/docker-compose.yml"

echo "Docker-compose up"
docker-compose --env-file "$ENV_PATH" -f "$COMPOSE_FILE" up -d postgres redis

echo "Waiting database up..."
until docker exec $(docker-compose -f "$COMPOSE_FILE" ps -q postgres) pg_isready -U postgres; do
  sleep 1
done

echo "Make migration"
./scripts/dev/migrator.sh

echo "Success!"

