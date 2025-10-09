package routes

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) SetupSubjectRoutes(app *fiber.App) {
	subjects := app.Group("/subjects")

	subjects.Get("/", r.SubjectRoutes.GetAllSubjects)
	subjects.Get("/me", middleware.AuthMiddleware(), r.SubjectRoutes.GetUserSubjects)
	subjects.Get("/:subjectId", r.SubjectRoutes.GetSubjectByID)
}
