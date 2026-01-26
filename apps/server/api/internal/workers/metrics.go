package workers

import (
	"time"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/workers"
	"github.com/gofiber/fiber/v3"
)

// getWorkerMetrics returns comprehensive metrics for all workers
func (wr *WorkerRoutes) GetWorkerMetrics(c fiber.Ctx) error {
	if wr.manager == nil {
		msg := "Worker manager not available for metrics retrieval"
		return lib.HandleServiceError(c, lib.ErrWorkerUnavailable, msg)
	}

	healthStatus := wr.manager.HealthStatus()
	if healthStatus == nil {
		msg := "Unable to retrieve worker metrics from manager"
		return lib.HandleServiceError(c, lib.ErrWorkerUnavailable, msg)
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
