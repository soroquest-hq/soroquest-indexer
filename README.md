# soroquest-indexer

Event indexer and REST API for the SoroQuest bounty platform.

## Overview

Listens to Soroban contract events, keeps a PostgreSQL database in sync,
and exposes a REST API that the frontend queries for fast bounty reads.

## Quick start

```bash
cp .env.example .env
# Edit .env with your database and RPC URLs

# Run migrations
make migrate-up

# Start the indexer
make run
```

## API docs

See [docs/api.md](docs/api.md).

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md).

## License

MIT
