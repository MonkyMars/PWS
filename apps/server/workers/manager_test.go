package workers

import (
	"context"
	"testing"
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/types"
)

// createTestConfig creates a test configuration for workers
func createTestConfig() *config.Config {
	return &config.Config{
		AppName:     "test-app",
		Environment: "test",
		Port:        "8080",
		LogLevel:    "debug",
		Audit: types.AuditConfig{
			BatchSize:     10,
			ChannelSize:   100,
			Enabled:       true,
			FlushTime:     1 * time.Second,
			MaxFailures:   3,
			MaxRetries:    2,
			RetentionDays: 30,
			RetryDelay:    1 * time.Second,
		},
		Health: types.HealthConfig{
			BatchSize:      10,
			ChannelSize:    100,
			Enabled:        true,
			FlushTime:      2 * time.Second,
			ReportInterval: 1 * time.Second,
			MaxFailures:    3,
			MaxRetries:     2,
			RetentionDays:  7,
			RetryDelay:     1 * time.Second,
		},
	}
}

// createTestLogger creates a test logger
func createTestLogger() *config.Logger {
	// In a real test, you might want to use a test logger that captures output
	return config.SetupLogger()
}

func TestNewWorkerManager(t *testing.T) {
	cfg := createTestConfig()
	logger := createTestLogger()

	manager := NewWorkerManager(cfg, logger)

	if manager == nil {
		t.Fatal("NewWorkerManager returned nil")
	}

	if manager.cfg != cfg {
		t.Error("Config not properly set")
	}

	if manager.logger != logger {
		t.Error("Logger not properly set")
	}

	if manager.running {
		t.Error("Manager should not be running initially")
	}
}

func TestWorkerManagerStart(t *testing.T) {
	cfg := createTestConfig()
	logger := createTestLogger()
	manager := NewWorkerManager(cfg, logger)

	// Test starting the manager
	err := manager.Start()
	if err != nil {
		t.Fatalf("Failed to start worker manager: %v", err)
	}

	// Verify manager is running
	if !manager.running {
		t.Error("Manager should be running after Start()")
	}

	// Verify workers are initialized
	if manager.auditWorker == nil {
		t.Error("Audit worker should be initialized")
	}

	if manager.healthWorker == nil {
		t.Error("Health worker should be initialized")
	}

	if manager.cleanupWorker == nil {
		t.Error("Cleanup worker should be initialized")
	}

	// Test starting again should return error
	err = manager.Start()
	if err == nil {
		t.Error("Starting already running manager should return error")
	}

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = manager.Stop(ctx)
	if err != nil {
		t.Errorf("Failed to stop worker manager: %v", err)
	}
}

func TestWorkerManagerStop(t *testing.T) {
	cfg := createTestConfig()
	logger := createTestLogger()
	manager := NewWorkerManager(cfg, logger)

	// Start the manager first
	err := manager.Start()
	if err != nil {
		t.Fatalf("Failed to start worker manager: %v", err)
	}

	// Test stopping the manager
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = manager.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop worker manager: %v", err)
	}

	// Verify manager is not running
	if manager.running {
		t.Error("Manager should not be running after Stop()")
	}

	// Test stopping again should not error
	err = manager.Stop(ctx)
	if err != nil {
		t.Error("Stopping already stopped manager should not return error")
	}
}

func TestWorkerManagerStopTimeout(t *testing.T) {
	cfg := createTestConfig()
	logger := createTestLogger()
	manager := NewWorkerManager(cfg, logger)

	// Start the manager
	err := manager.Start()
	if err != nil {
		t.Fatalf("Failed to start worker manager: %v", err)
	}

	// Test stopping with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	err = manager.Stop(ctx)
	if err == nil {
		t.Error("Stop with very short timeout should return context deadline exceeded error")
	}

	// Clean up properly
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	err = manager.Stop(ctx2)
	if err == nil {
		t.Error("Stop with short timeout should return context deadline exceeded error")
	}
}

