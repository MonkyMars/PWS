package workers

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/workers"
	"github.com/gofiber/fiber/v3"
)

// WorkerRoutes defines routes related to worker management and monitoring.
// It follows clean architecture principles by depending on interfaces rather than concrete implementations.
// This makes the code more testable and maintainable.
type WorkerRoutes struct {
	manager    workers.WorkerManagerInterface
	middleware *middleware.Middleware
}

// NewWorkerRoutes creates a new WorkerRoutes instance with dependency injection.
// In production, this will use the real services, but in tests,
// it can use mock implementations for better unit testing.
func NewWorkerRoutesWithDefaults() *WorkerRoutes {
	return &WorkerRoutes{
		manager:    workers.GetGlobalManager(),
		middleware: middleware.NewMiddleware(),
	}
}

// This method organizes routes logically and follows RESTful conventions.
// It groups related functionality and applies appropriate middleware.
func (wr *WorkerRoutes) RegisterRoutes(app *fiber.App) {
	// Worker health monitoring routes
	workerGroup := app.Group("/workers", wr.middleware.AdminMiddleware())

	// Overall worker health status
	workerGroup.Get("/health", wr.GetWorkerHealth)

	// Individual worker health status
	workerGroup.Get("/audit/health", wr.GetAuditWorkerHealth)
	workerGroup.Get("/health-monitor/health", wr.GetHealthWorkerHealth)
	workerGroup.Get("/cleanup/health", wr.GetCleanupWorkerHealth)

	// Worker metrics and statistics
	workerGroup.Get("/metrics", wr.GetWorkerMetrics)
	workerGroup.Get("/audit/metrics", wr.GetAuditWorkerMetrics)
	workerGroup.Get("/health-monitor/metrics", wr.GetHealthWorkerMetrics)

	// Health monitoring service information
	workerGroup.Get("/health-monitor/services", wr.GetMonitoredServices)
	workerGroup.Get("/health-monitor/services/:service", wr.GetServiceStatistics)

	// Administrative actions
	workerGroup.Post("/cleanup/trigger", wr.TriggerCleanup)
}
