#!/usr/bin/env bash

go run ./cmd/migrator/main.go \
  --config-path="${CONFIG_PATH}" \
  --migrations-path="${MIGRATIONS_PATH}"
