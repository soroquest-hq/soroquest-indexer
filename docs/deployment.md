# SoroQuest Indexer — Deployment Guide

This guide describes how to deploy the SoroQuest Indexer daemon and PostgreSQL database on Render.

---

## 1. Prerequisites

- A [Render](https://render.com) account.
- The `soroquest-hq/soroquest-indexer` repository connected to your Render account.
- Deployed Soroban escrow contract ID on Stellar Testnet (`CA7C6LGWPTZVUTZX7EXJVLWCQG6DVWPROBPO3FGXYKEY4W5H4HRPO5UV`).

---

## 2. Option A: Blueprint Deployment (Recommended)

1. Navigate to the **Blueprints** tab in the Render Dashboard.
2. Click **New Blueprint Instance**.
3. Select the `soroquest-hq/soroquest-indexer` repository.
4. Render will parse [`render.yaml`](../render.yaml) and automatically create:
   - A managed PostgreSQL database (`soroquest-db`).
   - A Go web service (`soroquest-indexer`).
5. Click **Apply**.

---

## 3. Option B: Manual Deployment

### Step 1: Create a PostgreSQL Database
1. Go to the Render Dashboard → **New** → **PostgreSQL**.
2. Set:
   - **Name:** `soroquest-db`
   - **Database:** `soroquest`
   - **User:** `soroquest_user`
   - **Plan:** Free
3. Once provisioned, copy the **Internal Database URL** (e.g. `postgres://...`).

### Step 2: Create the Web Service
1. In Render Dashboard, click **New** → **Web Service**.
2. Connect `soroquest-hq/soroquest-indexer`.
3. Set the build and start settings:
   - **Runtime:** Go
   - **Build Command:** `go build -o bin/indexer ./cmd/indexer`
   - **Start Command:** `./bin/indexer`
   - **Plan:** Free
4. Configure Environment Variables:
   - `PORT`: `8080`
   - `DATABASE_URL`: *(paste the Internal Database URL from Step 1)*
   - `STELLAR_RPC_URL`: `https://soroban-testnet.stellar.org`
   - `STELLAR_NETWORK_PASSPHRASE`: `Test SDF Network ; September 2015`
   - `CONTRACT_ID`: `CA7C6LGWPTZVUTZX7EXJVLWCQG6DVWPROBPO3FGXYKEY4W5H4HRPO5UV`
5. Click **Deploy Web Service**.

---

## 4. Database Migrations

The indexer automatically executes pending SQL migrations located in `migrations/` on startup. The tables created are:
- `bounties`: Bounty state and details.
- `events`: Historical contract event logs.
- `cursor`: Ingestion checkpoint tracking the latest processed ledger sequence.

---

## 5. Verification

1. Once the service builds and starts, obtain the public URL (e.g., `https://soroquest-indexer.onrender.com`).
2. Test the `/health` endpoint:
   ```bash
   curl -i https://soroquest-indexer.onrender.com/health
   ```
   Expected response:
   ```json
   {
     "status": "ok",
     "last_ledger": 123456,
     "lag": 2
   }
   ```
3. Test the bounties list endpoint:
   ```bash
   curl -i https://soroquest-indexer.onrender.com/api/bounties
   ```

---

## 6. Update the Frontend on Vercel

Once your Render Indexer URL is confirmed:
1. Open the **`soroquest-app`** project on [Vercel](https://vercel.com).
2. Go to **Settings** → **Environment Variables**.
3. Set `NEXT_PUBLIC_INDEXER_URL` = `https://<your-indexer-service>.onrender.com`.
4. Trigger a redeploy on Vercel.

