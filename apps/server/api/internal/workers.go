package internal

import (
	"time"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/workers"
	"github.com/gofiber/fiber/v3"
)

type WorkerRoutes struct{}

func NewWorkerRoutes() *WorkerRoutes {
	return &WorkerRoutes{}
}

func (wr *WorkerRoutes) GetWorkerHealth(c fiber.Ctx) error {
	manager := workers.GetGlobalManager()
	if manager == nil {
		return response.ServiceUnavailable(c, "Worker manager not available")
	}

	healthStatus := manager.HealthStatus()
	if healthStatus == nil {
		return response.ServiceUnavailable(c, "Unable to retrieve worker health status")
	}

	// Determine overall status
	isHealthy := false
	if healthVal, exists := healthStatus["is_healthy"]; exists && healthVal != nil {
		if healthy, ok := healthVal.(bool); ok {
			isHealthy = healthy
		}
	}
	if !isHealthy {
		return response.ServiceUnavailable(c, "Workers are not healthy")
	}

	return response.SuccessWithMessage(c, "Worker health status retrieved", healthStatus)
}

// getAuditWorkerHealth returns the health status of the audit worker
func (wr *WorkerRoutes) GetAuditWorkerHealth(c fiber.Ctx) error {
	healthStatus := workers.AuditHealthStatus()
	if healthStatus == nil {
		return response.ServiceUnavailable(c, "Unable to retrieve audit worker status")
	}

	isHealthy := false
	if healthVal, exists := healthStatus["is_healthy"]; exists && healthVal != nil {
		if healthy, ok := healthVal.(bool); ok {
			isHealthy = healthy
		}
	}

	if !isHealthy {
		return response.ServiceUnavailable(c, "Audit worker is not healthy")
	}

	return response.SuccessWithMessage(c, "Audit worker health status retrieved", healthStatus)
}

// getHealthWorkerHealth returns the health status of the health monitoring worker
func (wr *WorkerRoutes) GetHealthWorkerHealth(c fiber.Ctx) error {
	healthStatus := workers.ServiceHealthStatus()
	if healthStatus == nil {
		return response.ServiceUnavailable(c, "Unable to retrieve health worker status")
	}

	isHealthy := false
	if healthVal, exists := healthStatus["is_healthy"]; exists && healthVal != nil {
		if healthy, ok := healthVal.(bool); ok {
			isHealthy = healthy
		}
	}

	if !isHealthy {
		return response.ServiceUnavailable(c, "Health worker is not healthy")
	}

	return response.SuccessWithMessage(c, "Health worker status retrieved", healthStatus)
}

// getCleanupWorkerHealth returns the health status of the cleanup worker
func (wr *WorkerRoutes) GetCleanupWorkerHealth(c fiber.Ctx) error {
	manager := workers.GetGlobalManager()
	if manager == nil {
		return response.ServiceUnavailable(c, "Worker manager not available")
	}

	healthStatus := manager.HealthStatus()
	if healthStatus == nil {
		return response.ServiceUnavailable(c, "Unable to retrieve worker health status")
	}

	// Extract cleanup worker status
	var cleanupStatus map[string]any
	if cleanupVal, exists := healthStatus["cleanup"]; exists && cleanupVal != nil {
		if cleanup, ok := cleanupVal.(map[string]any); ok {
			cleanupStatus = cleanup
		}
	}

	if cleanupStatus == nil {
		cleanupStatus = map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "cleanup worker not available",
		}
	}

	isHealthy := false
	if healthVal, exists := cleanupStatus["is_healthy"]; exists && healthVal != nil {
		if healthy, ok := healthVal.(bool); ok {
			isHealthy = healthy
		}
	}

	if !isHealthy {
		return response.ServiceUnavailable(c, "Cleanup worker is not healthy")
	}

	return response.SuccessWithMessage(c, "Cleanup worker health status retrieved", cleanupStatus)
}

// getWorkerMetrics returns comprehensive metrics for all workers
func (wr *WorkerRoutes) GetWorkerMetrics(c fiber.Ctx) error {
	manager := workers.GetGlobalManager()
	if manager == nil {
		return response.ServiceUnavailable(c, "Worker manager not available")
	}

	healthStatus := manager.HealthStatus()
	if healthStatus == nil {
		return response.ServiceUnavailable(c, "Unable to retrieve worker metrics")
	}

	// Extract metrics from health status with null checks
	timestamp := time.Now()
	if ts, exists := healthStatus["timestamp"]; exists && ts != nil {
		timestamp = ts.(time.Time)
	}

	isHealthy := false
	if healthVal, exists := healthStatus["is_healthy"]; exists && healthVal != nil {
		if healthy, ok := healthVal.(bool); ok {
			isHealthy = healthy
		}
	}

	metrics := map[string]any{
		"timestamp": timestamp,
		"workers": map[string]any{
			"audit":   wr.ExtractWorkerMetrics(healthStatus["audit"]),
			"health":  wr.ExtractWorkerMetrics(healthStatus["health"]),
			"cleanup": wr.ExtractWorkerMetrics(healthStatus["cleanup"]),
		},
		"overall_healthy": isHealthy,
	}

	return response.SuccessWithMessage(c, "Worker metrics retrieved", metrics)
}

