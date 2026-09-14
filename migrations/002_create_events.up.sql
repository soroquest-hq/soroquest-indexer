CREATE TABLE IF NOT EXISTS events (
    id              BIGSERIAL PRIMARY KEY,
    bounty_id       BIGINT NOT NULL REFERENCES bounties(id),
    event_type      TEXT NOT NULL CHECK (event_type IN ('posted','claimed','completed','cancelled')),
    ledger          BIGINT NOT NULL,
    tx_hash         TEXT NOT NULL,
    payload         JSONB NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_events_bounty_id ON events(bounty_id);
CREATE INDEX IF NOT EXISTS idx_events_ledger    ON events(ledger);
