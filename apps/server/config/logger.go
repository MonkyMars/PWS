// Package config provides centralized logger configuration for the PWS application.
package config

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/MonkyMars/PWS/types"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

// Logger wraps the standard library's structured logger with additional functionality
// specific to HTTP request logging and application-specific log formatting.
type Logger struct {
	*slog.Logger
}

// SetupLogger creates and configures a new Logger instance based on the centralized configuration.
// It sets up structured logging with appropriate log levels and custom formatting for timestamps.
//
// Returns a configured Logger instance ready for use throughout the application.
func SetupLogger() *Logger {
	cfg := Get()

	var level slog.Level
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Format time as compact readable format
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(time.Now().Format("15:04:05")),
				}
			}
			// Add app name to all log entries
			return a
		},
		AddSource: true,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler).With("app", cfg.AppName)

	return &Logger{logger}
}

// HTTPMiddleware returns a Fiber middleware handler for HTTP request logging.
// This middleware logs each HTTP request with timing information, status codes,
// and client details in a structured format.
//
// Returns a Fiber middleware handler function that can be added to the application middleware stack.
func (l *Logger) HTTPMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get status code
		status := c.Response().StatusCode()

		// Get method and path
		method := c.Method()
		path := c.Path()

		// Get IP
		ip := c.IP()

		// Log with different levels based on status code
		logLevel := slog.LevelInfo
		if status >= 400 && status < 500 {
			logLevel = slog.LevelWarn
		} else if status >= 500 {
			logLevel = slog.LevelError
		}

		// Create log message
		message := fmt.Sprintf("%s %s - %d %v %s",
			method,
			path,
			status,
			duration.Round(time.Microsecond),
			ip,
		)

		// Log with structured attributes
		attrs := []slog.Attr{
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("duration", duration),
			slog.String("ip", ip),
		}

		l.LogAttrs(context.TODO(), logLevel, message, attrs...)

		return err
	}
}

// ServerStart logs a message indicating that the server is starting.
// This should be called before the server begins listening for connections.
func (l *Logger) ServerStart() {
	cfg := Get()
	l.Info("Server starting",
		slog.String("port", cfg.Port),
		slog.String("environment", cfg.Environment),
		slog.String("address", cfg.GetServerAddress()),
	)
}

// ServerReady logs a message indicating that the server is ready to accept connections.
// This should be called after the server has successfully started and is listening.
func (l *Logger) ServerReady() {
	cfg := Get()
	l.Info("Server ready",
		slog.String("url", fmt.Sprintf("http://localhost:%s", cfg.Port)),
		slog.String("environment", cfg.Environment),
	)
}

// ServerError logs server-level errors with appropriate error-level logging.
// This should be used for critical server errors that may affect application availability.
//
// Parameters:
//   - err: The error that occurred
func (l *Logger) ServerError(err error) {
	l.Error("Server error", slog.String("error", err.Error()))
}

// DatabaseConnected logs successful database connection
func (l *Logger) DatabaseConnected() {
	cfg := Get()
	l.Info("Database connected",
		slog.String("host", cfg.Database.Host),
		slog.Int("port", cfg.Database.Port),
		slog.String("database", cfg.Database.Name),
		slog.Int("max_conns", cfg.Database.MaxConns),
	)
}

// DatabaseError logs database-related errors
func (l *Logger) DatabaseError(operation string, err error) {
	l.Error("Database error",
		slog.String("operation", operation),
		slog.String("error", err.Error()),
	)
}

// RouteRegistered logs debug information about route registration.
// This is useful during development and debugging to understand which routes are available.
//
// Parameters:
//   - method: The HTTP method for the route (GET, POST, etc.)
//   - path: The URL path pattern for the route
func (l *Logger) RouteRegistered(method, path string) {
	l.Debug("Route registered",
		slog.String("method", method),
		slog.String("path", path),
	)
}

// ConfigLoaded logs that configuration has been successfully loaded
func (l *Logger) ConfigLoaded() {
	cfg := Get()
	l.Info("Configuration loaded",
		slog.String("app_name", cfg.AppName),
		slog.String("environment", cfg.Environment),
		slog.String("log_level", cfg.LogLevel),
	)
}

// Shutdown logs application shutdown
func (l *Logger) Shutdown(reason string) {
	l.Info("Application shutting down",
		slog.String("reason", reason),
	)
}

// Performance logs performance metrics
func (l *Logger) Performance(operation string, duration time.Duration) {
	l.Debug("Performance metric",
		slog.String("operation", operation),
		slog.Duration("duration", duration),
	)
}

