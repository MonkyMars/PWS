package workers

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/workers"
	"github.com/gofiber/fiber/v3"
)

func (wr *WorkerRoutes) GetWorkerHealth(c fiber.Ctx) error {
	manager := workers.GetGlobalManager()
	if manager == nil {
		return lib.HandleServiceError(c, lib.ErrServiceUnavailable)
	}

	healthStatus := manager.HealthStatus()
	if healthStatus == nil {
		return lib.HandleServiceError(c, lib.ErrServiceUnavailable)
	}

	// Determine overall status
	isHealthy := false
	if healthVal, exists := healthStatus["is_healthy"]; exists && healthVal != nil {
		if healthy, ok := healthVal.(bool); ok {
			isHealthy = healthy
		}
	}
	if !isHealthy {
		return lib.HandleServiceError(c, lib.ErrServiceUnavailable)
	}

	return response.SuccessWithMessage(c, "Worker health status retrieved", healthStatus)
}

// getAuditWorkerHealth returns the health status of the audit worker
func (wr *WorkerRoutes) GetAuditWorkerHealth(c fiber.Ctx) error {
	healthStatus := workers.AuditHealthStatus()
	if healthStatus == nil {
		return lib.HandleServiceError(c, lib.ErrServiceUnavailable)
	}

	isHealthy := false
	if healthVal, exists := healthStatus["is_healthy"]; exists && healthVal != nil {
		if healthy, ok := healthVal.(bool); ok {
			isHealthy = healthy
		}
	}

	if !isHealthy {
		return lib.HandleServiceError(c, lib.ErrServiceUnavailable)
	}

	return response.SuccessWithMessage(c, "Audit worker health status retrieved", healthStatus)
}

// getHealthWorkerHealth returns the health status of the health monitoring worker
func (wr *WorkerRoutes) GetHealthWorkerHealth(c fiber.Ctx) error {
	healthStatus := workers.ServiceHealthStatus()
	if healthStatus == nil {
		return lib.HandleServiceError(c, lib.ErrServiceUnavailable)
	}

	isHealthy := false
	if healthVal, exists := healthStatus["is_healthy"]; exists && healthVal != nil {
		if healthy, ok := healthVal.(bool); ok {
			isHealthy = healthy
		}
	}

	if !isHealthy {
		return lib.HandleServiceError(c, lib.ErrServiceUnavailable)
	}

	return response.SuccessWithMessage(c, "Health worker status retrieved", healthStatus)
}

// getCleanupWorkerHealth returns the health status of the cleanup worker
func (wr *WorkerRoutes) GetCleanupWorkerHealth(c fiber.Ctx) error {
	manager := workers.GetGlobalManager()
	if manager == nil {
		return lib.HandleServiceError(c, lib.ErrServiceUnavailable)
	}

	healthStatus := manager.HealthStatus()
	if healthStatus == nil {
		return lib.HandleServiceError(c, lib.ErrServiceUnavailable)
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
		return lib.HandleServiceError(c, lib.ErrServiceUnavailable)
	}

	return response.SuccessWithMessage(c, "Cleanup worker health status retrieved", cleanupStatus)
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
		return lib.HandleServiceError(c, lib.ErrMissingParameter)
	}

	stats, err := workers.GetServiceStats(serviceName)
	if err != nil {
		return lib.HandleServiceError(c, err)
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
		return lib.HandleServiceError(c, err)
	}

	return response.Accepted(c, "Cleanup triggered successfully")
}
