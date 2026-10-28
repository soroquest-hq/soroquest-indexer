# ADR 0001: At-Least-Once Event Ingestion and Idempotency

## Status
Accepted

## Context
Soroban contract events are polled via the `getEvents` RPC method. RPC nodes can experience network partitions, rate limits, or transient timeouts. The indexer must guarantee that no contract events (bounty posted, claimed, approved, cancelled) are permanently skipped, and that duplicate event deliveries do not corrupt relational database state.

## Decision
1. **Ledger Cursor as Single Source of Progress**:
   - The indexer tracks the latest successfully ingested ledger sequence number in the `indexer_cursor` table.
   - The polling window always requests events starting from `last_ledger + 1`.

2. **Atomic Ingestion Transactions**:
   - Updating the bounty status, inserting into the immutable `events` table, and advancing the `indexer_cursor` are wrapped in a single database transaction (`tx.Begin() ... tx.Commit()`).
   - If a failure occurs mid-batch, the transaction rolls back, and the cursor does not advance.

3. **Idempotent Upserts**:
   - Event insertions use unique constraints on `(bounty_id, event_type, ledger_sequence)`.
   - Repeated processing of the same ledger produces identical database state without side effects.

## Consequences
- At-least-once processing guarantees zero missing bounties.
- Crash-resilient state recovery upon indexer restarts.
