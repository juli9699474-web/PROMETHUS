## PROMETHUS

<p align="center">
  <img src="assets/logo.svg" alt="PROMETHUS Logo" width="120" />
</p>

<p align="center">
  <b>The swarm thinks. The swarm learns. The swarm evolves.</b><br/>
  A simple and universal swarm intelligence engine — mission-control UI + replayable event forensics.
</p>

<p align="center">
  <a href="https://github.com/juli9699474-web/PROMETHUS">GitHub</a>
  ·
  <a href="docs/getting-started.md">Getting Started</a>
  ·
  <a href="docs/architecture.md">Architecture</a>
  ·
  <a href="docs/api-reference.md">API Reference</a>
</p>

<p align="center">
  <a href="https://github.com/juli9699474-web/PROMETHUS/stargazers">
    <img alt="GitHub stars" src="https://img.shields.io/github/stars/juli9699474-web/PROMETHUS?style=flat" />
  </a>
  <a href="https://github.com/juli9699474-web/PROMETHUS/watchers">
    <img alt="GitHub watchers" src="https://img.shields.io/github/watchers/juli9699474-web/PROMETHUS?style=flat" />
  </a>
  <a href="https://github.com/juli9699474-web/PROMETHUS/fork">
    <img alt="GitHub forks" src="https://img.shields.io/github/forks/juli9699474-web/PROMETHUS?style=flat" />
  </a>
</p>

## ⚡ Overview

PROMETHUS is an agentic swarm OS scaffold focused on one thing: **observable, replayable swarm behavior**.
You can run a local “mission control”, create agents/tasks, watch the event stream live, then rewind and inspect
what happened with filters, cursor pagination, and metrics (latency + percentiles).

## 🌐 What you can do today

- **Run a local mission control**: dashboard + gateway
- **Create agents and assign tasks** via REST
- **Watch live swarm events** over WebSocket
- **Replay history** with filters (`eventType`, `agentId`, time range) + cursor paging
- **Observe replay system health**: rate limits, durable metrics, percentiles (p50/p95/p99)

## 🔄 Workflow

1. Start infrastructure (Postgres, Redis, etc.)
2. Run migrations
3. Start gateway (Go) + dashboard (Next.js)
4. Create agents / assign tasks
5. Inspect live + replay modes, export findings

## 🚀 Quick Start

### Option 1: Source deployment (recommended)

Prerequisites:
- Node.js 18+
- Go 1.22+
- Docker + Docker Compose

```bash
cp .env.example .env
make setup
make migrate
npm install
npm run dev
```

Run a quick validation:

```bash
npm run smoke
```

Open:
- **Dashboard**: `http://localhost:3000`
- **Gateway**: `http://localhost:3001`
- **Swarm WS**: `ws://localhost:3001/ws/swarm`

### Option 2: Docker (infra only)

```bash
make setup
make migrate
```

Then run gateway + web with `npm run dev`.

## 📸 Screenshots

- Add screenshots/GIFs to `docs/diagrams/` and link them here:
  - `docs/diagrams/dashboard.png`
  - `docs/diagrams/replay.png`
  - `docs/diagrams/metrics.png`

## 🎬 Demo

- Record a short GIF of the dashboard event stream + replay mode and embed it at the top of this README.

## 🔐 API Key (optional, recommended)

Set `PROMETHEUS_API_KEY` to protect replay metrics endpoints. Send it as `X-API-Key`.

## 📦 Deploy (Kubernetes baseline)

Manifests live under `infra/kubernetes/`.

## Phase 2 Added

- Protobuf contracts in `gateway/proto/`
- SQL bootstrap migration in `infra/sql/001_init.sql`
- Gateway WebSocket stream at `GET /ws/swarm`
- Runtime event bus in `core/prometheus-runtime/src/runtime/events.rs`
- Next.js mission-control shell in `web/`

## Phase 3 Added

- Real gateway API for agent/task lifecycle (`/api/v1/agents`, `/api/v1/agents/:id/tasks`)
- Event hub fanout for live WebSocket updates (`/ws/swarm`)
- Dashboard component split (`SwarmVisualizer`, `MarketTicker`, `ConsciousnessMap`)
- In-app control panel to trigger agent/task events end-to-end

## Phase 4 Added

- Repository abstraction in gateway for persistence-ready handlers
- Postgres-backed repository with in-memory fallback
- Append-only `agent_events` writes for every emitted swarm event
- SQL migration updated with `agent_events` table and indexes

## Phase 5 Added

- Historical event replay endpoint (`GET /api/v1/events`)
- Read-model endpoints for persisted agents and tasks
- Dashboard now bootstraps from REST state, then continues with live WebSocket updates

## Phase 6 Added

- Filtered event queries (`eventType`, `agentId`, `sinceMs`, `untilMs`)
- Cursor pagination for replay (`cursorMs`, `nextCursorMs`)
- Dashboard playback controls for live vs replay mode and history pagination

## Phase 7 Added

- Server-side event query validation for UUID `agentId` and time-range bounds
- Additional event-table indexes for filtered replay performance
- Dashboard date-range inputs mapped to `sinceMs` and `untilMs`

## Phase 8 Added

- Replay endpoint guardrails: bounded `limit` and per-IP rate limiting
- Replay query metrics endpoint (`/api/v1/metrics/replay`)
- Dashboard saved replay presets for one-click forensic filters

## Phase 9 Added

- Durable replay metric persistence in Postgres (`replay_query_metrics`)
- Identity-aware replay quotas using `X-API-Key` with IP fallback
- Dashboard replay latency and rejection trend panel

## Phase 10 Added

- Replay metrics timeseries endpoint for historical charting
- Replay latency percentile endpoint (`p50`, `p95`, `p99`)
- Dashboard persisted preset selection and percentile visibility

## Phase 11 Added

- API-key middleware for protected replay metrics endpoints
- Gateway smoke tests for health and API-key protection
- Migration runner script (`scripts/migrate.sh`)
- Baseline Kubernetes namespace, deployments, services, and ingress manifests
