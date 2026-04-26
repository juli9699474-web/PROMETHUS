package router

import "sync"

type EventHub struct {
	mu          sync.RWMutex
	subscribers map[chan SwarmEvent]struct{}
}

func NewEventHub() *EventHub {
	return &EventHub{subscribers: map[chan SwarmEvent]struct{}{}}
}

func (h *EventHub) Subscribe() chan SwarmEvent {
	ch := make(chan SwarmEvent, 64)
	h.mu.Lock()
	h.subscribers[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *EventHub) Unsubscribe(ch chan SwarmEvent) {
	h.mu.Lock()
	delete(h.subscribers, ch)
	close(ch)
	h.mu.Unlock()
}

func (h *EventHub) Publish(evt SwarmEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subscribers {
		select {
		case ch <- evt:
		default:
		}
	}
}
