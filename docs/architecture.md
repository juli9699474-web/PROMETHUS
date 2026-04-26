# PROMETHEUS Architecture (Phase 11)

PROMETHEUS uses a polyglot architecture:

- Rust runtime for agent lifecycle, market, protocol, consciousness, and evolution loops.
- Python AI subsystem for model routing, prompt evolution, and tool forging.
- Go gateway for REST, UAP ingress, and real-time WebSocket fanout.
- Next.js dashboard for live visibility and operator controls.

## Runtime Dataflow

1. Task enters gateway (`/api/v1/...`).
2. Gateway forwards to runtime/UAP boundary.
3. Runtime schedules and auctions tasks across agents.
4. Agent execution emits swarm events to event bus.
5. Gateway streams events to clients via `/ws/swarm`.
6. Outcomes are persisted in Postgres and indexed in memory systems.
