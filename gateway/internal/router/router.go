package router

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Repo        Repository
	Hub         *EventHub
	ReplayGuard *ReplayRateGuard
	Metrics     *ReplayMetrics
}

func NewDependencies(repo Repository) *Dependencies {
	return &Dependencies{
		Repo:        repo,
		Hub:         NewEventHub(),
		ReplayGuard: NewReplayRateGuard(),
		Metrics:     &ReplayMetrics{},
	}
}

func SetupRouter(deps *Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/swarm/state", func(c *gin.Context) {
			stats, err := deps.Repo.Stats(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load swarm state"})
				return
			}
			c.JSON(http.StatusOK, stats)
		})
		v1.POST("/agents", func(c *gin.Context) {
			var req createAgentRequest
			if err := c.BindJSON(&req); err != nil || req.Name == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			agent, err := deps.Repo.CreateAgent(c.Request.Context(), req.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create agent"})
				return
			}
			deps.emit("agent.created", map[string]any{"agent": agent})
			c.JSON(http.StatusCreated, agent)
		})
		v1.GET("/agents/:id", func(c *gin.Context) {
			agent, err := deps.Repo.GetAgent(c.Request.Context(), c.Param("id"))
			if errors.Is(err, ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
				return
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch agent"})
				return
			}
			c.JSON(http.StatusOK, agent)
		})
		v1.GET("/agents", func(c *gin.Context) {
			agents, err := deps.Repo.ListAgents(c.Request.Context(), parseLimit(c.Query("limit"), 100))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list agents"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"items": agents})
		})
		v1.POST("/agents/:id/tasks", func(c *gin.Context) {
			agentID := c.Param("id")
			if _, err := deps.Repo.GetAgent(c.Request.Context(), agentID); errors.Is(err, ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate agent"})
				return
			}
			var req assignTaskRequest
			if err := c.BindJSON(&req); err != nil || req.Title == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
			reward := req.Reward
			if reward <= 0 {
				reward = 10.0
			}
			task, err := deps.Repo.CreateTask(c.Request.Context(), agentID, req.Title, req.Description, reward)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
				return
			}
			deps.emit("task.assigned", map[string]any{"task": task})
			c.JSON(http.StatusCreated, task)
		})
		v1.GET("/tasks", func(c *gin.Context) {
			tasks, err := deps.Repo.ListTasks(c.Request.Context(), parseLimit(c.Query("limit"), 100))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tasks"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"items": tasks})
		})
		v1.GET("/events", func(c *gin.Context) {
			clientKey := c.GetHeader("X-API-Key")
			if clientKey == "" {
				clientKey = c.ClientIP()
			}
			if !deps.ReplayGuard.Allow(clientKey, time.Now().UTC()) {
				deps.Metrics.Reject()
				_ = deps.Repo.AppendReplayMetric(c.Request.Context(), ReplayMetricRecord{
					ClientKey: clientKey, QueryLimit: 0, ResultCount: 0, LatencyMs: 0, Rejected: true,
				})
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "replay query rate limit exceeded"})
				return
			}
			start := time.Now()
			query := EventQuery{
				Limit:     boundedLimit(c.Query("limit"), 200),
				EventType: c.Query("eventType"),
				AgentID:   c.Query("agentId"),
				SinceMs:   parseInt64(c.Query("sinceMs"), 0),
				UntilMs:   parseInt64(c.Query("untilMs"), 0),
				CursorMs:  parseInt64(c.Query("cursorMs"), 0),
			}
			if err := validateEventQuery(query); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			events, err := deps.Repo.QueryEvents(c.Request.Context(), query)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
				return
			}
			nextCursor := int64(0)
			if len(events) > 0 {
				nextCursor = events[len(events)-1].TsUnixMs
			}
			latency := time.Since(start)
			deps.Metrics.Observe(len(events), latency)
			_ = deps.Repo.AppendReplayMetric(c.Request.Context(), ReplayMetricRecord{
				ClientKey: clientKey, QueryLimit: query.Limit, ResultCount: len(events), LatencyMs: latency.Milliseconds(),
				Rejected: false, FilterEventType: query.EventType,
			})
			c.JSON(http.StatusOK, gin.H{"items": events, "nextCursorMs": nextCursor})
		})
		protected := v1.Group("")
		protected.Use(apiKeyMiddleware())
		protected.GET("/metrics/replay", func(c *gin.Context) {
			c.JSON(http.StatusOK, deps.Metrics.Snapshot())
		})
		protected.GET("/metrics/replay/timeseries", func(c *gin.Context) {
			limit := parseLimit(c.Query("limit"), 100)
			points, err := deps.Repo.GetReplayMetricSeries(c.Request.Context(), limit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load replay timeseries"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"items": points})
		})
		protected.GET("/metrics/replay/percentiles", func(c *gin.Context) {
			limit := parseLimit(c.Query("limit"), 500)
			points, err := deps.Repo.GetReplayMetricSeries(c.Request.Context(), limit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load replay percentiles"})
				return
			}
			latencies := make([]int64, 0, len(points))
			for _, p := range points {
				if p.Rejected {
					continue
				}
				latencies = append(latencies, p.LatencyMs)
			}
			c.JSON(http.StatusOK, CalculatePercentiles(latencies))
		})
	}

	r.POST("/uap/connect", func(c *gin.Context) { c.JSON(200, gin.H{"status": "connected"}) })
	r.POST("/uap/message", func(c *gin.Context) { c.JSON(200, gin.H{"status": "accepted"}) })
	r.GET("/ws/swarm", swarmWebSocketHandler(deps))

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			deps.emit("swarm.heartbeat", map[string]any{"status": "alive"})
		}
	}()

	return r
}
