SHELL ?= /bin/bash

build:
	@go build -o ./bin/go-movies ./cmd/api

build-prod:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/go-movies ./cmd/api

run: build
	@./bin/go-movies

tidy:
	@go mod tidy
