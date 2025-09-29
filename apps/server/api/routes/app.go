package routes

import (
	"github.com/MonkyMars/PWS/api/internal"
	"github.com/MonkyMars/PWS/config"
	"github.com/gofiber/fiber/v3"
)

func SetupAppRoutes(app *fiber.App) {
	cfg := config.Get()
	if cfg.Environment == "development" {
		app.Get("/health", internal.GetSystemHealth)
		app.Get("/health/database", internal.GetDatabaseHealth)
		app.Get("/health/audit", internal.GetAuditHealth)
		app.Get("/logs", internal.GetLogs)
	}

	app.Use(internal.NotFoundHandler)
}