func TestWorkerManagerAddAuditLog(t *testing.T) {
	cfg := createTestConfig()
	logger := createTestLogger()
	manager := NewWorkerManager(cfg, logger)

	err := manager.Start()
	if err != nil {
		t.Fatalf("Failed to start worker manager: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = manager.Stop(ctx)
		if err != nil {
			t.Errorf("Failed to stop worker manager: %v", err)
		}
	}()

	// Test adding audit log
	auditLog := types.AuditLog{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "Test audit log",
		Attrs:     map[string]any{"test": "value"},
	}

	// This should not panic or error
	manager.AddAuditLog(auditLog)

	// Test adding audit log with empty message (should be ignored)
	emptyLog := types.AuditLog{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "",
	}

	manager.AddAuditLog(emptyLog)
}

func TestWorkerManagerRecordHealthMetric(t *testing.T) {
	cfg := createTestConfig()
	logger := createTestLogger()
	manager := NewWorkerManager(cfg, logger)

	err := manager.Start()
	if err != nil {
		t.Fatalf("Failed to start worker manager: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = manager.Stop(ctx)
		if err != nil {
			t.Errorf("Failed to stop worker manager: %v", err)
		}
	}()

	// Register a service first
	if manager.healthWorker != nil {
		manager.healthWorker.RegisterService("test-service")
	}

	// Test recording health metric
	manager.RecordHealthMetric("test-service", 200, 100*time.Millisecond)
	manager.RecordHealthMetric("test-service", 404, 50*time.Millisecond)
}

func TestWorkerManagerHealthStatus(t *testing.T) {
	cfg := createTestConfig()
	logger := createTestLogger()
	manager := NewWorkerManager(cfg, logger)

	// Test health status when not running
	status := manager.HealthStatus()
	if status["manager_running"].(bool) {
		t.Error("Manager should not be running initially")
	}

	// Start manager and test health status
	err := manager.Start()
	if err != nil {
		t.Fatalf("Failed to start worker manager: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = manager.Stop(ctx)
		if err != nil {
			t.Errorf("Failed to stop worker manager: %v", err)
		}
	}()

	status = manager.HealthStatus()
	if !status["manager_running"].(bool) {
		t.Error("Manager should be running after Start()")
	}

	// Check that all worker statuses are included
	if _, ok := status["audit"]; !ok {
		t.Error("Audit worker status should be included")
	}

	if _, ok := status["health"]; !ok {
		t.Error("Health worker status should be included")
	}

	if _, ok := status["cleanup"]; !ok {
		t.Error("Cleanup worker status should be included")
	}

	if _, ok := status["is_healthy"]; !ok {
		t.Error("Overall health status should be included")
	}
}

func TestWorkerManagerWithDisabledWorkers(t *testing.T) {
	cfg := createTestConfig()
	cfg.Audit.Enabled = false
	cfg.Health.Enabled = false

	logger := createTestLogger()
	manager := NewWorkerManager(cfg, logger)

	err := manager.Start()
	if err != nil {
		t.Fatalf("Failed to start worker manager: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = manager.Stop(ctx)
		if err != nil {
			t.Errorf("Failed to stop worker manager: %v", err)
		}
	}()

	// Adding audit log should be safe even when disabled
	auditLog := types.AuditLog{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "Test audit log",
	}
	manager.AddAuditLog(auditLog)

	// Recording health metric should be safe even when disabled
	manager.RecordHealthMetric("test-service", 200, 100*time.Millisecond)
}

func TestGetGlobalManager(t *testing.T) {
	// Test that global manager is created and reused
	manager1 := GetGlobalManager()
	manager2 := GetGlobalManager()

	if manager1 != manager2 {
		t.Error("GetGlobalManager should return the same instance")
	}

	if manager1 == nil {
		t.Error("GetGlobalManager should not return nil")
	}
}

func TestBackwardCompatibilityFunctions(t *testing.T) {
	// Test that backward compatibility functions don't panic
	// These functions use the global manager

	// Test starting individual workers
	StartAuditWorker()
	StartHealthLogWorker()
	StartCleanupScheduler()

	// Add some data
	auditLog := types.AuditLog{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "Test audit log",
	}
	AddAuditLog(auditLog)

	// Get health status
	auditHealth := AuditHealthStatus()
	if auditHealth == nil {
		t.Error("AuditHealthStatus should not return nil")
	}

	serviceHealth := ServiceHealthStatus()
	if serviceHealth == nil {
		t.Error("ServiceHealthStatus should not return nil")
	}

	// Test cleanup
	err := TriggerCleanupNow()
	if err == nil {
		t.Error("TriggerCleanupNow should return error when no cleanup worker")
	}

	// Stop workers
	StopAuditWorker()
	StopHealthLogWorker()
	StopAuditCleanupScheduler()
}

func TestWorkerManagerConcurrency(t *testing.T) {
	cfg := createTestConfig()
	logger := createTestLogger()
	manager := NewWorkerManager(cfg, logger)

	err := manager.Start()
	if err != nil {
		t.Fatalf("Failed to start worker manager: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = manager.Stop(ctx)
		if err != nil {
			t.Errorf("Failed to stop worker manager: %v", err)
		}
	}()

	// Test concurrent access to manager functions
	const numGoroutines = 10
	const numOperations = 100

	done := make(chan bool, numGoroutines)

	// Start multiple goroutines doing concurrent operations
	for i := range numGoroutines {
		go func(workerID int) {
			defer func() { done <- true }()

			for j := range numOperations {
				// Add audit logs
				auditLog := types.AuditLog{
					Timestamp: time.Now(),
					Level:     "info",
					Message:   "Concurrent test log",
					Attrs:     map[string]any{"worker": workerID, "operation": j},
				}
				manager.AddAuditLog(auditLog)

				// Record health metrics
				serviceName := "test-service"
				if manager.healthWorker != nil {
					manager.healthWorker.RegisterService(serviceName)
					manager.RecordHealthMetric(serviceName, 200, time.Duration(j)*time.Millisecond)
				}

				// Get health status
				_ = manager.HealthStatus()
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for range numGoroutines {
		select {
		case <-done:
			// Goroutine completed
		case <-time.After(10 * time.Second):
			t.Fatal("Timeout waiting for concurrent operations to complete")
		}
	}
}

func BenchmarkWorkerManagerAddAuditLog(b *testing.B) {
	cfg := createTestConfig()
	logger := createTestLogger()
	manager := NewWorkerManager(cfg, logger)

	err := manager.Start()
	if err != nil {
		b.Fatalf("Failed to start worker manager: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = manager.Stop(ctx)
		if err != nil {
			b.Errorf("Failed to stop worker manager: %v", err)
		}
	}()

	auditLog := types.AuditLog{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "Benchmark test log",
		Attrs:     map[string]any{"test": "benchmark"},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			manager.AddAuditLog(auditLog)
		}
	})
}

func BenchmarkWorkerManagerHealthStatus(b *testing.B) {
	cfg := createTestConfig()
	logger := createTestLogger()
	manager := NewWorkerManager(cfg, logger)

	err := manager.Start()
	if err != nil {
		b.Fatalf("Failed to start worker manager: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = manager.Stop(ctx)
		if err != nil {
			b.Errorf("Failed to stop worker manager: %v", err)
		}
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = manager.HealthStatus()
		}
	})
}
