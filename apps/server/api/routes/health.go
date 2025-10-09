package routes

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) SetupHealthRoutes(app *fiber.App) {
	health := app.Group("/health", middleware.AdminMiddleware())
	health.Get("/", r.AppRoutes.GetSystemHealth)
	health.Get("/database", r.AppRoutes.GetDatabaseHealth)
	health.Get("/logs", r.AppRoutes.GetLogs)
}
