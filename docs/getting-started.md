# Getting Started

## Prerequisites

- Docker + Docker Compose
- Rust toolchain
- Python 3.11+
- Go 1.22+
- Node.js 20+

## Boot

1. `make setup`
2. Apply SQL bootstrap: `bash scripts/migrate.sh`
3. Set `DATABASE_URL=postgres://prometheus:devpassword@localhost:5432/prometheus`
4. Optional API protection: set `PROMETHEUS_API_KEY=<strong-key>`
5. `cargo check --workspace`
6. `go run ./gateway/cmd/gateway`
7. `cd web && npm install && npm run dev`

Then open `http://localhost:3000` and confirm lifecycle and heartbeat events from `ws://localhost:3001/ws/swarm`.
