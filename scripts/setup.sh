#!/bin/bash

go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

cp scripts/pre-commit .git/hooks/pre-commit
cp scripts/pre-push .git/hooks/pre-push

chmod +x .git/hooks/pre-commit
chmod +x .git/hooks/pre-push