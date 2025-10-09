package health

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
)

func (hr *HealthRoutes) GetLogs(c fiber.Ctx) error {
	auditService := services.NewAuditService()
	logs, err := auditService.GetLogs()
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve audit logs")
	}
	return response.Success(c, logs)
}
