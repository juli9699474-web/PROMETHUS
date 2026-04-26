# PROMETHEUS AGI OS

The swarm thinks. The swarm learns. The swarm evolves.

## Quick Start

```bash
make setup
make all
```

## Current Status

This repository now includes a Phase-1 executable skeleton:

- Rust runtime with genome, market, consciousness, evolution, and UAP modules
- Python evolution and tool forge scaffolding
- Go gateway with core endpoints
- Base docs and compose stack

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
