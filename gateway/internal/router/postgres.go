package router

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, databaseURL string) (*PostgresRepository, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create pg pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping pg: %w", err)
	}
	return &PostgresRepository{pool: pool}, nil
}

func (r *PostgresRepository) Close() {
	r.pool.Close()
}

func (r *PostgresRepository) CreateAgent(ctx context.Context, name string) (Agent, error) {
	id := uuid.New()
	genome := map[string]any{"name": name, "version": "phase4"}
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO agents (id, genome, status, created_at, fitness_score) VALUES ($1, $2, 'idle', NOW(), 0.5)`,
		id,
		genome,
	)
	if err != nil {
		return Agent{}, err
	}
	return Agent{
		ID:        id.String(),
		Name:      name,
		Status:    "idle",
		Fitness:   0.5,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (r *PostgresRepository) GetAgent(ctx context.Context, id string) (Agent, error) {
	var (
		status    string
		createdAt time.Time
		fitness   float64
		genomeRaw []byte
	)
	err := r.pool.QueryRow(
		ctx,
		`SELECT status, created_at, fitness_score, genome FROM agents WHERE id = $1`,
		id,
	).Scan(&status, &createdAt, &fitness, &genomeRaw)
	if err != nil {
		if isNoRows(err) {
			return Agent{}, ErrNotFound
		}
		return Agent{}, err
	}
	name := "agent"
	var genome map[string]any
	if json.Unmarshal(genomeRaw, &genome) == nil {
		if v, ok := genome["name"].(string); ok && v != "" {
			name = v
		}
	}
	return Agent{
		ID:        id,
		Name:      name,
		Status:    status,
		Fitness:   fitness,
		CreatedAt: createdAt,
	}, nil
}

func (r *PostgresRepository) CreateTask(ctx context.Context, agentID, title, description string, reward float64) (Task, error) {
	id := uuid.New()
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO tasks (id, title, description, status, assigned_to, token_reward, created_at)
		 VALUES ($1, $2, $3, 'queued', $4, $5, NOW())`,
		id,
		title,
		description,
		agentID,
		reward,
	)
	if err != nil {
		return Task{}, err
	}
	return Task{
		ID:          id.String(),
		AgentID:     agentID,
		Title:       title,
		Description: description,
		Status:      "queued",
		Reward:      reward,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func (r *PostgresRepository) ListAgents(ctx context.Context, limit int) ([]Agent, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.pool.Query(
		ctx,
		`SELECT id, status, created_at, fitness_score, genome
		 FROM agents ORDER BY created_at DESC LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Agent, 0, limit)
	for rows.Next() {
		var (
			id        string
			status    string
			createdAt time.Time
			fitness   float64
			genomeRaw []byte
		)
		if err := rows.Scan(&id, &status, &createdAt, &fitness, &genomeRaw); err != nil {
			return nil, err
		}
		name := "agent"
		var genome map[string]any
		if json.Unmarshal(genomeRaw, &genome) == nil {
			if v, ok := genome["name"].(string); ok && v != "" {
				name = v
			}
		}
		out = append(out, Agent{
			ID: id, Name: name, Status: status, Fitness: fitness, CreatedAt: createdAt,
		})
	}
	return out, rows.Err()
}

func (r *PostgresRepository) ListTasks(ctx context.Context, limit int) ([]Task, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.pool.Query(
		ctx,
		`SELECT id, assigned_to, title, description, status, token_reward, created_at
		 FROM tasks ORDER BY created_at DESC LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Task, 0, limit)
	for rows.Next() {
		var (
			id, title, description, status string
			agentID                        *string
			reward                         float64
			createdAt                      time.Time
		)
		if err := rows.Scan(&id, &agentID, &title, &description, &status, &reward, &createdAt); err != nil {
			return nil, err
		}
		aid := ""
		if agentID != nil {
			aid = *agentID
		}
		out = append(out, Task{
			ID: id, AgentID: aid, Title: title, Description: description, Status: status, Reward: reward, CreatedAt: createdAt,
		})
	}
	return out, rows.Err()
}

func (r *PostgresRepository) ListEvents(ctx context.Context, limit int) ([]SwarmEvent, error) {
	return r.QueryEvents(ctx, EventQuery{Limit: limit})
}

func (r *PostgresRepository) QueryEvents(ctx context.Context, query EventQuery) ([]SwarmEvent, error) {
	limit := query.Limit
	if limit <= 0 {
		limit = 100
	}
	sql := `SELECT event_type, payload, EXTRACT(EPOCH FROM timestamp) * 1000
		FROM agent_events WHERE 1=1`
	args := make([]any, 0, 8)
	argN := 1
	if query.EventType != "" {
		sql += fmt.Sprintf(" AND event_type = $%d", argN)
		args = append(args, query.EventType)
		argN++
	}
	if query.AgentID != "" {
		sql += fmt.Sprintf(" AND agent_id = $%d", argN)
		args = append(args, query.AgentID)
		argN++
	}
	if query.SinceMs > 0 {
		sql += fmt.Sprintf(" AND timestamp >= to_timestamp($%d::double precision / 1000.0)", argN)
		args = append(args, query.SinceMs)
		argN++
	}
	if query.UntilMs > 0 {
		sql += fmt.Sprintf(" AND timestamp <= to_timestamp($%d::double precision / 1000.0)", argN)
		args = append(args, query.UntilMs)
		argN++
	}
	if query.CursorMs > 0 {
		sql += fmt.Sprintf(" AND timestamp < to_timestamp($%d::double precision / 1000.0)", argN)
		args = append(args, query.CursorMs)
		argN++
	}
	sql += " ORDER BY timestamp DESC"
	sql += fmt.Sprintf(" LIMIT $%d", argN)
	args = append(args, limit)
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]SwarmEvent, 0, limit)
	for rows.Next() {
		var (
			eventType string
			payload   []byte
			tsRaw     float64
		)
		if err := rows.Scan(&eventType, &payload, &tsRaw); err != nil {
			return nil, err
		}
		parsed := map[string]any{}
		_ = json.Unmarshal(payload, &parsed)
		out = append(out, SwarmEvent{
			EventType: eventType,
			Payload:   parsed,
			TsUnixMs:  int64(tsRaw),
		})
	}
	return out, rows.Err()
}

