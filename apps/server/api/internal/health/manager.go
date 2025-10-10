package health

import (
	"time"

	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
)

var (
	appStartTime       = time.Now()
	requestCount int64 = 0
)

// HealthRoutes handles HTTP routing for health-related endpoints.
// It follows clean architecture principles by depending on interfaces rather than concrete implementations.
// This makes the code more testable and maintainable.
type HealthRoutes struct {
	auditService services.AuditServiceInterface
}

// NewAuthRoutesWithDefaults creates an AuthRoutes instance with default dependencies.
// This is a convenience constructor for production use where you want to use
// the default implementations of all services.
func NewHealthRoutesWithDefaults() *HealthRoutes {
	return &HealthRoutes{
		auditService: services.NewAuditService(),
	}
}

// This method organizes routes logically and follows RESTful conventions.
// It groups related functionality and applies appropriate middleware.
func (hr *HealthRoutes) RegisterRoutes(app *fiber.App) {
	health := app.Group("/health")
	health.Get("/", hr.GetSystemHealth)
	health.Get("/database", hr.GetDatabaseHealth)
	health.Get("/logs", hr.GetLogs)
}
