CREATE TABLE IF NOT EXISTS bounties (
    id              BIGINT PRIMARY KEY,
    owner           TEXT NOT NULL,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL,
    amount          NUMERIC NOT NULL,
    token           TEXT NOT NULL,
    status          TEXT NOT NULL CHECK (status IN ('open','claimed','completed','cancelled')),
    claimant        TEXT,
    created_at      BIGINT NOT NULL,
    claim_deadline  BIGINT NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bounties_status ON bounties(status);
CREATE INDEX IF NOT EXISTS idx_bounties_owner  ON bounties(owner);
