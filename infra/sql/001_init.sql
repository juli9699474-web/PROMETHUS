CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    genome JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'idle',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_active_at TIMESTAMPTZ,
    fitness_score DOUBLE PRECISION DEFAULT 0.5,
    token_balance DOUBLE PRECISION DEFAULT 100.0,
    total_tasks_completed INTEGER DEFAULT 0,
    total_tasks_failed INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    priority INTEGER DEFAULT 5,
    assigned_to UUID REFERENCES agents(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    result JSONB,
    token_reward DOUBLE PRECISION DEFAULT 10.0
);

CREATE TABLE IF NOT EXISTS market_trades (
    id BIGSERIAL PRIMARY KEY,
    task_id UUID REFERENCES tasks(id),
    winning_agent UUID REFERENCES agents(id),
    winning_bid DOUBLE PRECISION,
    auction_duration_ms INTEGER,
    num_bidders INTEGER,
    settled_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS agent_events (
    id BIGSERIAL PRIMARY KEY,
    agent_id UUID NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agent_events_agent_id ON agent_events (agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_events_timestamp ON agent_events (timestamp);
CREATE INDEX IF NOT EXISTS idx_agent_events_event_type ON agent_events (event_type);
CREATE INDEX IF NOT EXISTS idx_agent_events_event_type_timestamp ON agent_events (event_type, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_agent_events_agent_id_timestamp ON agent_events (agent_id, timestamp DESC);

CREATE TABLE IF NOT EXISTS replay_query_metrics (
    id BIGSERIAL PRIMARY KEY,
    client_key TEXT NOT NULL,
    query_limit INTEGER NOT NULL,
    result_count INTEGER NOT NULL,
    latency_ms BIGINT NOT NULL,
    rejected BOOLEAN NOT NULL DEFAULT FALSE,
    filter_event_type TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_replay_query_metrics_created_at ON replay_query_metrics (created_at DESC);
