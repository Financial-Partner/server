#!/bin/bash
set -e

if ! command -v swag &> /dev/null; then
    echo "swag could not be found, installing it..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

echo "Generating Swagger documentation..."
swag init -g cmd/server/main.go -o ./swagger

echo "Swagger documentation generated successfully." 