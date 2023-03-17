#!/bin/bash

set -e
go fmt ./...
go test ./internal/...
go build -o minecraft-overseer ./cmd/server
docker build . --tag meschbach/minecraft-overseer:dev
docker run --rm -p25565:25565 \
  -e "INSTANCE_NAME=docker-dev" \
  -e "PORT_SPEC=:25565" \
  -v $PWD/test/config:/mc/config:ro -v $PWD/test/data:/mc/instance \
  -v $PWD/test/secrets:/mc/secrets:ro \
  meschbach/minecraft-overseer:dev
