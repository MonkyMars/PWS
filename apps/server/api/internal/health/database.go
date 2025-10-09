package health

import (
	"time"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

func (hr *HealthRoutes) GetDatabaseHealth(c fiber.Ctx) error {
	now := time.Now()
	if err := services.Ping(); err != nil {
		return response.ServiceUnavailable(c, "Database connection error: "+err.Error())
	}

	return response.Success(c, types.DatabaseHealthResponse{
		Status:  "ok",
		Message: "Database connection is healthy",
		Elapsed: time.Since(now).String(),
	})
}
