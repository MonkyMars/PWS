package internal

import (
	"fmt"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
)

func GetLogs(c fiber.Ctx) error {
	auditService := services.NewAuditService()
	logs, err := auditService.GetLogs()
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve audit logs")
	}

	items := []any{}
	for _, log := range logs {
		items = append(items, log)
	}

	pages := len(items) / 50
	if len(items)%50 != 0 {
		pages++
	}

	c.Set("X-Total-Count", fmt.Sprintf("%d", len(items)))
	c.Set("X-Total-Pages", fmt.Sprintf("%d", pages))
	c.Set("X-Page-Size", "50")

	return response.Paginated(c, items, len(items), pages, 50)
}
