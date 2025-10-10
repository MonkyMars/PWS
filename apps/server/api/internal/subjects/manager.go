package subjects

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
)

// HealthRoutes handles HTTP routing for health-related endpoints.
// It follows clean architecture principles by depending on interfaces rather than concrete implementations.
// This makes the code more testable and maintainable.
type SubjectRoutes struct {
	subjectService services.SubjectServiceInterface
}

// NewAuthRoutesWithDefaults creates an AuthRoutes instance with default dependencies.
// This is a convenience constructor for production use where you want to use
// the default implementations of all services.
func NewSubjectRoutesWithDefaults() *SubjectRoutes {
	return &SubjectRoutes{
		subjectService: services.NewSubjectService(),
	}
}

// This method organizes routes logically and follows RESTful conventions.
// It groups related functionality and applies appropriate middleware.
func (sr *SubjectRoutes) RegisterRoutes(app *fiber.App) {
	subjects := app.Group("/subjects")

	subjects.Get("/", sr.GetAllSubjects)
	subjects.Get("/me", middleware.AuthMiddleware(), sr.GetUserSubjects)
	subjects.Get("/:subjectId", sr.GetSubjectByID)
}
