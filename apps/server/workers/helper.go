package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/types"
)

// DeadLetterQueue handles failed audit logs that couldn't be written to the database
type DeadLetterQueue struct {
	filePath    string
	maxFileSize int64
	maxFiles    int
	mu          sync.Mutex
	logger      *config.Logger
	cfg         *config.Config
}

// DeadLetterEntry represents an entry in the dead letter queue
type DeadLetterEntry struct {
	Timestamp     time.Time      `json:"timestamp"`
	FailureTime   time.Time      `json:"failure_time"`
	RetryCount    int            `json:"retry_count"`
	LastError     string         `json:"last_error"`
	OriginalLog   types.AuditLog `json:"original_log"`
	FailureReason string         `json:"failure_reason"`
}

// NewDeadLetterQueue creates a new dead letter queue instance
func NewDeadLetterQueue(cfg *config.Config, logger *config.Logger) *DeadLetterQueue {
	// Default file path
	filePath := getEnvString("DLQ_FILE_PATH", "/tmp/pws_dead_letter_queue.jsonl")

	// Default max file size (10MB)
	maxFileSize := int64(getEnvInt("DLQ_MAX_FILE_SIZE", 10*1024*1024))

	// Default max number of files to keep
	maxFiles := getEnvInt("DLQ_MAX_FILES", 5)

	return &DeadLetterQueue{
		filePath:    filePath,
		maxFileSize: maxFileSize,
		maxFiles:    maxFiles,
		logger:      logger,
		cfg:         cfg,
	}
}

// AddFailedLog adds a failed audit log to the dead letter queue
func (dlq *DeadLetterQueue) AddFailedLog(log types.AuditLog, failureReason string, lastError error) error {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()

	entry := DeadLetterEntry{
		Timestamp:     time.Now(),
		FailureTime:   time.Now(),
		RetryCount:    0,
		LastError:     "",
		OriginalLog:   log,
		FailureReason: failureReason,
	}

	if lastError != nil {
		entry.LastError = lastError.Error()
	}

	return dlq.writeEntry(entry)
}

// AddFailedBatch adds a batch of failed audit logs to the dead letter queue
func (dlq *DeadLetterQueue) AddFailedBatch(logs []types.AuditLog, failureReason string, lastError error) error {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()

	errorStr := ""
	if lastError != nil {
		errorStr = lastError.Error()
	}

	for _, log := range logs {
		entry := DeadLetterEntry{
			Timestamp:     time.Now(),
			FailureTime:   time.Now(),
			RetryCount:    0,
			LastError:     errorStr,
			OriginalLog:   log,
			FailureReason: failureReason,
		}

		if err := dlq.writeEntry(entry); err != nil {
			// Log the error but continue with other entries
			dlq.logger.Error("Failed to write dead letter queue entry",
				"error", err,
				"original_log_message", log.Message)
		}
	}

	return nil
}

// writeEntry writes a single entry to the dead letter queue file
func (dlq *DeadLetterQueue) writeEntry(entry DeadLetterEntry) error {
	// Ensure directory exists
	dir := filepath.Dir(dlq.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Check if we need to rotate the file
	if err := dlq.rotateIfNeeded(); err != nil {
		dlq.logger.Warn("Failed to rotate dead letter queue file", "error", err)
	}

	// Open file for appending
	file, err := os.OpenFile(dlq.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open dead letter queue file: %w", err)
	}
	defer file.Close()

	// Marshal entry to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal dead letter entry: %w", err)
	}

	// Write JSON line
	if _, err := file.Write(append(jsonData, '\n')); err != nil {
		return fmt.Errorf("failed to write to dead letter queue file: %w", err)
	}

	return nil
}

