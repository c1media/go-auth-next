# Makefile for Go Auth Template

.PHONY: run migrate migrate-only tidy build

run:
	air -c .air.toml || (echo 'Air not installed. Run: go install github.com/air-verse/air@latest')

migrate:
	go run cmd/server/main.go --migrate

migrate-only:
	go run cmd/server/main.go --migrate-only

tidy:
	go mod tidy

build:
	go build -o bin/authserver cmd/server/main.go 