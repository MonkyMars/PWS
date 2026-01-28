package deadlines

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// DeleteDeadline handles deleting a specific deadline by ID
// DELETE /deadlines/:id
func (dr *DeadlineRoutes) DeleteDeadlineById(c fiber.Ctx) error {
	deadlineId := c.Params("id")
	if deadlineId == "" {
		return lib.HandleServiceError(c, nil, "deadline id parameter is required")
	}

	err := dr.deadlineService.DeleteDeadlineById(deadlineId)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to delete deadline")
	}

	return response.NoContent(c)
}

// DeleteDeadlinesByUser handles deleting all deadlines for a specific user
// DELETE /deadlines/user/:user_id
func (dr *DeadlineRoutes) DeleteDeadlinesByUser(c fiber.Ctx) error {
	userId := c.Params("user_id")
	if userId == "" {
		return lib.HandleServiceError(c, nil, "user_id parameter is required")
	}

	userUuid, err := uuid.Parse(userId)
	if err != nil {
		return lib.HandleServiceError(c, err, "invalid user_id parameter")
	}

	err = dr.deadlineService.DeleteDeadlinesFromUser(userUuid)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to delete deadlines for user")
	}

	return response.NoContent(c)
}
