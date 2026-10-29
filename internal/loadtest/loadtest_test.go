package loadtest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRunAgainstLocalServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result, err := Run(context.Background(), Config{
		URL:         server.URL,
		Connections: 4,
		Duration:    500 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if result.TotalRequests == 0 {
		t.Fatal("expected at least one request to complete")
	}
	if result.Non2xx != 0 {
		t.Errorf("expected all requests to succeed, got %d non-2xx", result.Non2xx)
	}
	if result.AverageLatencyMs <= 0 {
		t.Errorf("expected positive average latency, got %f", result.AverageLatencyMs)
	}
	if result.AverageRequestsPerSec < 0 {
		t.Errorf("expected non-negative average rps, got %f", result.AverageRequestsPerSec)
	}
}

func TestRunAgainstUnreachableServerCountsAllAsFailures(t *testing.T) {
	result, err := Run(context.Background(), Config{
		URL:         "http://127.0.0.1:1",
		Connections: 2,
		Duration:    200 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.TotalRequests == 0 {
		t.Fatal("expected connection attempts to be counted even when they fail")
	}
	if result.Non2xx != result.TotalRequests {
		t.Errorf("expected every attempted request to count as non-2xx, got %d/%d", result.Non2xx, result.TotalRequests)
	}
}
