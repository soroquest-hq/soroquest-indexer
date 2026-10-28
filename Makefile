.PHONY: all build run test fmt lint docker-up docker-down clean verify help

all: fmt lint test build

build:
	go build -o bin/indexer ./cmd/indexer
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/indexer

test:
	go test -v -race ./...

fmt:
	gofmt -w -s .

lint:
	go vet ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

clean:
	rm -rf bin/

verify: fmt lint test
	@echo "All indexer verification checks passed successfully."

help:
	@echo "SoroQuest Indexer Development Commands:"
	@echo "  make build       - Compile indexer and api binaries into bin/"
	@echo "  make test        - Run unit test suite with race detector"
	@echo "  make lint        - Run go vet"
	@echo "  make docker-up   - Launch postgres and indexer via docker compose"
	@echo "  make verify      - Run fmt, lint, and test checks"
