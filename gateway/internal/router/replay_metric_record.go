package router

type ReplayMetricRecord struct {
	ClientKey      string
	QueryLimit     int
	ResultCount    int
	LatencyMs      int64
	Rejected       bool
	FilterEventType string
}
