#!/bin/bash

# Air (hot reloading in development)
go install github.com/air-verse/air@latest

# Mockgen (mock generation)
go install go.uber.org/mock/mockgen@latest

# Wire (dependency injection)
go install github.com/google/wire/cmd/wire@latest
