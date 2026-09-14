# ARCHITECTURE ESSENTIALS — soroquest-indexer

## Stack
- Go 1.22+ | PostgreSQL | `chi` router | `pgx` driver | `golang-migrate`
- Stellar: `github.com/stellar/go` for RPC + XDR decoding
- Deploy: Render (free PostgreSQL + web service = real public URL)

## Key directories
```
cmd/indexer/main.go       → entry point
internal/ingest/          → ingestion loop + event handlers
internal/db/              → all database queries
internal/api/             → HTTP handlers
internal/stellar/         → RPC client wrapper
migrations/               → SQL up/down files
```

## Database tables
```sql
bounties  (id, owner, title, description, amount, token, status,
           claimant, created_at, claim_deadline, updated_at)

events    (id, bounty_id, event_type, ledger, tx_hash, payload JSONB, created_at)

cursor    (id=1, last_ledger, updated_at)   -- always exactly one row
```

## Ingestion loop — critical pattern
```
read cursor.last_ledger
  → fetch contract events from (last_ledger+1) to latest
    → for each event:
        BEGIN TRANSACTION
          dispatch to handler (posted|claimed|completed|cancelled)
          update bounties table
          insert into events table
          update cursor.last_ledger
        COMMIT
  → sleep 5s
  → on RPC error: log, backoff, reconnect, resume from cursor
```
Everything inside one DB transaction. If anything fails, cursor does
not advance and event is retried. No silent skips.

## Event handlers
```go
handlePosted(ctx, db, event)    // INSERT bounty, status=open
handleClaimed(ctx, db, event)   // UPDATE status=claimed, claimant=address
handleCompleted(ctx, db, event) // UPDATE status=completed
handleCancelled(ctx, db, event) // UPDATE status=cancelled, claimant=null
```

## API endpoints
```
GET /api/v1/bounties             ?status=&limit=&offset=
GET /api/v1/bounties/:id
GET /api/v1/bounties/:id/events
GET /api/v1/stats
GET /health                      → { last_ledger, lag }
```

## Amount handling — critical
- Store as PostgreSQL `NUMERIC`, return as decimal string in JSON
- i128 can exceed JS Number.MAX_SAFE_INTEGER — never use float
- Frontend uses BigInt to read it

## Environment variables
```
DATABASE_URL
STELLAR_RPC_URL
CONTRACT_ID
NETWORK_PASSPHRASE
PORT
START_LEDGER   (0 = from contract deployment)
```

## Critical rules
- One DB transaction per event — cursor never advances on failure
- Never write to the contract — indexer is read-only against the chain
- Never drop events silently — log and retry on any processing error
- CORS must allow the app's Vercel domain
- `/health` must expose `last_ledger` so app can detect lag
- Amount is always a string (NUMERIC in DB, string in JSON, BigInt in frontend)
