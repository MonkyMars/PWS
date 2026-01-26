package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
)

// Start starts the audit worker
func (aw *AuditWorker) Start() error {
	aw.mu.Lock()
	defer aw.mu.Unlock()

	if aw.running {
		return fmt.Errorf("audit worker already running")
	}

	if !aw.cfg.Audit.Enabled {
		return nil
	}

	aw.running = true
	aw.wg.Add(1)
	go aw.run()

	return nil
}

// Stop gracefully stops the audit worker
func (aw *AuditWorker) Stop(ctx context.Context) error {
	aw.mu.Lock()
	if !aw.running {
		aw.mu.Unlock()
		return nil
	}
	aw.cancel()
	aw.mu.Unlock()

	// Wait for worker to finish with timeout
	done := make(chan struct{})
	go func() {
		aw.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		aw.logger.Info("Audit worker stopped successfully")
		return nil
	case <-ctx.Done():
		aw.logger.Warn("Audit worker stop timed out")
		return ctx.Err()
	}
}

// AddLog adds an audit log entry to the processing queue
func (aw *AuditWorker) AddLog(entry types.AuditLog) {
	if !aw.cfg.Audit.Enabled {
		return
	}

	// Validate input
	if entry.Message == "" {
		return
	}

	// Set timestamp if not provided
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	aw.mu.RLock()
	running := aw.running
	failures := aw.stats.FailureCount
	aw.mu.RUnlock()

	if !running {
		aw.logger.Warn("Audit worker not running, dropping log entry",
			"level", entry.Level,
			"message", entry.Message)
		return
	}

	// Circuit breaker: if too many failures, drop new entries temporarily
	if failures >= aw.cfg.Audit.MaxFailures {
		aw.mu.Lock()
		aw.stats.TotalDropped++
		aw.mu.Unlock()

		aw.logger.Warn("Audit worker in failure mode, dropping log entry",
			"level", entry.Level,
			"message", entry.Message,
			"failure_count", failures)
		return
	}

	select {
	case aw.auditChan <- entry:
		// Successfully added to queue
	default:
		// Channel is full, update dropped count and log warning
		aw.mu.Lock()
		aw.stats.TotalDropped++
		aw.mu.Unlock()

		aw.logger.Warn("Audit log channel is full, dropping log entry",
			"level", entry.Level,
			"message", entry.Message,
			"queue_size", len(aw.auditChan))
	}
}

// HealthStatus returns the current health status of the audit worker
func (aw *AuditWorker) HealthStatus() map[string]any {
	if aw == nil {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "audit worker is nil",
		}
	}

	if aw.cfg == nil {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "audit worker configuration is nil",
		}
	}

	aw.mu.RLock()
	defer aw.mu.RUnlock()

	queueSize := 0
	if aw.auditChan != nil {
		queueSize = len(aw.auditChan)
	}

	isHealthy := aw.cfg.Audit.Enabled && aw.running && aw.stats.FailureCount < aw.cfg.Audit.MaxFailures

	return map[string]any{
		"enabled":         aw.cfg.Audit.Enabled,
		"worker_running":  aw.running,
		"queue_size":      queueSize,
		"queue_capacity":  aw.cfg.Audit.ChannelSize,
		"last_flush_time": aw.stats.LastFlushTime,
		"failure_count":   aw.stats.FailureCount,
		"total_processed": aw.stats.TotalProcessed,
		"total_dropped":   aw.stats.TotalDropped,
		"is_healthy":      isHealthy,
		"configuration": map[string]any{
			"batch_size":     aw.cfg.Audit.BatchSize,
			"flush_time":     aw.cfg.Audit.FlushTime.String(),
			"max_retries":    aw.cfg.Audit.MaxRetries,
			"max_failures":   aw.cfg.Audit.MaxFailures,
			"retention_days": aw.cfg.Audit.RetentionDays,
		},
	}
}

// run is the main worker loop
func (aw *AuditWorker) run() {
	defer aw.wg.Done()
	defer func() {
		aw.mu.Lock()
		aw.running = false
		close(aw.auditChan)
		aw.mu.Unlock()
	}()

	batch := make([]types.AuditLog, 0, aw.cfg.Audit.BatchSize)
	ticker := time.NewTicker(aw.cfg.Audit.FlushTime)
	defer ticker.Stop()

	for {
		select {
		case entry := <-aw.auditChan:
			batch = append(batch, entry)

			// Flush immediately if batch full
			if len(batch) >= aw.cfg.Audit.BatchSize {
				aw.flushBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				aw.flushBatch(batch)
				batch = batch[:0]
			}

		case <-aw.ctx.Done():
			// Flush remaining entries before shutting down
			if len(batch) > 0 {
				aw.flushBatch(batch)
			}
			// Drain any remaining entries in channel
			for {
				select {
				case entry := <-aw.auditChan:
					batch = append(batch, entry)
					if len(batch) >= aw.cfg.Audit.BatchSize {
						aw.flushBatch(batch)
						batch = batch[:0]
					}
				default:
					if len(batch) > 0 {
						aw.flushBatch(batch)
					}
					return
				}
			}
		}
	}
}

// flushBatch writes a batch of audit logs to the database
func (aw *AuditWorker) flushBatch(entries []types.AuditLog) {
	if len(entries) == 0 {
		return
	}

	var err error
	var successfulInserts int64

	for attempt := 0; attempt < aw.cfg.Audit.MaxRetries; attempt++ {
		successfulInserts, err = aw.tryFlushBatchWithCount(entries)
		if err == nil {
			aw.mu.Lock()
			aw.stats.FailureCount = 0 // Reset failure count on success
			aw.stats.LastFlushTime = time.Now()
			aw.stats.TotalProcessed += successfulInserts
			aw.mu.Unlock()

			aw.logger.Debug("Flushed audit log batch",
				"count", len(entries),
				"successful_inserts", successfulInserts,
				"attempt", attempt+1)
			return
		}

		// Log retry attempt
		aw.logger.Warn("Audit batch flush failed, retrying",
			"attempt", attempt+1,
			"max_retries", aw.cfg.Audit.MaxRetries,
			"error", err,
			"batch_size", len(entries))

		if attempt < aw.cfg.Audit.MaxRetries-1 {
			// Exponential backoff: 100ms, 200ms, 400ms
			backoff := time.Duration(100*(1<<attempt)) * time.Millisecond
			time.Sleep(backoff)
		}
	}

	// After all retries failed, update failure count
	aw.mu.Lock()
	aw.stats.FailureCount++
	aw.mu.Unlock()

	// After all retries failed, log the error but don't crash
	aw.logger.Error("Failed to flush audit log batch after retries",
		"error", err,
		"batch_size", len(entries),
		"max_retries", aw.cfg.Audit.MaxRetries,
		"total_failures", aw.stats.FailureCount)
}

// tryFlushBatchWithCount attempts to flush a batch and returns the count of successful inserts
func (aw *AuditWorker) tryFlushBatchWithCount(entries []types.AuditLog) (int64, error) {
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
			"source":     entry.Source,
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
		aw.logger.Debug("Skipped invalid audit entries during flush",
			"skipped_count", skippedEntries,
			"total_entries", len(entries))
	}

	return successfulInserts, nil
}
