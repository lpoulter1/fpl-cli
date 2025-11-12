SHELL := /bin/bash

.PHONY: build test lint

build:
	go build ./cmd/fplcli

test:
	go test ./...

lint:
	golangci-lint run

