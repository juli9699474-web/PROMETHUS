package router

type ReplayMetricPoint struct {
	CreatedAtMs int64 `json:"createdAtMs"`
	LatencyMs   int64 `json:"latencyMs"`
	Rejected    bool  `json:"rejected"`
}
