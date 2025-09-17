package routes

import (
	"github.com/MonkyMars/PWS/api/internal"
	"github.com/gofiber/fiber/v3"
)

func SetupAppRoutes(app *fiber.App) {
	app.Get("/health", internal.GetSystemHealth)
	app.Use(internal.NotFoundHandler)
}
