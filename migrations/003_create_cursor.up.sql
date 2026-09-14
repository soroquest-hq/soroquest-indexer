CREATE TABLE IF NOT EXISTS cursor (
    id          INT PRIMARY KEY DEFAULT 1,
    last_ledger BIGINT NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Ensure there is always exactly one row
INSERT INTO cursor (id, last_ledger) VALUES (1, 0)
ON CONFLICT (id) DO NOTHING;