// AuditError logs error messages to both the standard logger and the audit system.
// This function creates an audit log entry that gets batched and stored in the database
// via the audit worker for persistent error tracking and analysis.
//
// Parameters:
//   - message: A descriptive error message
//   - attrs: Additional structured attributes to include in both logs
func (l *Logger) AuditError(message string, attrs ...any) {
	// Log to standard logger first
	l.Error(message, attrs...)

	// Create audit log entry with validation
	auditAttrs := make(map[string]any)

	// Process attrs in pairs (key, value)
	for i := 0; i < len(attrs)-1; i += 2 {
		if key, ok := attrs[i].(string); ok && key != "" {
			auditAttrs[key] = attrs[i+1]
		}
	}

	// Capture source information
	source := ""
	if _, file, line, ok := runtime.Caller(1); ok {
		if idx := strings.LastIndex(file, "/"); idx >= 0 {
			file = file[idx+1:]
		}
		source = fmt.Sprintf("%s:%d", file, line)
	}

	auditLog := types.AuditLog{
		Timestamp: time.Now(),
		Level:     "ERROR",
		Message:   message,
		Attrs:     auditAttrs,
		Source:    source,
	}

	entryHash := generateEntryHash(auditLog)
	auditLog.EntryHash = entryHash

	// Send to audit worker (non-blocking)
	addAuditLogFunc := getAddAuditLogFunc()
	if addAuditLogFunc != nil {
		addAuditLogFunc(auditLog)
	}
}

func (l *Logger) AuditWarn(message string, attrs ...any) {
	// Log to standard logger first
	l.Warn(message, attrs...)

	// Create audit log entry with validation
	auditAttrs := make(map[string]any)

	// Process attrs in pairs (key, value)
	for i := 0; i < len(attrs)-1; i += 2 {
		if key, ok := attrs[i].(string); ok && key != "" {
			auditAttrs[key] = attrs[i+1]
		}
	}

	// Capture source information
	source := ""
	if _, file, line, ok := runtime.Caller(1); ok {
		if idx := strings.LastIndex(file, "/"); idx >= 0 {
			file = file[idx+1:]
		}
		source = fmt.Sprintf("%s:%d", file, line)
	}

	auditLog := types.AuditLog{
		Timestamp: time.Now(),
		Level:     "WARN",
		Message:   message,
		Attrs:     auditAttrs,
		Source:    source,
	}

	entryHash := generateEntryHash(auditLog)
	auditLog.EntryHash = entryHash

	// Send to audit worker (non-blocking)
	addAuditLogFunc := getAddAuditLogFunc()
	if addAuditLogFunc != nil {
		addAuditLogFunc(auditLog)
	}
}

// getAddAuditLogFunc returns the AddAuditLog function to avoid circular imports
// This uses a lazy loading approach to access the workers.AddAuditLog function
func getAddAuditLogFunc() func(types.AuditLog) {
	auditMutex.RLock()
	defer auditMutex.RUnlock()
	return globalAddAuditLogFunc
}

// Global variable to hold the AddAuditLog function reference
var (
	globalAddAuditLogFunc func(types.AuditLog)
	auditMutex            sync.RWMutex
)

// SetAuditLogFunc sets the audit log function to avoid circular dependencies.
// This should be called during application initialization to wire up the audit logging.
func SetAuditLogFunc(fn func(types.AuditLog)) {
	if fn == nil {
		return // Don't set nil function
	}
	auditMutex.Lock()
	defer auditMutex.Unlock()
	globalAddAuditLogFunc = fn
}

// generateEntryHash creates a unique hash for an audit log entry
// This is used for deduplication to prevent the same entry from being inserted multiple times
func generateEntryHash(entry types.AuditLog) string {
	// Validate required fields
	if entry.Message == "" || entry.Level == "" {
		return "" // Return empty hash for invalid entries
	}

	// Normalize timestamp to second precision to handle minor timing differences
	timestamp := entry.Timestamp.Truncate(time.Second).Unix()

	// Handle nil attrs map
	attrs := entry.Attrs
	if attrs == nil {
		attrs = make(map[string]any)
	}

	// Create a deterministic representation of the entry
	data := map[string]any{
		"timestamp": timestamp,
		"level":     entry.Level,
		"message":   entry.Message,
		"attrs":     attrs,
	}

	// Convert to JSON for consistent hashing
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Fallback to a simple string representation with sanitized message
		sanitizedMessage := entry.Message
		if len(sanitizedMessage) > 100 {
			sanitizedMessage = sanitizedMessage[:100] + "..."
		}
		fallbackHash := fmt.Sprintf("%d_%s_%s", timestamp, entry.Level, sanitizedMessage)

		// Ensure fallback hash is not too long for database column
		if len(fallbackHash) > 64 {
			// Use a hash of the fallback string if it's too long
			hash := sha256.Sum256([]byte(fallbackHash))
			return fmt.Sprintf("%x", hash)
		}
		return fallbackHash
	}

	// Generate SHA256 hash
	hash := sha256.Sum256(jsonData)
	return fmt.Sprintf("%x", hash)
}
