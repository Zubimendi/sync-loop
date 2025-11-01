# SyncLoop — Airbyte-lite for Internal Apps

SyncLoop is a lightweight, self-hostable data sync platform that connects your internal tools — databases, SaaS apps, and spreadsheets — without needing a full data engineering team.

> “Move data from any internal app to any other internal app in <5 min, keep it in sync, and prove it’s correct.”

---

## 🚨 The Problem

Ops and Finance teams in 20–200-employee companies run dozens of “data chores” every week:
- CSV exports, SFTP dumps, Excel attachments
- Silent job failures and stale data
- Scattered credentials and no audit trail

**SyncLoop** replaces these brittle scripts with a single, auditable sync service that “just works.”

---

## 👩‍💻 Target Personas

| Persona | Role | Pain Point |
|----------|------|-------------|
| **Sam** | Ops Manager | Needs yesterday’s sales CSV in Snowflake by 9 am |
| **Dana** | Finance Analyst | Reconciles QuickBooks vs Shopify weekly |
| **Lex** | IT Admin | Worries about GDPR, credentials, and uptime |

---

## 🚀 MVP Scope

### ✅ Core Features

**Web UI (React-TypeScript)**
- Connector gallery: Postgres, MySQL, Snowflake, S3, Excel, Google Sheets, Salesforce, REST
- Configuration wizard: auth test → schema detect → schedule
- Sync history: run logs, row counts, job status
- Email + Slack alerts on failure

**Runtime (Go)**
- Worker pool (Temporal.io) for incremental + full refresh
- Row-level checksum validation
- Secrets encrypted at rest (AES-256-GCM + KMS envelope)
- Basic transformations (column rename, static filters)

**Deployment**
- Single binary + Postgres (with embedded migrations)
- Docker Compose for BYO infra (Helm chart later)
- OTEL traces → Axiom (we dogfood our own telemetry)

---

## 🧱 Entity Model (Simplified)

Main entities and relationships:

workspace → users → connectors → sync_jobs → sync_runs → destinations → field_mappings


Postgres schema includes:
- **connector** (type, config, ownership)
- **sync_job** (schedule, status)
- **sync_run** (row counts, logs, checksum)
- **destination** (output config)
- **users / workspaces** (auth, roles, billing)
  
See `/docs/schema.sql` for full definition.

---

## 📊 Success Metrics (MVP)

| Metric | Target |
|---------|---------|
| Time-to-first-sync | ≤ 15 min (unassisted) |
| Failed-job alert latency | ≤ 5 min |
| Manual credential rotation | ≤ 1 × per connector / year |
| Paying pilot | ≥ 3 customers @ $199 / mo by month 4 |


## 🧩 Tech Stack

| Layer | Technology |
|--------|-------------|
| Frontend | React + TypeScript + Vite |
| Backend | Go (net/http, Temporal.io) |
| Database | Postgres |
| Auth | JWT + bcrypt |
| Encryption | AES-256-GCM + KMS envelope |
| Telemetry | OpenTelemetry → Axiom |
| Packaging | Docker Compose (Helm coming soon) |

---

## 🧰 Development

### Requirements
- Go ≥ 1.22  
- Node ≥ 20 + pnpm  
- Docker Desktop or Podman  
- Postgres ≥ 15

### Local Setup

```bash
# Backend
cd backend
go run main.go

# Frontend
cd frontend
pnpm install
pnpm dev

# Optional: launch Postgres + Axiom mock
docker compose up


# Backend unit tests
go test ./...

# Frontend tests
pnpm test


docker compose up -d

```

### Outputs:

Web UI: http://localhost:8080

API: http://localhost:8000/api

Logs: docker logs syncloop-api


### 🗺️ Roadmap
Phase	Focus
- Sprint 0	Schema, scaffolding, auth
- Sprint 1	Connectors: Postgres, S3, Excel
- Sprint 2	Incremental sync + alerts
- Sprint 3	Billing + telemetry
- Q2 2025	Real-time CDC + visual mapper UI

🧾 License

MIT © 2025 SyncLoop
