package router

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SwarmStore struct {
	mu     sync.RWMutex
	agents map[string]Agent
	tasks  map[string]Task
}

func NewSwarmStore() *SwarmStore {
	return &SwarmStore{
		agents: map[string]Agent{},
		tasks:  map[string]Task{},
	}
}

func (s *SwarmStore) createAgent(name string) Agent {
	s.mu.Lock()
	defer s.mu.Unlock()
	a := Agent{
		ID:        uuid.NewString(),
		Name:      name,
		Status:    "idle",
		Fitness:   0.5,
		CreatedAt: time.Now().UTC(),
	}
	s.agents[a.ID] = a
	return a
}

func (s *SwarmStore) getAgent(id string) (Agent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.agents[id]
	return a, ok
}

func (s *SwarmStore) createTask(agentID, title, description string, reward float64) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := Task{
		ID:          uuid.NewString(),
		AgentID:     agentID,
		Title:       title,
		Description: description,
		Status:      "queued",
		Reward:      reward,
		CreatedAt:   time.Now().UTC(),
	}
	s.tasks[t.ID] = t
	return t
}

func (s *SwarmStore) stats() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]any{
		"activeAgents": len(s.agents),
		"queuedTasks":  len(s.tasks),
		"mode":         "phase3-prototype",
	}
}

func (s *SwarmStore) AppendEvent(_ context.Context, _ string, _ string, _ map[string]any) error {
	return nil
}

func (s *SwarmStore) CreateAgent(ctx context.Context, name string) (Agent, error) {
	_ = ctx
	return s.createAgent(name), nil
}

func (s *SwarmStore) GetAgent(ctx context.Context, id string) (Agent, error) {
	_ = ctx
	a, ok := s.getAgent(id)
	if !ok {
		return Agent{}, ErrNotFound
	}
	return a, nil
}

func (s *SwarmStore) CreateTask(ctx context.Context, agentID, title, description string, reward float64) (Task, error) {
	_ = ctx
	return s.createTask(agentID, title, description, reward), nil
}

func (s *SwarmStore) ListAgents(ctx context.Context, limit int) ([]Agent, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()
	if limit <= 0 {
		limit = 100
	}
	out := make([]Agent, 0, limit)
	for _, a := range s.agents {
		out = append(out, a)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (s *SwarmStore) ListTasks(ctx context.Context, limit int) ([]Task, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()
	if limit <= 0 {
		limit = 100
	}
	out := make([]Task, 0, limit)
	for _, t := range s.tasks {
		out = append(out, t)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (s *SwarmStore) ListEvents(ctx context.Context, limit int) ([]SwarmEvent, error) {
	_ = ctx
	_ = limit
	return []SwarmEvent{}, nil
}

func (s *SwarmStore) QueryEvents(ctx context.Context, query EventQuery) ([]SwarmEvent, error) {
	_ = ctx
	_ = query
	return []SwarmEvent{}, nil
}

func (s *SwarmStore) Stats(ctx context.Context) (map[string]any, error) {
	_ = ctx
	return s.stats(), nil
}

func (s *SwarmStore) AppendReplayMetric(ctx context.Context, metric ReplayMetricRecord) error {
	_ = ctx
	_ = metric
	return nil
}

func (s *SwarmStore) GetReplayMetricSeries(ctx context.Context, limit int) ([]ReplayMetricPoint, error) {
	_ = ctx
	_ = limit
	return []ReplayMetricPoint{}, nil
}
