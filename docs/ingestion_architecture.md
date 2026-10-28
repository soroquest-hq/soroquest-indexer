# SoroQuest Indexer: Event Ingestion Architecture & Resiliency

This document details the architectural design, polling strategy, error handling, and replay mechanics of `soroquest-indexer`.

```mermaid
flowchart TD
    A["Soroban RPC getEvents"] -->|Filter by Contract ID| B["Ingest Loop (5s Poll)"]
    B --> C{"Events Returned?"}
    C -->|No| D["Sleep & Wait for Next Ledger"]
    C -->|Yes| E["Decode XDR ScVal"]
    E --> F["Parse Topics & Data (Post, Claim, Approve, Cancel)"]
    F --> G["PostgreSQL Transaction"]
    G --> H["Upsert Bounty State"]
    G --> I["Insert Immutable Event Log"]
    G --> J["Update Ingestion Cursor (ledger_sequence)"]
    J --> K["Commit Transaction"]
```

---

## Cursor Tracking & At-Least-Once Delivery
- **Ledger Watermark**: The cursor is stored in the `indexer_cursor` table.
- **Atomic Progress**: State updates, event logs, and the cursor update execute within a single PostgreSQL transaction (`BEGIN ... COMMIT`).
- **Idempotency**: If the service restarts mid-poll, transactions roll back. The next poll begins from `last_ledger`, preventing duplicated state transitions or skipped bounties.
- **XDR Decoding**: Soroban events emit raw XDR ScVals. The decoder validates schema structures before database writes.
