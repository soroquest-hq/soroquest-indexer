# SoroQuest Indexer — REST API Reference

Base URL: `https://soroquest-indexer.onrender.com/api/v1`

All responses are JSON. Errors use standard HTTP status codes:
```json
{ "error": "message" }
```

## Endpoints

### GET /api/v1/bounties

List bounties with optional filtering and pagination.

**Query parameters:**
- `status` — filter by status: `open`, `claimed`, `completed`, `cancelled`
- `limit` — number of results (default: 20, max: 100)
- `offset` — pagination offset (default: 0)

**Response:**
```json
{
  "bounties": [Bounty],
  "total": 42,
  "limit": 20,
  "offset": 0
}
```

### GET /api/v1/bounties/:id

Get a single bounty by ID.

**Response:** `Bounty` or `404`

### GET /api/v1/bounties/:id/events

Get the full event history for a bounty.

**Response:**
```json
{ "events": [Event] }
```

### GET /api/v1/stats

Platform-wide aggregate statistics.

**Response:**
```json
{
  "total": 100,
  "open": 42,
  "claimed": 10,
  "completed": 45,
  "cancelled": 3,
  "total_usdc": "450000000"
}
```

### GET /health

Liveness check. Shows ingestion lag.

**Response:**
```json
{
  "status": "ok",
  "last_ledger": 54321000,
  "lag": 2
}
```

## Data types

### Bounty
```json
{
  "id": 1,
  "owner": "GXXX...",
  "title": "Fix login bug",
  "description": "...",
  "amount": "10000000",
  "token": "CBIELTK6...",
  "status": "open",
  "claimant": null,
  "created_at": 54321000,
  "claim_deadline": 0,
  "updated_at": "2026-01-01T00:00:00Z"
}
```

### Event
```json
{
  "id": 1,
  "bounty_id": 1,
  "event_type": "posted",
  "ledger": 54321000,
  "tx_hash": "abc123...",
  "payload": {},
  "created_at": "2026-01-01T00:00:00Z"
}
```
