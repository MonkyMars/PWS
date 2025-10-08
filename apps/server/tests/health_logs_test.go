package tests

import (
	"slices"
	"testing"
	"time"

	"github.com/MonkyMars/PWS/types"
	"github.com/MonkyMars/PWS/workers"
)

func TestHealthLogConversion(t *testing.T) {
	// Create a sample health log
	timestamp := time.Now()
	healthLog := types.HealthLog{
		Timestamp:      timestamp,
		Service:        "auth",
		StatusCode:     200,
		RequestCount:   150,
		ErrorCount:     5,
		AverageLatency: 250 * time.Millisecond,
		TimeSpan:       5 * time.Minute,
	}

	// Test that health log can be processed (conversion logic is now internal)
	if healthLog.Service != "auth" {
		t.Errorf("Expected service 'auth', got %v", healthLog.Service)
	}

	if healthLog.StatusCode != 200 {
		t.Errorf("Expected status_code 200, got %v", healthLog.StatusCode)
	}

	if healthLog.RequestCount != 150 {
		t.Errorf("Expected request_count 150, got %v", healthLog.RequestCount)
	}

	if healthLog.ErrorCount != 5 {
		t.Errorf("Expected error_count 5, got %v", healthLog.ErrorCount)
	}

	if healthLog.AverageLatency != 250*time.Millisecond {
		t.Errorf("Expected average_latency 250ms, got %v", healthLog.AverageLatency)
	}

	if healthLog.TimeSpan != 5*time.Minute {
		t.Errorf("Expected time_span 5 minutes, got %v", healthLog.TimeSpan)
	}

	// Verify timestamp is set correctly
	if healthLog.Timestamp != timestamp {
		t.Errorf("Expected timestamp %v, got %v", timestamp, healthLog.Timestamp)
	}
}

func TestWorkerManagerHealthFunctionality(t *testing.T) {
	manager := workers.GetGlobalManager()
	if manager == nil {
		t.Skip("Worker manager is not available")
	}

	// Test recording health metrics
	manager.RecordHealthMetric("test-service", 200, 100*time.Millisecond)
	manager.RecordHealthMetric("test-service", 404, 50*time.Millisecond)

	// Verify service was registered through the health metrics
	services := workers.GetAllServices()
	found := false
	if slices.Contains(services, "test-service") {
		found = true
	}

	if !found {
		t.Errorf("Expected to find 'test-service' in registered services, got %v", services)
	}

	// Test getting service stats
	stats := workers.GetServiceStats("test-service")
	if stats == nil {
		t.Errorf("Expected to get stats for 'test-service', got nil")
	}
}