// rotateIfNeeded rotates the dead letter queue file if it exceeds the maximum size
func (dlq *DeadLetterQueue) rotateIfNeeded() error {
	// Check current file size
	info, err := os.Stat(dlq.filePath)
	if err != nil {
		// File doesn't exist yet, no need to rotate
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to stat dead letter queue file: %w", err)
	}

	// Check if rotation is needed
	if info.Size() < dlq.maxFileSize {
		return nil
	}

	// Rotate files
	return dlq.rotateFiles()
}

// rotateFiles rotates the dead letter queue files
func (dlq *DeadLetterQueue) rotateFiles() error {
	base := dlq.filePath
	ext := filepath.Ext(base)
	nameWithoutExt := base[:len(base)-len(ext)]

	// Remove the oldest file if we're at the limit
	oldestFile := fmt.Sprintf("%s.%d%s", nameWithoutExt, dlq.maxFiles-1, ext)
	if _, err := os.Stat(oldestFile); err == nil {
		if err := os.Remove(oldestFile); err != nil {
			dlq.logger.Warn("Failed to remove oldest dead letter queue file",
				"file", oldestFile, "error", err)
		}
	}

	// Rotate existing files
	for i := dlq.maxFiles - 2; i >= 1; i-- {
		oldName := fmt.Sprintf("%s.%d%s", nameWithoutExt, i, ext)
		newName := fmt.Sprintf("%s.%d%s", nameWithoutExt, i+1, ext)

		if _, err := os.Stat(oldName); err == nil {
			if err := os.Rename(oldName, newName); err != nil {
				dlq.logger.Warn("Failed to rotate dead letter queue file",
					"from", oldName, "to", newName, "error", err)
			}
		}
	}

	// Rotate the current file
	rotatedName := fmt.Sprintf("%s.1%s", nameWithoutExt, ext)
	if err := os.Rename(base, rotatedName); err != nil {
		return fmt.Errorf("failed to rotate current dead letter queue file: %w", err)
	}

	dlq.logger.Info("Rotated dead letter queue file", "new_file", rotatedName)
	return nil
}

// ReadEntries reads all entries from the dead letter queue files
func (dlq *DeadLetterQueue) ReadEntries() ([]DeadLetterEntry, error) {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()

	var allEntries []DeadLetterEntry

	// Read from main file and rotated files
	files := dlq.getAllFiles()

	for _, filePath := range files {
		entries, err := dlq.readEntriesFromFile(filePath)
		if err != nil {
			dlq.logger.Warn("Failed to read dead letter queue file",
				"file", filePath, "error", err)
			continue
		}
		allEntries = append(allEntries, entries...)
	}

	return allEntries, nil
}

// getAllFiles returns all dead letter queue files (main + rotated)
func (dlq *DeadLetterQueue) getAllFiles() []string {
	files := []string{}

	// Add main file if it exists
	if _, err := os.Stat(dlq.filePath); err == nil {
		files = append(files, dlq.filePath)
	}

	// Add rotated files
	base := dlq.filePath
	ext := filepath.Ext(base)
	nameWithoutExt := base[:len(base)-len(ext)]

	for i := 1; i < dlq.maxFiles; i++ {
		rotatedFile := fmt.Sprintf("%s.%d%s", nameWithoutExt, i, ext)
		if _, err := os.Stat(rotatedFile); err == nil {
			files = append(files, rotatedFile)
		}
	}

	return files
}

