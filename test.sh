#!/bin/bash

PKGS=$(go list ./... | grep -v "cmd/server" | grep -v "mocks")

go test $PKGS -coverprofile=coverage.out

COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
COVERAGE=${COVERAGE%.*}

if [ $? -eq 0 ] && [ "$COVERAGE" -ge 90 ]; then
    echo "Tests passed with ${COVERAGE}% coverage"
    rm coverage.out
    exit 0
else
    echo "Tests failed or coverage below 90% (current: ${COVERAGE}%)"
    exit 1
fi
