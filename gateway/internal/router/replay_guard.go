package router

import (
	"sync"
	"time"
)

const (
	maxReplayLimit       = 500
	maxReplayRequestsMin = 120
)

type replayWindow struct {
	windowStart time.Time
	count       int
}

type ReplayRateGuard struct {
	mu      sync.Mutex
	clients map[string]replayWindow
}

func NewReplayRateGuard() *ReplayRateGuard {
	return &ReplayRateGuard{
		clients: map[string]replayWindow{},
	}
}

func (g *ReplayRateGuard) Allow(clientKey string, now time.Time) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	w := g.clients[clientKey]
	if w.windowStart.IsZero() || now.Sub(w.windowStart) >= time.Minute {
		g.clients[clientKey] = replayWindow{windowStart: now, count: 1}
		return true
	}
	if w.count >= maxReplayRequestsMin {
		return false
	}
	w.count++
	g.clients[clientKey] = w
	return true
}