// readEntriesFromFile reads entries from a specific file
func (dlq *DeadLetterQueue) readEntriesFromFile(filePath string) ([]DeadLetterEntry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var entries []DeadLetterEntry
	lines := strings.Split(string(data), "\n")

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var entry DeadLetterEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			dlq.logger.Warn("Failed to unmarshal dead letter queue entry",
				"file", filePath, "line", lineNum+1, "error", err)
			continue
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// RetryFailedLogs attempts to retry failed logs from the dead letter queue
func (dlq *DeadLetterQueue) RetryFailedLogs(ctx context.Context, maxRetries int) (int, error) {
	entries, err := dlq.ReadEntries()
	if err != nil {
		return 0, fmt.Errorf("failed to read dead letter queue entries: %w", err)
	}

	if len(entries) == 0 {
		return 0, nil
	}

	dlq.logger.Info("Starting retry of failed audit logs", "count", len(entries))

	var successCount int
	var retryEntries []DeadLetterEntry

	for _, entry := range entries {
		// Skip entries that have exceeded max retries
		if entry.RetryCount >= maxRetries {
			dlq.logger.Debug("Skipping entry that exceeded max retries",
				"message", entry.OriginalLog.Message,
				"retry_count", entry.RetryCount)
			retryEntries = append(retryEntries, entry)
			continue
		}

		// Attempt to process the log again
		if err := dlq.retrySingleEntry(ctx, entry); err != nil {
			// Update retry count and add back to queue
			entry.RetryCount++
			entry.LastError = err.Error()
			entry.FailureTime = time.Now()
			retryEntries = append(retryEntries, entry)

			dlq.logger.Debug("Failed to retry audit log entry",
				"message", entry.OriginalLog.Message,
				"retry_count", entry.RetryCount,
				"error", err)
		} else {
			successCount++
			dlq.logger.Debug("Successfully retried audit log entry",
				"message", entry.OriginalLog.Message)
		}
	}

	// Rewrite the dead letter queue with failed retries
	if err := dlq.rewriteQueue(retryEntries); err != nil {
		dlq.logger.Error("Failed to rewrite dead letter queue", "error", err)
		// Don't return error here as the retry operation partially succeeded
	}

	dlq.logger.Info("Completed retry of failed audit logs",
		"success_count", successCount,
		"remaining_count", len(retryEntries))

	return successCount, nil
}

// retrySingleEntry attempts to process a single dead letter queue entry
func (dlq *DeadLetterQueue) retrySingleEntry(ctx context.Context, entry DeadLetterEntry) error {

	return fmt.Errorf("retry mechanism not fully implemented - placeholder")
}

// rewriteQueue rewrites the entire dead letter queue with the given entries
func (dlq *DeadLetterQueue) rewriteQueue(entries []DeadLetterEntry) error {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()

	// Remove all existing files
	files := dlq.getAllFiles()
	for _, file := range files {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			dlq.logger.Warn("Failed to remove dead letter queue file",
				"file", file, "error", err)
		}
	}

	// Write new entries
	for _, entry := range entries {
		if err := dlq.writeEntry(entry); err != nil {
			return fmt.Errorf("failed to write entry during queue rewrite: %w", err)
		}
	}

	return nil
}

// GetStats returns statistics about the dead letter queue
func (dlq *DeadLetterQueue) GetStats() map[string]any {
	entries, err := dlq.ReadEntries()
	if err != nil {
		return map[string]any{
			"error": err.Error(),
		}
	}

	totalSize := int64(0)
	files := dlq.getAllFiles()
	for _, file := range files {
		if info, err := os.Stat(file); err == nil {
			totalSize += info.Size()
		}
	}

	// Calculate retry counts
	retryCounts := make(map[int]int)
	for _, entry := range entries {
		retryCounts[entry.RetryCount]++
	}

	return map[string]any{
		"total_entries":    len(entries),
		"total_files":      len(files),
		"total_size_bytes": totalSize,
		"file_path":        dlq.filePath,
		"max_file_size":    dlq.maxFileSize,
		"max_files":        dlq.maxFiles,
		"retry_counts":     retryCounts,
		"files":            files,
	}
}

// Clear removes all entries from the dead letter queue
func (dlq *DeadLetterQueue) Clear() error {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()

	files := dlq.getAllFiles()
	var errors []error

	for _, file := range files {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			errors = append(errors, fmt.Errorf("failed to remove %s: %w", file, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to clear dead letter queue: %v", errors)
	}

	dlq.logger.Info("Cleared dead letter queue", "removed_files", len(files))
	return nil
}

// Helper functions for environment variables
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
