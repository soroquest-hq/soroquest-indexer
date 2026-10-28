#!/usr/bin/env bash
set -euo pipefail

# SoroQuest Indexer Demo Data Seeding Helper
# Seeds local PostgreSQL with sample testnet bounties for UI verification

DATABASE_URL="${DATABASE_URL:-postgres://soroquest_user:soroquest_password@localhost:5432/soroquest?sslmode=disable}"

echo "=== Seeding SoroQuest Demo Bounties ==="

if ! command -v psql &> /dev/null; then
    echo "Notice: psql not installed locally. Bounties can be posted directly on-chain via https://soroquest-app.vercel.app."
    exit 0
fi

psql "$DATABASE_URL" << 'EOF'
INSERT INTO bounties (id, owner, title, description, amount, token, status, created_at, claim_deadline)
VALUES 
  (1, 'GA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJUO5U', 'Implement Telegram Notifications for Bounty Claims', 'Build a lightweight webhook listener or bot that sends a real-time message to the community Telegram channel.', 50000000, 'CBIELTK6YBZJU5UP2WWQEUCYKLPU6AUNZ2BQ4WWFEIE3USCIHMXQDAMA', 'open', 124500, 0),
  (2, 'GA7QYNF7SOWQ3GLR2BGMZEHXAVIRZA4KVWLTJJFC7MGXUA74P7UJUO5U', 'Add Dark Mode Theme Support to SoroQuest Web UI', 'Extend the Tailwind configuration and UI components to support system dark mode and toggleable theme.', 100000000, 'CBIELTK6YBZJU5UP2WWQEUCYKLPU6AUNZ2BQ4WWFEIE3USCIHMXQDAMA', 'claimed', 123100, 135000)
ON CONFLICT (id) DO NOTHING;
EOF

echo "Demo bounties seeded successfully."
