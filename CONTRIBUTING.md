# Contributing

## Setup

1. Install Go 1.22+
2. Install PostgreSQL
3. Copy `.env.example` to `.env` and fill in values
4. Run `make migrate-up` to create the schema
5. Run `make run` to start the indexer

## Testing

Run `make test`. Tests require a running PostgreSQL instance.

## Code style

Run `make fmt` and `make lint` before submitting a PR.
