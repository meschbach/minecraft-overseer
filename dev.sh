#!/bin/bash

set -e
go fmt ./...
go test ./...
go build -o minecraft-overseer ./cmd/server
docker build . --tag meschbach/minecraft-overseer:dev
docker run --rm -p25565:25565 -v $PWD/test/config:/mc/config:ro -v $PWD/test/secrets:/mc/secrets -v $PWD/test/data:/mc/instance meschbach/minecraft-overseer:dev
