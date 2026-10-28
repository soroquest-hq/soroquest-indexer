# SoroQuest Indexer: Environment Variables Reference

This document provides a reference for all environment configuration keys used by `soroquest-indexer`.

| Variable | Required | Default | Description |
| :--- | :---: | :---: | :--- |
| `DATABASE_URL` | **Yes** | — | PostgreSQL connection string (`postgres://user:pass@host:5432/db?sslmode=disable`). |
| `STELLAR_RPC_URL` | **Yes** | — | Soroban JSON-RPC endpoint (`https://soroban-testnet.stellar.org`). |
| `CONTRACT_ID` | **Yes** | — | SoroQuest Escrow contract address (`CA7C6LGWPTZVUTZX7EXJVLWCQG6DVWPROBPO3FGXYKEY4W5H4HRPO5UV`). |
| `NETWORK_PASSPHRASE` | **Yes** | — | Stellar network passphrase (`Test SDF Network ; September 2015`). |
| `PORT` | No | `8080` | HTTP port for the REST API server. |
| `ALLOWED_ORIGIN` | No | `*` | Allowed CORS origins for browser dApp requests (e.g. `https://soroquest-app.vercel.app`). |
| `POLL_INTERVAL_MS` | No | `5000` | Polling frequency for Soroban RPC ledger queries (in milliseconds). |
