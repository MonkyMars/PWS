package deadlines

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
	"github.com/MonkyMars/PWS/types"
)


func CreateDeadline(c fiber.Ctx) error {
	logger := config.SetupLogger()

	body, err := middleware.GetValidatedRequest[types.CreateDeadlineRequest](c)
	if err != nil {
		logger.Error("Failed to get validated request", "error", err)
		return lib.HandleValidationError(c, err, "request")
	}

	if body == nil {
		return response.NotFound(c, "Data not found")
	}

	deadlineService := services.NewDeadlineService()
	err = deadlineService.CreateDeadline(body)
	if err != nil {
		return response.InternalServerError(c, "Failed to create deadline: "+err.Error())
	}
	return response.Accepted(c, "Deadline creation accepted")
}
