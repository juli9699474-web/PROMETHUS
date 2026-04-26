package router

import (
	"sort"
	"sync"
	"time"
)

type ReplayMetrics struct {
	mu                  sync.Mutex
	TotalQueries        int64         `json:"totalQueries"`
	RejectedQueries     int64         `json:"rejectedQueries"`
	TotalReturnedEvents int64         `json:"totalReturnedEvents"`
	TotalLatencyMs      int64         `json:"totalLatencyMs"`
	LastQueryAt         time.Time     `json:"lastQueryAt"`
	AvgLatencyMs        float64       `json:"avgLatencyMs"`
}

type ReplayPercentiles struct {
	P50 float64 `json:"p50"`
	P95 float64 `json:"p95"`
	P99 float64 `json:"p99"`
}

func (m *ReplayMetrics) Observe(resultCount int, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalQueries++
	m.TotalReturnedEvents += int64(resultCount)
	m.TotalLatencyMs += latency.Milliseconds()
	m.LastQueryAt = time.Now().UTC()
	if m.TotalQueries > 0 {
		m.AvgLatencyMs = float64(m.TotalLatencyMs) / float64(m.TotalQueries)
	}
}

func (m *ReplayMetrics) Reject() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RejectedQueries++
}

func (m *ReplayMetrics) Snapshot() ReplayMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	return ReplayMetrics{
		TotalQueries:        m.TotalQueries,
		RejectedQueries:     m.RejectedQueries,
		TotalReturnedEvents: m.TotalReturnedEvents,
		TotalLatencyMs:      m.TotalLatencyMs,
		LastQueryAt:         m.LastQueryAt,
		AvgLatencyMs:        m.AvgLatencyMs,
	}
}

func CalculatePercentiles(latencies []int64) ReplayPercentiles {
	if len(latencies) == 0 {
		return ReplayPercentiles{}
	}
	copyVals := make([]int64, len(latencies))
	copy(copyVals, latencies)
	sort.Slice(copyVals, func(i, j int) bool { return copyVals[i] < copyVals[j] })
	return ReplayPercentiles{
		P50: percentile(copyVals, 0.50),
		P95: percentile(copyVals, 0.95),
		P99: percentile(copyVals, 0.99),
	}
}

func percentile(sorted []int64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * p)
	return float64(sorted[idx])
}
