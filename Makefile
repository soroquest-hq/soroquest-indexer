.PHONY: build run test migrate-up migrate-down fmt lint

build:
	go build -o soroquest-indexer ./cmd/indexer

run:
	go run ./cmd/indexer

test:
	go test ./...

fmt:
	gofmt -w .

lint:
	go vet ./...

migrate-up:
	golang-migrate -path migrations -database "$$DATABASE_URL" up

migrate-down:
	golang-migrate -path migrations -database "$$DATABASE_URL" down
