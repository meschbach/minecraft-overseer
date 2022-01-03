#!/bin/bash

set -e
go fmt ./...
docker build . --tag meschbach/minecraft-overseer:dev
docker run --rm -v $PWD/test/config:/mc/config:ro -v $PWD/test/data:/mc/instance meschbach/minecraft-overseer:dev
