package deadlines

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

func (dr *DeadlineRoutes) UpdateDeadlineById(c fiber.Ctx) error {
	deadlineId := c.Params("id")
	if deadlineId == "" {
		return lib.HandleServiceError(c, nil, "deadline id parameter is required")
	}

	var updateData types.Deadline
	if err := c.Bind().Body(&updateData); err != nil {
		return lib.HandleServiceError(c, err, "failed to parse request body")
	}

	err := dr.deadlineService.UpdateDeadlineById(deadlineId, updateData)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to update deadline")
	}

	return response.NoContent(c)
}
