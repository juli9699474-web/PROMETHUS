package router

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	CreateAgent(ctx context.Context, name string) (Agent, error)
	GetAgent(ctx context.Context, id string) (Agent, error)
	ListAgents(ctx context.Context, limit int) ([]Agent, error)
	CreateTask(ctx context.Context, agentID, title, description string, reward float64) (Task, error)
	ListTasks(ctx context.Context, limit int) ([]Task, error)
	ListEvents(ctx context.Context, limit int) ([]SwarmEvent, error)
	QueryEvents(ctx context.Context, query EventQuery) ([]SwarmEvent, error)
	Stats(ctx context.Context) (map[string]any, error)
	AppendEvent(ctx context.Context, agentID string, eventType string, payload map[string]any) error
	AppendReplayMetric(ctx context.Context, metric ReplayMetricRecord) error
	GetReplayMetricSeries(ctx context.Context, limit int) ([]ReplayMetricPoint, error)
}
