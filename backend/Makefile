.PHONY: run test tidy migrate-up migrate-down build

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api

migrate-up:
	go run ./cmd/migrate -dir up

migrate-down:
	go run ./cmd/migrate -dir down