// getAuditWorkerMetrics returns detailed metrics for the audit worker
func (wr *WorkerRoutes) GetAuditWorkerMetrics(c fiber.Ctx) error {
	healthStatus := workers.AuditHealthStatus()

	metrics := map[string]any{
		"processed_total": healthStatus["total_processed"],
		"dropped_total":   healthStatus["total_dropped"],
		"failure_count":   healthStatus["failure_count"],
		"queue_size":      healthStatus["queue_size"],
		"queue_capacity":  healthStatus["queue_capacity"],
		"last_flush_time": healthStatus["last_flush_time"],
		"configuration":   healthStatus["configuration"],
	}

	return response.SuccessWithMessage(c, "Audit worker metrics retrieved", metrics)
}

// getHealthWorkerMetrics returns detailed metrics for the health monitoring worker
func (wr *WorkerRoutes) GetHealthWorkerMetrics(c fiber.Ctx) error {
	healthStatus := workers.ServiceHealthStatus()

	metrics := map[string]any{
		"service_count":   healthStatus["service_count"],
		"queue_size":      healthStatus["queue_size"],
		"queue_capacity":  healthStatus["queue_capacity"],
		"last_flush_time": healthStatus["last_flush_time"],
		"configuration":   healthStatus["configuration"],
	}

	return response.SuccessWithMessage(c, "Health worker metrics retrieved", metrics)
}

// getMonitoredServices returns a list of all services being monitored
func (wr *WorkerRoutes) GetMonitoredServices(c fiber.Ctx) error {
	services := workers.GetAllServices()

	if services == nil {
		services = []string{}
	}

	data := map[string]any{
		"services": services,
		"count":    len(services),
	}

	return response.SuccessWithMessage(c, "Monitored services retrieved", data)
}

// getServiceStatistics returns statistics for a specific monitored service
func (wr *WorkerRoutes) GetServiceStatistics(c fiber.Ctx) error {
	serviceName := c.Params("service")
	if serviceName == "" {
		return response.BadRequest(c, "Service name is required")
	}

	stats := workers.GetServiceStats(serviceName)
	if stats == nil {
		return response.NotFound(c, "Service not found")
	}

	// Convert stats to a map for JSON response
	statsData := map[string]any{
		"name":          stats.Name,
		"base_path":     stats.BasePath,
		"request_count": stats.RequestCount,
		"error_count":   stats.ErrorCount,
		"last_status":   stats.LastStatus,
		"start_time":    stats.StartTime,
	}

	// Calculate average latency if there are requests
	if stats.RequestCount > 0 {
		avgLatencyMs := stats.TotalLatency.Milliseconds() / stats.RequestCount
		statsData["average_latency_ms"] = avgLatencyMs
		statsData["error_rate"] = float64(stats.ErrorCount) / float64(stats.RequestCount)
	} else {
		statsData["average_latency_ms"] = 0
		statsData["error_rate"] = 0.0
	}

	return response.SuccessWithMessage(c, "Service statistics retrieved", statsData)
}

// triggerCleanup manually triggers cleanup operations
func (wr *WorkerRoutes) TriggerCleanup(c fiber.Ctx) error {
	err := workers.TriggerCleanupNow()
	if err != nil {
		return response.InternalServerError(c, "Failed to trigger cleanup")
	}

	return response.Accepted(c, "Cleanup triggered successfully")
}

// extractWorkerMetrics extracts relevant metrics from worker health status
func (wr *WorkerRoutes) ExtractWorkerMetrics(workerHealth any) map[string]any {
	if workerHealth == nil {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
		}
	}

	health, ok := workerHealth.(map[string]any)
	if !ok {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "invalid health data format",
		}
	}

	metrics := map[string]any{
		"enabled":        false,
		"worker_running": false,
		"is_healthy":     false,
	}

	if health != nil {
		if enabled, ok := health["enabled"]; ok {
			metrics["enabled"] = enabled
		}
		if running, ok := health["worker_running"]; ok {
			metrics["worker_running"] = running
		}
		if healthy, ok := health["is_healthy"]; ok {
			metrics["is_healthy"] = healthy
		}
	}

	// Add optional metrics if they exist
	if queueSize, ok := health["queue_size"]; ok {
		metrics["queue_size"] = queueSize
	}
	if queueCapacity, ok := health["queue_capacity"]; ok {
		metrics["queue_capacity"] = queueCapacity
	}
	if lastFlush, ok := health["last_flush_time"]; ok {
		metrics["last_flush_time"] = lastFlush
	}

	return metrics
}
