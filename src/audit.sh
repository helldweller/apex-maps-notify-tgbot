#!/bin/bash

# Verify dependencies
go mod verify

# Build
go build -v ./...

# Run go vet
go vet ./...

# Install staticcheck
go install honnef.co/go/tools/cmd/staticcheck@latest

# Run staticcheck
staticcheck ./...

# Install golint
go install golang.org/x/lint/golint@latest

# Run golint
golint -set_exit_status -min_confidence 0.5 ./...

# Run tests
go test -race -vet=off ./...
