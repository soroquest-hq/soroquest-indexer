# ARCHITECTURE — soroquest-indexer

## Tech stack

| Layer | Choice | Reason |
|---|---|---|
| Language | Go 1.22+ | Strong concurrency, good RPC client libraries, fits org stack |
| Database | PostgreSQL | Reliable, queryable, good Go support via `pgx` |
| DB migrations | `golang-migrate` | Simple, file-based, reversible |
| HTTP server | `net/http` + `chi` router | Lightweight, idiomatic Go |
| Stellar RPC | `github.com/stellar/go` | Official Go SDK, handles XDR decoding |
| Config | `env` vars + `.env` file | Simple, works on Render/Railway/Fly.io |
| Deployment | Render (free tier) | Free PostgreSQL + web service, public URL |

## Project structure

```
soroquest-indexer/
├── cmd/
│   └── indexer/
│       └── main.go             # Entry point — wires everything together
├── internal/
│   ├── ingest/
│   │   ├── ingest.go           # Main ingestion loop
│   │   ├── cursor.go           # Last-processed ledger tracking
│   │   └── handlers.go         # One handler per event type
│   ├── db/
│   │   ├── db.go               # Connection pool setup
│   │   ├── bounty.go           # Bounty read/write queries
│   │   └── event.go            # Event log queries
│   ├── api/
│   │   ├── server.go           # HTTP server setup
│   │   ├── bounties.go         # Bounty handlers
│   │   └── health.go           # Health check handler
│   └── stellar/
│       └── client.go           # Soroban RPC client wrapper
├── migrations/
│   ├── 001_create_bounties.up.sql
│   ├── 001_create_bounties.down.sql
│   ├── 002_create_events.up.sql
│   └── 002_create_events.down.sql
├── docs/
│   └── api.md                  # REST API reference
├── go.mod
├── go.sum
├── README.md
├── CONTRIBUTING.md
├── SECURITY.md
└── LICENSE
```

## Database schema

### `bounties` table

```sql
CREATE TABLE bounties (
    id              BIGINT PRIMARY KEY,
    owner           TEXT NOT NULL,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL,
    amount          NUMERIC NOT NULL,       -- stored as raw i128 string
    token           TEXT NOT NULL,
    status          TEXT NOT NULL,          -- open|claimed|completed|cancelled
    claimant        TEXT,                   -- null until claimed
    created_at      BIGINT NOT NULL,        -- ledger sequence
    claim_deadline  BIGINT NOT NULL,        -- 0 = no deadline
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bounties_status ON bounties(status);
CREATE INDEX idx_bounties_owner  ON bounties(owner);
```

### `events` table

```sql
CREATE TABLE events (
    id              BIGSERIAL PRIMARY KEY,
    bounty_id       BIGINT NOT NULL REFERENCES bounties(id),
    event_type      TEXT NOT NULL,          -- posted|claimed|completed|cancelled
    ledger          BIGINT NOT NULL,
    tx_hash         TEXT NOT NULL,
    payload         JSONB NOT NULL,         -- full event data
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_bounty_id ON events(bounty_id);
CREATE INDEX idx_events_ledger    ON events(ledger);
```

### `cursor` table

```sql
CREATE TABLE cursor (
    id              INT PRIMARY KEY DEFAULT 1,   -- always one row
    last_ledger     BIGINT NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## Ingestion loop

```
startup
  ↓
read last_ledger from cursor table
  ↓
connect to Stellar RPC
  ↓
loop:
  fetch events for contract from (last_ledger + 1) to latest
    ↓
  for each event:
    decode XDR → Go struct
      ↓
    begin transaction
      dispatch to event handler (posted/claimed/completed/cancelled)
        ↓
      update bounties table
        ↓
      insert into events table
        ↓
      update cursor.last_ledger
    commit transaction
      ↓
  sleep 5s (one Stellar ledger close ≈ 5s)
  ↓
on RPC error: log, backoff, reconnect, resume from cursor
```

The entire event processing for a single event is wrapped in a database
transaction. If anything fails, the cursor does not advance and the
event is retried on the next loop iteration. This ensures no event is
silently skipped.

## Event handlers

```go
// internal/ingest/handlers.go

func handlePosted(ctx context.Context, db *DB, event ContractEvent) error
// Creates a new Bounty row with status=open

