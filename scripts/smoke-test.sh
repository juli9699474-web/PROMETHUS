#!/usr/bin/env bash
set -euo pipefail

API_URL="${API_URL:-http://localhost:3001}"
API_KEY="${PROMETHEUS_API_KEY:-}"

echo "== PROMETHUS smoke test =="
echo "API_URL=$API_URL"

curl -fsS "$API_URL/api/v1/swarm/state" >/dev/null
echo "OK: swarm state"

AGENT_JSON="$(curl -fsS -X POST "$API_URL/api/v1/agents" -H "Content-Type: application/json" -d '{"name":"SmokeAgent"}')"
AGENT_ID="$(python - <<'PY'
import json,sys
print(json.loads(sys.stdin.read())["id"])
PY
<<<"$AGENT_JSON")"
echo "OK: create agent ($AGENT_ID)"

curl -fsS -X POST "$API_URL/api/v1/agents/$AGENT_ID/tasks" \
  -H "Content-Type: application/json" \
  -d '{"title":"SmokeTask","description":"smoke","reward":1.0}' >/dev/null
echo "OK: assign task"

curl -fsS "$API_URL/api/v1/events?limit=10" >/dev/null
echo "OK: events query"

if [[ -n "$API_KEY" ]]; then
  curl -fsS "$API_URL/api/v1/metrics/replay" -H "X-API-Key: $API_KEY" >/dev/null
  curl -fsS "$API_URL/api/v1/metrics/replay/percentiles?limit=50" -H "X-API-Key: $API_KEY" >/dev/null
  echo "OK: metrics (api key)"
else
  echo "SKIP: metrics endpoints (PROMETHEUS_API_KEY not set)"
fi

echo "Smoke test passed."

