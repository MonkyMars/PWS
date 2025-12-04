package workers

import (
	"context"
	"fmt"
	"maps"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// RouteService tracks metrics for a specific route service
type RouteService struct {
	Name         string
	BasePath     string
	RequestCount int64
	ErrorCount   int64
	TotalLatency time.Duration
	LastStatus   int
	StartTime    time.Time
	mutex        sync.RWMutex
}

// Start starts the health worker
func (hw *HealthWorker) Start() error {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	if hw.running {
		return fmt.Errorf("health worker already running")
	}

	if !hw.cfg.Health.Enabled {
		return nil
	}

	hw.running = true

	// Start the health reporter
	hw.wg.Add(1)
	go hw.healthReporter()

	// Start the log processor
	hw.wg.Add(1)
	go hw.logProcessor()

	return nil
}

// Stop gracefully stops the health worker
func (hw *HealthWorker) Stop(ctx context.Context) error {
	hw.mu.Lock()
	if !hw.running {
		hw.mu.Unlock()
		return nil
	}
	hw.cancel()
	hw.mu.Unlock()

	// Wait for worker to finish with timeout
	done := make(chan struct{})
	go func() {
		hw.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		hw.logger.Info("Health worker stopped successfully")
		return nil
	case <-ctx.Done():
		hw.logger.Warn("Health worker stop timed out")
		return ctx.Err()
	}
}

// DiscoverRoutes automatically discovers all base routes from the fiber app
func (hw *HealthWorker) DiscoverRoutes(app *fiber.App) {
	if !hw.cfg.Health.Enabled {
		return
	}

	// Get all routes from the fiber app
	routes := app.GetRoutes()
	hw.logger.Info("Starting route discovery", "total_routes", len(routes))

	discoveredServices := make(map[string]bool)

	for _, route := range routes {
		basePath := hw.extractBasePath(route.Path)

		if basePath != "" && !discoveredServices[basePath] {
			hw.RegisterService(basePath)
			discoveredServices[basePath] = true
			hw.logger.Info("Registered new service", "service", basePath)
		}
	}

	services := make([]string, 0, len(discoveredServices))
	for service := range discoveredServices {
		services = append(services, service)
	}
	hw.logger.Info("Auto-discovered route services", "services", services, "count", len(services))
}

// RegisterService registers a service for health monitoring
func (hw *HealthWorker) RegisterService(serviceName string) {
	hw.mu.Lock()
	defer hw.mu.Unlock()

	if _, exists := hw.services[serviceName]; !exists {
		hw.services[serviceName] = &RouteService{
			Name:      serviceName,
			BasePath:  "/" + serviceName,
			StartTime: time.Now(),
		}
	}
}

// RecordRequest records a request for a specific service
func (hw *HealthWorker) RecordRequest(serviceName string, statusCode int, latency time.Duration) {
	if !hw.cfg.Health.Enabled {
		return
	}

	hw.mu.RLock()
	service, exists := hw.services[serviceName]
	hw.mu.RUnlock()

	if !exists {
		return
	}

	service.mutex.Lock()
	defer service.mutex.Unlock()

	service.RequestCount++
	service.TotalLatency += latency
	service.LastStatus = statusCode

	if statusCode >= 400 {
		service.ErrorCount++
	}
}

// HealthStatus returns the current health status of the health worker
func (hw *HealthWorker) HealthStatus() map[string]any {
	if hw == nil {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "health worker is nil",
		}
	}

	if hw.cfg == nil {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "health worker configuration is nil",
		}
	}

	hw.mu.RLock()
	defer hw.mu.RUnlock()

	queueSize := 0
	if hw.healthChan != nil {
		queueSize = len(hw.healthChan)
	}

	serviceCount := len(hw.services)
	isHealthy := hw.cfg.Health.Enabled && hw.running

	return map[string]any{
		"enabled":         hw.cfg.Health.Enabled,
		"worker_running":  hw.running,
		"queue_size":      queueSize,
		"queue_capacity":  hw.cfg.Health.ChannelSize,
		"last_flush_time": hw.lastFlushTime,
		"service_count":   serviceCount,
		"is_healthy":      isHealthy,
		"configuration": map[string]any{
			"report_interval": hw.cfg.Health.ReportInterval.String(),
			"flush_time":      hw.cfg.Health.FlushTime.String(),
			"channel_size":    hw.cfg.Health.ChannelSize,
		},
	}
}

// GetServiceStats returns current statistics for a service
func (hw *HealthWorker) GetServiceStats(serviceName string) *RouteService {
	hw.mu.RLock()
	defer hw.mu.RUnlock()

	service, exists := hw.services[serviceName]
	if !exists {
		return nil
	}

	// Return a copy to avoid data races
	service.mutex.RLock()
	defer service.mutex.RUnlock()

	return &RouteService{
		Name:         service.Name,
		BasePath:     service.BasePath,
		RequestCount: service.RequestCount,
		ErrorCount:   service.ErrorCount,
		TotalLatency: service.TotalLatency,
		LastStatus:   service.LastStatus,
		StartTime:    service.StartTime,
	}
}

// GetAllServices returns a list of all registered services
func (hw *HealthWorker) GetAllServices() []string {
	hw.mu.RLock()
	defer hw.mu.RUnlock()

	services := make([]string, 0, len(hw.services))
	for name := range hw.services {
		services = append(services, name)
	}
	return services
}

// healthReporter generates health reports every configured interval
func (hw *HealthWorker) healthReporter() {
	defer hw.wg.Done()

	ticker := time.NewTicker(hw.cfg.Health.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-hw.ctx.Done():
			return
		case <-ticker.C:
			hw.generateHealthReports()
		}
	}
}

