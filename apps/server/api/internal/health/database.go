package health

import (
	"fmt"
	"time"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

func (hr *HealthRoutes) GetDatabaseHealth(c fiber.Ctx) error {
	now := time.Now()
	if err := services.Ping(); err != nil {
		msg := fmt.Sprintf("Database health check failed: %v", err)
		return lib.HandleServiceError(c, lib.ErrDatabaseConnection, msg)
	}

	return response.Success(c, types.DatabaseHealthResponse{
		Status:  "ok",
		Message: "Database connection is healthy",
		Elapsed: time.Since(now).String(),
	})
}
