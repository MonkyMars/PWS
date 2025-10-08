package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
)

// Start starts the cleanup worker
func (cw *CleanupWorker) Start() error {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	if cw.running {
		return fmt.Errorf("cleanup worker already running")
	}

	if !cw.cfg.Audit.Enabled || cw.cfg.Audit.RetentionDays <= 0 {
		return nil // No cleanup needed
	}

	cw.running = true
	cw.wg.Add(1)
	go cw.run()

	return nil
}

// Stop gracefully stops the cleanup worker
func (cw *CleanupWorker) Stop(ctx context.Context) error {
	cw.mu.Lock()
	if !cw.running {
		cw.mu.Unlock()
		return nil
	}
	cw.cancel()
	cw.mu.Unlock()

	// Wait for worker to finish with timeout
	done := make(chan struct{})
	go func() {
		cw.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		cw.logger.Info("Cleanup worker stopped successfully")
		return nil
	case <-ctx.Done():
		cw.logger.Warn("Cleanup worker stop timed out")
		return ctx.Err()
	}
}

// TriggerCleanup manually triggers a cleanup operation
func (cw *CleanupWorker) TriggerCleanup() error {
	cw.logger.Info("Manual cleanup triggered")

	if err := cw.cleanupOldAuditLogs(); err != nil {
		cw.logger.Error("Manual cleanup failed", "error", err)
		return err
	}

	cw.logger.Info("Manual cleanup completed successfully")
	return nil
}

// HealthStatus returns the current health status of the cleanup worker
func (cw *CleanupWorker) HealthStatus() map[string]any {
	if cw == nil {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "cleanup worker is nil",
		}
	}

	if cw.cfg == nil {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "cleanup worker configuration is nil",
		}
	}

	cw.mu.RLock()
	defer cw.mu.RUnlock()

	enabled := cw.cfg.Audit.Enabled && cw.cfg.Audit.RetentionDays > 0
	isHealthy := enabled && cw.running

	return map[string]any{
		"enabled":        enabled,
		"worker_running": cw.running,
		"is_healthy":     isHealthy,
		"configuration": map[string]any{
			"retention_days": cw.cfg.Audit.RetentionDays,
		},
	}
}

// run is the main cleanup worker loop
func (cw *CleanupWorker) run() {
	defer cw.wg.Done()
	defer func() {
		cw.mu.Lock()
		cw.running = false
		cw.mu.Unlock()
	}()

	cw.logger.Info("Starting audit log cleanup scheduler")

	for {
		// Calculate time until next midnight (00:00)
		now := time.Now()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		duration := nextMidnight.Sub(now)

		// Wait until midnight or context cancellation
		select {
		case <-time.After(duration):
			// Run cleanup
			if err := cw.cleanupOldAuditLogs(); err != nil {
				cw.logger.Error("Scheduled cleanup failed", "error", err)
			} else {
				cw.logger.Info("Scheduled cleanup completed successfully")
			}
		case <-cw.ctx.Done():
			cw.logger.Info("Cleanup scheduler stopped")
			return
		}
	}
}

// cleanupOldAuditLogs removes audit logs older than the retention period
func (cw *CleanupWorker) cleanupOldAuditLogs() error {
	if !cw.cfg.Audit.Enabled || cw.cfg.Audit.RetentionDays <= 0 {
		return nil // No cleanup needed
	}

	cutoff := time.Now().AddDate(0, 0, -cw.cfg.Audit.RetentionDays)
	query := services.Query().
		SetOperation("delete").
		SetTable("audit_logs").
		SetWhereRaw("audit_logs.timestamp < ?", cutoff)

	result, err := database.ExecuteQuery[types.AuditLog](query)
	if err != nil {
		cw.logger.Error("Failed to clean up old audit logs", "error", err)
		return fmt.Errorf("cleanup failed: %w", err)
	}

	cw.logger.Info("Cleaned up old audit logs", "deleted_count", result.Count)
	return nil
}

// Backward compatibility function
func CleanupOldAuditLogs() error {
	manager := GetGlobalManager()
	if manager.cleanupWorker != nil {
		return manager.cleanupWorker.cleanupOldAuditLogs()
	}
	return fmt.Errorf("cleanup worker not available")
}