// generateHealthReports creates health reports for all services
func (hw *HealthWorker) generateHealthReports() {
	hw.mu.RLock()
	services := make(map[string]*RouteService)
	maps.Copy(services, hw.services)
	hw.mu.RUnlock()

	hw.logger.Info("Generating health reports", "services_count", len(services))

	for serviceName, service := range services {
		healthLog := hw.createHealthLog(serviceName, service)

		select {
		case hw.healthChan <- healthLog:
			// Successfully sent
		default:
			hw.logger.Warn("Health report channel full, dropping report", "service", serviceName)
		}
	}
}

// createHealthLog creates a health log from service metrics
func (hw *HealthWorker) createHealthLog(serviceName string, service *RouteService) types.HealthLog {
	service.mutex.RLock()
	defer service.mutex.RUnlock()

	var averageLatency time.Duration
	if service.RequestCount > 0 {
		averageLatency = time.Duration(service.TotalLatency.Milliseconds() / service.RequestCount)
	}

	statusCode := service.LastStatus
	if statusCode == 0 {
		statusCode = 200 // Default to OK if no requests recorded
	}

	// Calculate time span since last flush
	var timeSpan time.Duration
	if hw.lastFlushTime.IsZero() {
		timeSpan = time.Since(service.StartTime)
	} else {
		timeSpan = time.Since(hw.lastFlushTime)
	}

	// Capture source information (file:line)
	source := ""
	if _, file, line, ok := runtime.Caller(1); ok {
		// Extract just the filename, not the full path
		if idx := strings.LastIndex(file, "/"); idx >= 0 {
			file = file[idx+1:]
		}
		source = fmt.Sprintf("%s:%d", file, line)
	}

	return types.HealthLog{
		Timestamp:      time.Now(),
		Service:        serviceName,
		StatusCode:     statusCode,
		RequestCount:   service.RequestCount,
		ErrorCount:     service.ErrorCount,
		AverageLatency: averageLatency,
		TimeSpan:       timeSpan,
		Source:         source,
	}
}

// logProcessor processes and flushes health logs
func (hw *HealthWorker) logProcessor() {
	defer hw.wg.Done()
	defer close(hw.healthChan)

	flushTicker := time.NewTicker(hw.cfg.Health.FlushTime)
	defer flushTicker.Stop()

	var logBatch []types.HealthLog

	for {
		select {
		case <-hw.ctx.Done():
			// Flush remaining logs before exiting
			if len(logBatch) > 0 {
				hw.flushLogs(logBatch)
			}
			return

		case healthLog := <-hw.healthChan:
			logBatch = append(logBatch, healthLog)

		case <-flushTicker.C:
			if len(logBatch) > 0 {
				hw.logger.Info("Flushing health logs", "count", len(logBatch))
				hw.flushLogs(logBatch)
				logBatch = logBatch[:0] // Clear the batch
			}
		}
	}
}

// flushLogs writes the batch of health logs to the database
func (hw *HealthWorker) flushLogs(logs []types.HealthLog) {
	if len(logs) == 0 {
		return
	}

	items := make([]any, len(logs))
	for i, log := range logs {
		items[i] = hw.convertHealthLogToMap(log)
	}

	query := services.Query().
		SetOperation("insert").
		SetTable(lib.TableHealthLogs).
		SetEntries(items)

	_, err := database.ExecuteQuery[any](query)
	if err != nil {
		hw.logger.Error("Failed to flush health logs", "error", err, "count", len(logs))
		return
	}

	// Update last flush time after successful flush
	hw.mu.Lock()
	hw.lastFlushTime = time.Now()
	hw.mu.Unlock()
}

// convertHealthLogToMap converts a HealthLog struct to map[string]any for database insertion
func (hw *HealthWorker) convertHealthLogToMap(log types.HealthLog) map[string]any {
	latencyMs := log.AverageLatency.Milliseconds()
	timeSpanSeconds := int(log.TimeSpan.Seconds())

	return map[string]any{
		"timestamp":       log.Timestamp,
		"service":         log.Service,
		"status_code":     log.StatusCode,
		"request_count":   log.RequestCount,
		"error_count":     log.ErrorCount,
		"average_latency": latencyMs,
		"time_span":       timeSpanSeconds,
		"source":          log.Source,
	}
}

// extractBasePath extracts the base path from a route path
func (hw *HealthWorker) extractBasePath(path string) string {
	// Remove leading slash and split by slash
	trimmed := strings.TrimPrefix(path, "/")
	if trimmed == "" {
		return ""
	}

	segments := strings.Split(trimmed, "/")
	if len(segments) == 0 {
		return ""
	}

	basePath := segments[0]

	// Skip health and system routes
	if basePath == "health" || basePath == "metrics" || basePath == "logs" {
		return ""
	}

	// Handle parameterized routes (remove :param)
	if strings.HasPrefix(basePath, ":") {
		return ""
	}

	return basePath
}

// Backward compatibility functions

// LogHealthEvent logs a single health event
func LogHealthEvent(entry types.HealthLog) {
	manager := GetGlobalManager()
	if manager.healthWorker == nil {
		return
	}

	select {
	case manager.healthWorker.healthChan <- entry:
	default:
		// Channel is full, drop the log entry
	}
}

// GetServiceStats returns current statistics for a service (backward compatibility)
func GetServiceStats(serviceName string) (*RouteService, error) {
	manager := GetGlobalManager()
	if manager.healthWorker != nil {
		return manager.healthWorker.GetServiceStats(serviceName), nil
	}
	return nil, lib.ErrServiceUnavailable
}

// GetAllServices returns a list of all registered services (backward compatibility)
func GetAllServices() []string {
	manager := GetGlobalManager()
	if manager.healthWorker != nil {
		return manager.healthWorker.GetAllServices()
	}
	return nil
}
