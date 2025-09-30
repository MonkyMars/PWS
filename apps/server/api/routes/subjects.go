package routes

import (
	"github.com/MonkyMars/PWS/api/internal"
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/gofiber/fiber/v3"
)

func SetupSubjectRoutes(app *fiber.App) {
	subjects := app.Group("/subjects")

	subjects.Get("/", internal.GetAllSubjects)
	subjects.Get("/me", middleware.AuthMiddleware(), internal.GetUserSubjects)
	subjects.Get("/:subjectId", internal.GetSubjectByID)
}
