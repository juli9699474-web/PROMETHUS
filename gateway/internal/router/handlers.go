package router

import (
	"context"
	"time"
)

type createAgentRequest struct {
	Name string `json:"name"`
}

type assignTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Reward      float64 `json:"reward"`
}

func (d *Dependencies) emit(eventType string, payload map[string]any) {
	_ = d.Repo.AppendEvent(context.Background(), extractAgentID(payload), eventType, payload)
	d.Hub.Publish(SwarmEvent{
		EventType: eventType,
		Payload:   payload,
		TsUnixMs:  time.Now().UTC().UnixMilli(),
	})
}

func extractAgentID(payload map[string]any) string {
	if raw, ok := payload["agent"]; ok {
		if agent, ok := raw.(Agent); ok {
			return agent.ID
		}
	}
	if raw, ok := payload["task"]; ok {
		if task, ok := raw.(Task); ok {
			return task.AgentID
		}
	}
	return "00000000-0000-0000-0000-000000000000"
}
