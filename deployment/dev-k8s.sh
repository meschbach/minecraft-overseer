#!/bin/bash

set -e

helm upgrade --atomic --timeout 30s --install test ./minecraft-overseer -f dev-k8s.yaml