func (r *PostgresRepository) Stats(ctx context.Context) (map[string]any, error) {
	var activeAgents int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM agents`).Scan(&activeAgents); err != nil {
		return nil, err
	}
	var queuedTasks int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tasks WHERE status = 'queued'`).Scan(&queuedTasks); err != nil {
		return nil, err
	}
	return map[string]any{
		"activeAgents": activeAgents,
		"queuedTasks":  queuedTasks,
		"mode":         "phase4-postgres",
	}, nil
}

func (r *PostgresRepository) AppendEvent(ctx context.Context, agentID string, eventType string, payload map[string]any) error {
	if _, err := uuid.Parse(agentID); err != nil {
		agentID = uuid.Nil.String()
	}
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO agent_events (agent_id, event_type, payload, timestamp) VALUES ($1, $2, $3, NOW())`,
		agentID,
		eventType,
		payload,
	)
	return err
}

func (r *PostgresRepository) AppendReplayMetric(ctx context.Context, metric ReplayMetricRecord) error {
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO replay_query_metrics
		(client_key, query_limit, result_count, latency_ms, rejected, filter_event_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
		metric.ClientKey,
		metric.QueryLimit,
		metric.ResultCount,
		metric.LatencyMs,
		metric.Rejected,
		metric.FilterEventType,
	)
	return err
}

func (r *PostgresRepository) GetReplayMetricSeries(ctx context.Context, limit int) ([]ReplayMetricPoint, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.pool.Query(
		ctx,
		`SELECT EXTRACT(EPOCH FROM created_at) * 1000, latency_ms, rejected
		 FROM replay_query_metrics
		 ORDER BY created_at DESC
		 LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]ReplayMetricPoint, 0, limit)
	for rows.Next() {
		var p ReplayMetricPoint
		var ts float64
		if err := rows.Scan(&ts, &p.LatencyMs, &p.Rejected); err != nil {
			return nil, err
		}
		p.CreatedAtMs = int64(ts)
		out = append(out, p)
	}
	return out, rows.Err()
}
