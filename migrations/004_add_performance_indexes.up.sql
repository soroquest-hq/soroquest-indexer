-- Compound indexes for fast status + timestamp queries
CREATE INDEX IF NOT EXISTS idx_bounties_status_created_at ON bounties (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_bounties_owner_status ON bounties (owner, status);
CREATE INDEX IF NOT EXISTS idx_bounties_claimant_status ON bounties (claimant, status) WHERE claimant IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_events_bounty_id ON events (bounty_id);
