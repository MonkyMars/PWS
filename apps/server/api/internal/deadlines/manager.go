package deadlines

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
)

// HealthRoutes handles HTTP routing for health-related endpoints.
// It follows clean architecture principles by depending on interfaces rather than concrete implementations.
// This makes the code more testable and maintainable.
type DeadlineRoutes struct {
	subjectService  services.SubjectServiceInterface
	deadlineService services.DeadlineServiceInterface
	middleware      *middleware.Middleware
	logger          *config.Logger
}

// NewAuthRoutesWithDefaults creates an AuthRoutes instance with default dependencies.
// This is a convenience constructor for production use where you want to use
// the default implementations of all services.
func NewDeadlineRoutesWithDefaults() *DeadlineRoutes {
	return &DeadlineRoutes{
		subjectService:  services.NewSubjectService(),
		deadlineService: services.NewDeadlineService(),
		middleware:      middleware.NewMiddleware(),
		logger:          config.SetupLogger(),
	}
}

// This method organizes routes logically and follows RESTful conventions.
// It groups related functionality and applies appropriate middleware.
func (dr *DeadlineRoutes) RegisterRoutes(app *fiber.App) {
	deadlines := app.Group("/deadlines", dr.middleware.AuthMiddleware())

	deadlines.Post("/", dr.middleware.RoleMiddleware(lib.RoleAdmin, lib.RoleTeacher), dr.CreateDeadline)
	deadlines.Get("/me", dr.FetchDeadlinesForUser)
	deadlines.Put("/:id", dr.UpdateDeadlineById)
	deadlines.Delete("/:id", dr.DeleteDeadlineById)
	deadlines.Delete("/user/:user_id", dr.DeleteDeadlinesByUser)
}
