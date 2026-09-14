# PRD — soroquest-indexer

## What we are building

A backend service that listens to the SoroQuest contract's events on
Stellar, maintains a local database of bounty state, and exposes a REST
API that the frontend queries instead of hitting the contract directly.

Without the indexer, every page load in the app requires multiple RPC
calls to iterate contract storage — which is slow, rate-limited, and
produces a poor user experience. The indexer solves this by maintaining
a synchronized copy of contract state that the app can query instantly.

## Who it is for

**Primary: the soroquest-app frontend**
Every read in the app goes to the indexer first. The indexer is an
internal service — it has no end users itself.

**Secondary: third-party developers**
Anyone building on top of SoroQuest who needs historical bounty data,
event streams, or analytics. The REST API is public and documented.

## What the product needs to do

### Ingest contract events

- Connect to Stellar RPC and stream events from the SoroQuest contract
- Process all four event types: `posted`, `claimed`, `completed`,
  `cancelled`
- On each event, update the local database to reflect the new state
- Track the last processed ledger so restarts resume without replaying
  the full history
- Handle RPC disconnections gracefully — reconnect and resume from the
  last processed ledger without gaps or duplicates

### Maintain bounty state

- Store the full `Bounty` record for every bounty that has ever existed
- Update the record on every state-changing event
- Keep an event log for historical queries and debugging
- Never lose data on restart — the database is the source of truth for
  all historical state

### Serve a REST API

The API is what the frontend and third-party developers consume.

**Bounty endpoints:**

```
GET /bounties              List bounties, filterable by status, paginated
GET /bounties/:id          Single bounty by ID
GET /bounties/:id/events   Full event history for a bounty
GET /stats                 Platform totals: total bounties, total USDC,
                           open count, completed count
GET /health                Liveness check — returns last indexed ledger
```

### Operational requirements

- Lag behind the chain by no more than a few ledgers under normal
  conditions
- Recover automatically from RPC failures without manual intervention
- Log enough information to diagnose indexing gaps after the fact
- Expose a `/health` endpoint that shows the last processed ledger so
  the frontend can detect when the indexer is behind

## What the product must NOT do

- Never be a write path — all writes go directly to the contract
- Never cache stale data without a way to detect it is stale
- Never drop events silently — if an event cannot be processed, log it
  and retry

## Acceptance criteria

- All four contract event types are indexed correctly
- Bounty state in the database matches on-chain state after each event
- API returns correct data for all endpoints
- Indexer recovers from a simulated RPC disconnect without gaps
- Deployed and accessible at a public URL
- `/health` endpoint returns the last indexed ledger
- `go test ./...` passes with zero failures
