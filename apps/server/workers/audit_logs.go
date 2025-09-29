package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
)

var (
	auditChan      chan types.AuditLog
	auditCtx       context.Context
	auditCancel    context.CancelFunc
	auditWg        sync.WaitGroup
	once           sync.Once
	workerRunning  bool
	workerMutex    sync.RWMutex
	lastFlushTime  time.Time
	failureCount   int
	totalProcessed int64
	totalDropped   int64
	cleanupCtx     context.Context
	cleanupCancel  context.CancelFunc
	cleanupRunning bool
)

func StartAuditWorker() {
	once.Do(func() {
		cfg := config.Get()

		// Check if audit logging is enabled
		if !cfg.Audit.Enabled {
			return
		}

		workerMutex.Lock()
		defer workerMutex.Unlock()

		auditCtx, auditCancel = context.WithCancel(context.Background())
		auditChan = make(chan types.AuditLog, cfg.Audit.ChannelSize)
		workerRunning = true
		failureCount = 0

		auditWg.Go(func() {
			defer func() {
				workerMutex.Lock()
				workerRunning = false
				workerMutex.Unlock()
			}()
			runAuditWorker(auditCtx)
		})

	})
}

func runAuditWorker(ctx context.Context) {
	cfg := config.Get()
	batch := make([]types.AuditLog, 0, cfg.Audit.BatchSize)
	ticker := time.NewTicker(cfg.Audit.FlushTime)
	defer ticker.Stop()

	for {
		select {
		case entry := <-auditChan:
			batch = append(batch, entry)

			// flush immediately if batch full
			if len(batch) >= cfg.Audit.BatchSize {
				flushBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				flushBatch(batch)
				batch = batch[:0]
			}

		case <-ctx.Done():
			// Flush remaining entries before shutting down
			if len(batch) > 0 {
				flushBatch(batch)
			}
			// Drain any remaining entries in channel
			for {
				select {
				case entry := <-auditChan:
					batch = append(batch, entry)
					if len(batch) >= cfg.Audit.BatchSize {
						flushBatch(batch)
						batch = batch[:0]
					}
				default:
					if len(batch) > 0 {
						flushBatch(batch)
					}
					return
				}
			}
		}
	}
}

