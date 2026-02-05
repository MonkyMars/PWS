package deadlines

import (
	"net/http"
	"time"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// CreateOrUpdateSubmission handles student submission (create or update) for a deadline
// POST /deadlines/:id/submission
func (dr *DeadlineRoutes) CreateOrUpdateSubmission(c fiber.Ctx) error {
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to get user claims")
	}
	if claims.Role != "student" {
		return lib.HandleServiceError(c, lib.ErrInsufficientPermissions, "only students can submit to deadlines")
	}

	deadlineIDStr := c.Params("id")
	deadlineID, err := uuid.Parse(deadlineIDStr)
	if err != nil {
		return lib.HandleServiceError(c, err, "invalid deadline id")
	}

	var req types.CreateSubmissionRequest
	if err := c.Bind().Body(&req); err != nil {
		return lib.HandleServiceError(c, err, "failed to parse submission request")
	}

	// Validate file count
	if len(req.FileIDs) == 0 {
		return response.BadRequest(c, "At least one file must be submitted")
	}
	if len(req.FileIDs) > 10 {
		return response.BadRequest(c, "You can submit a maximum of 10 files")
	}

	// Get current time for timestamps
	now := time.Now().UTC().Format(time.RFC3339)

	// Call service to create or update submission
	submission, err := dr.deadlineService.CreateOrUpdateSubmission(deadlineID, claims.Sub, req, now)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to create or update submission")
	}

	// TODO: Notify teachers/admins of new/updated submission

	return c.Status(http.StatusAccepted).JSON(submission)
}
