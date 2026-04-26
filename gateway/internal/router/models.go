package router

import "time"

type Agent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Fitness   float64   `json:"fitnessScore"`
	CreatedAt time.Time `json:"createdAt"`
}

type Task struct {
	ID          string    `json:"id"`
	AgentID     string    `json:"agentId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Reward      float64   `json:"reward"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SwarmEvent struct {
	EventType string         `json:"eventType"`
	Payload   map[string]any `json:"payload"`
	TsUnixMs  int64          `json:"tsUnixMs"`
}

type EventQuery struct {
	Limit     int
	EventType string
	AgentID   string
	SinceMs   int64
	UntilMs   int64
	CursorMs  int64
}
