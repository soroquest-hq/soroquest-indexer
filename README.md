# SoroQuest Indexer

[![CI](https://github.com/soroquest-hq/soroquest-indexer/actions/workflows/ci.yml/badge.svg)](https://github.com/soroquest-hq/soroquest-indexer/actions/workflows/ci.yml)

The SoroQuest Indexer listens to Soroban smart contract events and indexes them into a PostgreSQL database, providing a fast REST API for the frontend.

## Prerequisites

- Go 1.22+
- PostgreSQL database

## Setup

1. Copy `.env.example` to `.env` and fill in your details:
   ```
   DATABASE_URL=postgres://user:pass@localhost:5432/soroquest
   PORT=8080
   STELLAR_RPC_URL=https://soroban-testnet.stellar.org
   CONTRACT_ID=<deployed_contract_address>
   STELLAR_NETWORK_PASSPHRASE="Test SDF Network ; September 2015"
   ```
2. Run database migrations in the `migrations/` folder against your database.

## Running

```bash
go run cmd/indexer/main.go
```

This starts both the background ingestor loop and the HTTP API on the specified port.

## API Endpoints

- `GET /health` - Service health and sync status
- `GET /api/bounties` - List bounties
- `GET /api/bounties/{id}` - Get a specific bounty
- `GET /api/bounties/{id}/events` - Get event history for a bounty
- `GET /api/stats` - Platform statistics
