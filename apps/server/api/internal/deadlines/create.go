package deadlines

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

func (dr *DeadlineRoutes) CreateDeadline(c fiber.Ctx) error {
	body, err := middleware.GetValidatedRequest[types.CreateDeadlineRequest](c)
	if err != nil {
		msg := "Failed to parse and validate deadline creation request body"
		return lib.HandleServiceError(c, err, msg)
	}

	if body == nil {
		return response.NotFound(c, "Data not found")
	}

	err = dr.deadlineService.CreateDeadline(body)
	if err != nil {
		return response.InternalServerError(c, "Failed to create deadline: "+err.Error())
	}
	return response.Accepted(c, "Deadline creation accepted")
}
