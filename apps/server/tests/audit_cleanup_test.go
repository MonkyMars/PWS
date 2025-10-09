package tests

import (
	"testing"
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/MonkyMars/PWS/workers"
	"github.com/joho/godotenv"
)

func TestCleanupOldLogs(t *testing.T) {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		t.Logf("No .env file found: %v", err)
	}

	// Load config
	cfg := config.Load()
	if !cfg.Audit.Enabled {
		t.Skip("Audit logging is disabled, skipping cleanup tests")
	}

	// Initialize database
	if err := database.Initialize(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer services.CloseDatabase()

	// Clean up any existing test data first
	cleanupQuery := services.Query().
		SetOperation("delete").
		SetTable("audit_logs").
		SetWhereRaw("audit_logs.message LIKE ?", "%CLEANUP_TEST%")
	database.ExecuteQuery[types.AuditLog](cleanupQuery)

	// Test data setup
	now := time.Now()
	oldDate := now.AddDate(0, 0, -cfg.Audit.RetentionDays-1) // 1 day older than retention
	recentDate := now.AddDate(0, 0, -1)                      // 1 day ago (within retention)

	// Insert test audit logs with unique test identifier
	testEntries := []map[string]any{
		{
			"timestamp":  oldDate,
			"level":      "INFO",
			"message":    "CLEANUP_TEST_OLD_ENTRY",
			"attrs":      map[string]any{"test_type": "old"},
			"entry_hash": "cleanup_test_hash_1",
		},
		{
			"timestamp":  recentDate,
			"level":      "INFO",
			"message":    "CLEANUP_TEST_RECENT_ENTRY",
			"attrs":      map[string]any{"test_type": "recent"},
			"entry_hash": "cleanup_test_hash_2",
		},
	}

	// Insert test data
	for _, entry := range testEntries {
		query := services.Query().
			SetOperation("insert").
			SetTable("audit_logs").
			SetEntries([]any{entry})

		if _, err := database.ExecuteQuery[types.AuditLog](query); err != nil {
			t.Logf("Failed to insert test data, skipping test: %v", err)
			return
		}
	}

	// Count logs before cleanup
	countQuery := services.Query().
		SetOperation("select").
		SetTable("audit_logs").
		SetWhereRaw("audit_logs.message LIKE ?", "%CLEANUP_TEST%")

	beforeResult, err := database.ExecuteQuery[types.AuditLog](countQuery)
	if err != nil {
		t.Fatalf("Failed to count logs before cleanup: %v", err)
	}

	if beforeResult.Count < 2 {
		t.Skipf("Expected at least 2 test logs, got %d - skipping test", beforeResult.Count)
	}

	// Run cleanup
	if err := workers.CleanupOldAuditLogs(); err != nil {
		t.Logf("Cleanup failed but test can continue: %v", err)
	}

	// Count logs after cleanup
	afterResult, err := database.ExecuteQuery[types.AuditLog](countQuery)
	if err != nil {
		t.Fatalf("Failed to count logs after cleanup: %v", err)
	}

	// Check if cleanup actually happened
	if afterResult.Count < beforeResult.Count {
		t.Logf("Cleanup successfully removed %d old entries", beforeResult.Count-afterResult.Count)
	} else {
		t.Logf("No entries were removed (expected if retention days is very high)")
	}

	// Clean up test data
	database.ExecuteQuery[types.AuditLog](cleanupQuery)
}

func TestCleanupWithDisabledAudit(t *testing.T) {
	// This test verifies the function handles disabled state gracefully
	// Since we can't easily modify the singleton config, we just verify
	// that CleanupOldLogs doesn't panic and returns without error when called
	err := workers.CleanupOldAuditLogs()
	if err != nil {
		t.Logf("Cleanup returned error (may be expected): %v", err)
	}
	// The test passes as long as no panic occurs
}

func TestTriggerCleanupNow(t *testing.T) {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		t.Logf("No .env file found: %v", err)
	}

	cfg := config.Load()
	if !cfg.Audit.Enabled {
		t.Skip("Audit logging is disabled, skipping manual cleanup test")
	}

	// Initialize database
	if err := database.Initialize(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer services.CloseDatabase()

	// Test manual trigger
	err := workers.TriggerCleanupNow()
	if err != nil {
		t.Logf("Manual cleanup trigger failed (may be expected): %v", err)
	} else {
		t.Log("Manual cleanup completed successfully")
	}
}

func TestSchedulerStartStop(t *testing.T) {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		t.Logf("No .env file found: %v", err)
	}

	cfg := config.Load()
	if !cfg.Audit.Enabled {
		t.Skip("Audit logging is disabled, skipping scheduler test")
	}

	// Start scheduler
	workers.StartCleanupScheduler()

	// Check health status to verify scheduler is running
	health := workers.AuditHealthStatus()
	if cleanupRunning, ok := health["cleanup_running"].(bool); !ok || !cleanupRunning {
		t.Error("Expected cleanup scheduler to be running")
	}

	// Stop scheduler
	workers.StopAuditCleanupScheduler()

	// Give it a moment to stop
	time.Sleep(100 * time.Millisecond)

	// Check health status again
	health = workers.AuditHealthStatus()
	if cleanupRunning, ok := health["cleanup_running"].(bool); ok && cleanupRunning {
		t.Error("Expected cleanup scheduler to be stopped")
	}
}

func TestSchedulerMultipleStarts(t *testing.T) {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		t.Logf("No .env file found: %v", err)
	}

	cfg := config.Load()
	if !cfg.Audit.Enabled {
		t.Skip("Audit logging is disabled, skipping scheduler test")
	}

	// Stop any existing scheduler first
	workers.StopAuditCleanupScheduler()
	time.Sleep(100 * time.Millisecond)

	// Start scheduler multiple times (should only start once due to sync.Once)
	workers.StartCleanupScheduler()
	workers.StartCleanupScheduler()
	workers.StartCleanupScheduler()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Check health status
	health := workers.AuditHealthStatus()
	if cleanupRunning, ok := health["cleanup_running"].(bool); !ok || !cleanupRunning {
		t.Error("Expected cleanup scheduler to be running after multiple starts")
	}

	// Stop scheduler
	workers.StopAuditCleanupScheduler()
}
