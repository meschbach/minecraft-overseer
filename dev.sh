#!/bin/bash

set -e
go fmt ./...
go build -o minecraft-overseer main.go
docker build . --tag meschbach/minecraft-overseer:dev
docker run --rm -p25565:25565 -v $PWD/test/config:/mc/config:ro -v $PWD/test/data:/mc/instance meschbach/minecraft-overseer:dev