func handleClaimed(ctx context.Context, db *DB, event ContractEvent) error
// Updates bounty: status=claimed, claimant=<address>

func handleCompleted(ctx context.Context, db *DB, event ContractEvent) error
// Updates bounty: status=completed

func handleCancelled(ctx context.Context, db *DB, event ContractEvent) error
// Updates bounty: status=cancelled, claimant=null
```

Each handler is a pure function that takes a decoded event and performs
exactly one database update. They are tested independently.

## REST API

Base path: `/api/v1`

```
GET /api/v1/bounties
  Query params: status (open|claimed|completed|cancelled), limit, offset
  Response: { bounties: Bounty[], total: number, limit: number, offset: number }

GET /api/v1/bounties/:id
  Response: Bounty | 404

GET /api/v1/bounties/:id/events
  Response: { events: Event[] }

GET /api/v1/stats
  Response: { total: number, open: number, claimed: number,
              completed: number, cancelled: number, total_usdc: string }

GET /health
  Response: { status: "ok", last_ledger: number, lag: number }
  lag = current_ledger - last_ledger (in ledgers)
```

All responses are JSON. Errors use standard HTTP status codes with a
JSON body: `{ "error": "message" }`.

CORS headers allow requests from the app's Vercel domain and localhost.

## Go types

```go
type Bounty struct {
    ID            uint64  `json:"id"`
    Owner         string  `json:"owner"`
    Title         string  `json:"title"`
    Description   string  `json:"description"`
    Amount        string  `json:"amount"`       // raw i128 as decimal string
    Token         string  `json:"token"`
    Status        string  `json:"status"`
    Claimant      *string `json:"claimant"`
    CreatedAt     uint64  `json:"created_at"`
    ClaimDeadline uint64  `json:"claim_deadline"`
    UpdatedAt     string  `json:"updated_at"`
}

type Event struct {
    ID        uint64          `json:"id"`
    BountyID  uint64          `json:"bounty_id"`
    EventType string          `json:"event_type"`
    Ledger    uint64          `json:"ledger"`
    TxHash    string          `json:"tx_hash"`
    Payload   json.RawMessage `json:"payload"`
    CreatedAt string          `json:"created_at"`
}
```

## Configuration

```bash
DATABASE_URL=postgresql://...
STELLAR_RPC_URL=https://soroban-testnet.stellar.org
CONTRACT_ID=C...
NETWORK_PASSPHRASE=Test SDF Network ; September 2015
PORT=8080
START_LEDGER=0   # 0 = start from contract deployment ledger
```

## Deployment

Deploy to Render:
1. Create a Render PostgreSQL instance (free tier)
2. Create a Render web service pointing to this repo
3. Set environment variables in Render dashboard
4. Render builds with `go build ./cmd/indexer` and runs the binary
5. Service is publicly accessible at `https://soroquest-indexer.onrender.com`

Run migrations on first deploy: `golang-migrate` reads from
`migrations/` directory.

## Testing approach

```go
// Use a real test database (PostgreSQL in Docker or test container)
// Never mock the database in handler tests

func TestHandlePosted(t *testing.T) {
    db := setupTestDB(t)
    event := makePostedEvent(1, "alice", 10_000_000)
    err := handlePosted(ctx, db, event)
    require.NoError(t, err)

    bounty, err := db.GetBounty(ctx, 1)
    require.NoError(t, err)
    assert.Equal(t, "open", bounty.Status)
    assert.Equal(t, "10000000", bounty.Amount)
}
```

The ingestion loop is tested with a fake RPC client that returns
pre-built events. The database handlers are tested against a real
PostgreSQL instance.

## Key design decisions

**Database transaction per event:** wrapping the full event processing
(bounty update + event log + cursor advance) in a single transaction
means the system is always consistent. There is no state where an event
was processed but the cursor wasn't advanced, or vice versa.

**Amount stored as string:** Soroban's `i128` can exceed JavaScript's
`Number.MAX_SAFE_INTEGER`. Amounts are stored as decimal strings in
PostgreSQL (`NUMERIC` type) and returned as strings in the API. The
frontend uses `BigInt` to handle them.

**Polling over WebSocket:** Stellar RPC supports event streaming but the
polling approach (fetch events since last ledger, sleep, repeat) is
simpler to implement correctly, easier to resume after failure, and
sufficient for a 5-second latency target.
