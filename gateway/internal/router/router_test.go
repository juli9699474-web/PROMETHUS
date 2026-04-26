package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSwarmStateEndpoint(t *testing.T) {
	deps := NewDependencies(NewSwarmStore())
	r := SetupRouter(deps)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/swarm/state", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestMetricsRequiresApiKeyWhenConfigured(t *testing.T) {
	t.Setenv("PROMETHEUS_API_KEY", "secret")
	deps := NewDependencies(NewSwarmStore())
	r := SetupRouter(deps)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/metrics/replay", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/metrics/replay", nil)
	req2.Header.Set("X-API-Key", os.Getenv("PROMETHEUS_API_KEY"))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
}
