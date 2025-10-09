// Package routes provides HTTP handlers for worker health monitoring endpoints.
// This file contains routes that expose the health status and metrics of background workers
// including audit logging, health monitoring, and cleanup workers.
package routes

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupWorkerRoutes configures worker monitoring endpoints
func (r *Router) SetupWorkerRoutes(app *fiber.App) {
	// Worker health monitoring routes
	workerGroup := app.Group("/workers", middleware.AdminMiddleware())

	// Overall worker health status
	workerGroup.Get("/health", r.WorkerRoutes.GetWorkerHealth)

	// Individual worker health status
	workerGroup.Get("/audit/health", r.WorkerRoutes.GetAuditWorkerHealth)
	workerGroup.Get("/health-monitor/health", r.WorkerRoutes.GetHealthWorkerHealth)
	workerGroup.Get("/cleanup/health", r.WorkerRoutes.GetCleanupWorkerHealth)

	// Worker metrics and statistics
	workerGroup.Get("/metrics", r.WorkerRoutes.GetWorkerMetrics)
	workerGroup.Get("/audit/metrics", r.WorkerRoutes.GetAuditWorkerMetrics)
	workerGroup.Get("/health-monitor/metrics", r.WorkerRoutes.GetHealthWorkerMetrics)

	// Health monitoring service information
	workerGroup.Get("/health-monitor/services", r.WorkerRoutes.GetMonitoredServices)
	workerGroup.Get("/health-monitor/services/:service", r.WorkerRoutes.GetServiceStatistics)

	// Administrative actions
	workerGroup.Post("/cleanup/trigger", r.WorkerRoutes.TriggerCleanup)
}