func flushBatch(entries []types.AuditLog) {
	logger := config.SetupLogger()
	cfg := config.Get()

	if len(entries) == 0 {
		return
	}

	var err error
	var successfulInserts int64

	for attempt := 0; attempt < cfg.Audit.MaxRetries; attempt++ {
		successfulInserts, err = tryFlushBatchWithCount(entries)
		if err == nil {
			workerMutex.Lock()
			failureCount = 0 // Reset failure count on success
			lastFlushTime = time.Now()
			totalProcessed += successfulInserts
			workerMutex.Unlock()

			logger.Debug("Flushed audit log batch",
				"count", len(entries),
				"successful_inserts", successfulInserts,
				"attempt", attempt+1)
			return
		}

		// Log retry attempt
		logger.Warn("Audit batch flush failed, retrying",
			"attempt", attempt+1,
			"max_retries", cfg.Audit.MaxRetries,
			"error", err,
			"batch_size", len(entries))

		if attempt < cfg.Audit.MaxRetries-1 {
			// Exponential backoff: 100ms, 200ms, 400ms
			backoff := time.Duration(100*(1<<attempt)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	// After all retries failed, update failure count
	workerMutex.Lock()
	failureCount++
	workerMutex.Unlock()

	// After all retries failed, log the error but don't crash
	logger.Error("Failed to flush audit log batch after retries",
		"error", err,
		"batch_size", len(entries),
		"max_retries", cfg.Audit.MaxRetries,
		"total_failures", failureCount)
}

func tryFlushBatchWithCount(entries []types.AuditLog) (int64, error) {
	// Convert AuditLog entries to the format expected by SetEntries
	auditEntries := make([]any, 0, len(entries))
	skippedEntries := 0

	for _, entry := range entries {
		// Validate entry before adding
		if entry.Message == "" {
			skippedEntries++
			continue // Skip invalid entries
		}

		auditEntry := map[string]any{
			"timestamp":  entry.Timestamp,
			"level":      entry.Level,
			"message":    entry.Message,
			"attrs":      entry.Attrs,
			"entry_hash": entry.EntryHash,
		}

		auditEntries = append(auditEntries, auditEntry)
	}

	if len(auditEntries) == 0 {
		return 0, nil // Nothing to flush
	}

	query := services.Query().
		SetOperation("insert").
		SetTable("audit_logs").
		SetEntries(auditEntries)

	result, err := database.ExecuteQuery[types.AuditLog](query)
	if err != nil {
		return 0, fmt.Errorf("database insert failed: %w", err)
	}

	// Return the actual number of rows inserted (may be less than auditEntries due to duplicates)
	successfulInserts := result.Count
	if skippedEntries > 0 {
		logger := config.SetupLogger()
		logger.Debug("Skipped invalid audit entries during flush",
			"skipped_count", skippedEntries,
			"total_entries", len(entries))
	}

	return successfulInserts, nil
}

// AddAuditLog adds an audit log entry to the processing queue
func AddAuditLog(entry types.AuditLog) {
	cfg := config.Get()
	logger := config.SetupLogger()

	// Check if audit logging is enabled
	if !cfg.Audit.Enabled {
		return
	}

	// Validate input
	if entry.Message == "" {
		return // Ignore invalid entries
	}

	// Set timestamp if not provided
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	// Check if worker is running
	workerMutex.RLock()
	running := workerRunning
	failures := failureCount
	workerMutex.RUnlock()

	if !running {
		logger.Warn("Audit worker not running, dropping log entry",
			"level", entry.Level,
			"message", entry.Message)
		return
	}

	// Circuit breaker: if too many failures, drop new entries temporarily
	if failures >= cfg.Audit.MaxFailures {
		workerMutex.Lock()
		totalDropped++
		workerMutex.Unlock()

		logger.Warn("Audit worker in failure mode, dropping log entry",
			"level", entry.Level,
			"message", entry.Message,
			"failure_count", failures)
		return
	}

	select {
	case auditChan <- entry:
		// Successfully added to queue
	default:
		// Channel is full, update dropped count and log warning
		workerMutex.Lock()
		totalDropped++
		workerMutex.Unlock()

		logger.Warn("Audit log channel is full, dropping log entry",
			"level", entry.Level,
			"message", entry.Message,
			"queue_size", len(auditChan))
	}
}

// StopAuditWorker gracefully stops the audit worker
func StopAuditWorker() {
	if auditCancel != nil {
		auditCancel()
		auditWg.Wait() // Wait for worker to finish

		workerMutex.Lock()
		if auditChan != nil {
			close(auditChan)
			auditChan = nil
		}
		workerRunning = false
		workerMutex.Unlock()
	}

	// Also stop the cleanup scheduler
	StopCleanupScheduler()
}

// HealthStatus returns the current health status of the audit worker
func HealthStatus() map[string]any {
	cfg := config.Get()
	workerMutex.RLock()
	defer workerMutex.RUnlock()

	queueSize := 0
	if auditChan != nil {
		queueSize = len(auditChan)
	}

	return map[string]any{
		"enabled":         cfg.Audit.Enabled,
		"worker_running":  workerRunning,
		"queue_size":      queueSize,
		"queue_capacity":  cfg.Audit.ChannelSize,
		"last_flush_time": lastFlushTime,
		"failure_count":   failureCount,
		"total_processed": totalProcessed,
		"total_dropped":   totalDropped,
		"is_healthy":      cfg.Audit.Enabled && workerRunning && failureCount < cfg.Audit.MaxFailures,
		"cleanup_running": cleanupRunning,
		"configuration": map[string]any{
			"batch_size":     cfg.Audit.BatchSize,
			"flush_time":     cfg.Audit.FlushTime.String(),
			"max_retries":    cfg.Audit.MaxRetries,
			"max_failures":   cfg.Audit.MaxFailures,
			"retention_days": cfg.Audit.RetentionDays,
		},
	}
}

// ResetFailureCount resets the failure count (for operational recovery)
func ResetFailureCount() {
	workerMutex.Lock()
	defer workerMutex.Unlock()
	failureCount = 0
}

func CleanupOldLogs() error {
	cfg := config.Get()
	logger := config.SetupLogger()

	if !cfg.Audit.Enabled || cfg.Audit.RetentionDays <= 0 {
		return nil // No cleanup needed
	}

	cutoff := time.Now().AddDate(0, 0, -cfg.Audit.RetentionDays)
	query := services.Query().
		SetOperation("delete").
		SetTable("audit_logs").
		SetWhereRaw("audit_logs.timestamp < ?", cutoff)
	result, err := database.ExecuteQuery[types.AuditLog](query)
	if err != nil {
		logger.Error("Failed to clean up old audit logs", "error", err)
		return fmt.Errorf("cleanup failed: %w", err)
	}

	logger.Info("Cleaned up old audit logs", "deleted_count", result.Count)
	return nil
}

// StartCleanupScheduler starts a goroutine that runs CleanupOldLogs daily at 00:00
func StartCleanupScheduler() {
	cfg := config.Get()

	// Check if audit logging is enabled
	if !cfg.Audit.Enabled {
		return
	}

	workerMutex.Lock()
	defer workerMutex.Unlock()

	// Don't start if already running
	if cleanupRunning {
		return
	}

	cleanupCtx, cleanupCancel = context.WithCancel(context.Background())
	cleanupRunning = true

	go func() {
		defer func() {
			workerMutex.Lock()
			cleanupRunning = false
			workerMutex.Unlock()
		}()

		logger := config.SetupLogger()
		logger.Info("Starting audit log cleanup scheduler")

		for {
			// Calculate time until next midnight (00:00)
			now := time.Now()
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			duration := nextMidnight.Sub(now)

			// Wait until midnight or context cancellation
			select {
			case <-time.After(duration):
				// Run cleanup
				if err := CleanupOldLogs(); err != nil {
					logger.Error("Scheduled cleanup failed", "error", err)
				} else {
					logger.Info("Scheduled cleanup completed successfully")
				}
			case <-cleanupCtx.Done():
				logger.Info("Cleanup scheduler stopped")
				return
			}
		}
	}()
}

// StopCleanupScheduler gracefully stops the cleanup scheduler
func StopCleanupScheduler() {
	workerMutex.Lock()
	defer workerMutex.Unlock()

	if cleanupCancel != nil {
		cleanupCancel()
	}
}

// TriggerCleanupNow manually triggers a cleanup operation (useful for testing or admin operations)
func TriggerCleanupNow() error {
	logger := config.SetupLogger()
	logger.Info("Manual cleanup triggered")

	if err := CleanupOldLogs(); err != nil {
		logger.Error("Manual cleanup failed", "error", err)
		return err
	}

	logger.Info("Manual cleanup completed successfully")
	return nil
}
