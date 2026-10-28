# Changelog

All notable changes to the SoroQuest Indexer service will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-10-28

### Added
- **Ingestion Engine**:
  - Polling loop against Soroban Testnet RPC with at-least-once delivery semantics.
  - XDR event decoder mapping `post`, `claim`, `approve`, and `cancel` events to relational database records.
  - Ledger watermark persistence in `indexer_cursor` table.
- **REST API**:
  - `GET /health`: Diagnostic health endpoint returning database status, last ingested ledger, and process uptime.
  - `GET /api/bounties`: Filterable, paginated bounty catalog with compound index queries.
  - `GET /api/bounties/{id}`: Direct bounty lookup by unique contract ID.
  - `GET /api/bounties/{id}/events`: Historical event audit trail per bounty.
- **Deployment & Cloud**:
  - Automated deployment Blueprint for Render (`render.yaml`) with PostgreSQL database.
  - Multi-stage `Dockerfile` and local `docker-compose.yml`.
  - Comprehensive OpenAPI 3.0 specification (`docs/openapi.yaml`).
