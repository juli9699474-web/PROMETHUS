# API Reference (Phase 11)

## REST

- `GET /api/v1/swarm/state`
- `GET /api/v1/agents?limit=50`
- `POST /api/v1/agents`
- `GET /api/v1/agents/:id`
- `POST /api/v1/agents/:id/tasks`
- `GET /api/v1/tasks?limit=50`
- `GET /api/v1/events?limit=200`
  - Filters: `eventType`, `agentId`, `sinceMs`, `untilMs`, `cursorMs`
  - Validation: `agentId` must be UUID, `sinceMs <= untilMs`
  - Guardrails: `limit` is bounded to `500`, replay queries are rate-limited by `X-API-Key` (fallback to client IP)
- `GET /api/v1/metrics/replay`
  - Returns replay query counters and average latency
  - Protected when `PROMETHEUS_API_KEY` is configured (send `X-API-Key`)
- `GET /api/v1/metrics/replay/timeseries?limit=100`
  - Returns persisted replay metric points (`createdAtMs`, `latencyMs`, `rejected`)
- `GET /api/v1/metrics/replay/percentiles?limit=500`
  - Returns `p50`, `p95`, `p99` replay latency percentiles
- `POST /uap/connect`
- `POST /uap/message`

## WebSocket

- `GET /ws/swarm`
  - Streams lifecycle and heartbeat events:
  - `eventType`, `payload`, `tsUnixMs`
