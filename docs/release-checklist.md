# Release Checklist (v1)

## 1. Environment

- [ ] `DATABASE_URL` configured for gateway runtime
- [ ] `PROMETHEUS_API_KEY` configured for production (recommended)
- [ ] Postgres reachable from gateway deployment
- [ ] Web app can reach gateway base URL

## 2. Schema & Migrations

- [ ] Run `scripts/migrate.sh` successfully against target DB
- [ ] Confirm tables exist: `agents`, `tasks`, `agent_events`, `replay_query_metrics`
- [ ] Confirm replay indexes exist for `agent_events`

## 3. Local Smoke Test

- [ ] Start stack (`make setup`, gateway, web)
- [ ] Run `scripts/smoke-test.sh`
- [ ] Verify smoke test exits `0`

## 4. CI Gates

- [ ] Rust workspace check passes
- [ ] Python syntax check passes
- [ ] Go tests pass
- [ ] Go race tests pass

## 5. API Health

- [ ] `GET /api/v1/swarm/state` returns `200`
- [ ] `POST /api/v1/agents` returns `201`
- [ ] `POST /api/v1/agents/:id/tasks` returns `201`
- [ ] `GET /api/v1/events` returns items + cursor
- [ ] `GET /api/v1/metrics/replay` returns counters
- [ ] `GET /api/v1/metrics/replay/percentiles` returns `p50/p95/p99`

## 6. Security

- [ ] API key middleware enabled in production
- [ ] Secrets supplied via secret store or environment only
- [ ] No plaintext credentials committed

## 7. Deploy

- [ ] Apply `infra/kubernetes/namespace.yaml`
- [ ] Apply deployments, services, ingress, and HPA manifests
- [ ] Validate pods healthy and endpoints reachable

## 8. Demo Readiness

- [ ] Dashboard loads persisted agents/tasks and live stream
- [ ] Replay filters, presets, and pagination verified
- [ ] Replay metrics panel updates (summary + trend + percentiles)
