#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/../server"

echo "Generating Swagger docs..."
swag init \
  --generalInfo main.go \
  --output docs \
  --outputTypes go,json,yaml \
  --parseDependency \
  --parseInternal

echo "Swagger docs generated in server/docs/"
