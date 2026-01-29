package deadlines

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// GetOwnSubmission handles fetching the current student's submission for a specific deadline
// GET /deadlines/:id/submission
func (dr *DeadlineRoutes) GetOwnSubmission(c fiber.Ctx) error {
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to get user claims")
	}
	if claims.Role != "student" {
		return lib.HandleServiceError(c, lib.ErrInsufficientPermissions, "only students can fetch their own submission")
	}

	deadlineIDStr := c.Params("id")
	deadlineID, err := uuid.Parse(deadlineIDStr)
	if err != nil {
		return lib.HandleServiceError(c, err, "invalid deadline id")
	}

	submission, err := dr.deadlineService.GetSubmissionByStudent(deadlineID, claims.Sub)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to fetch submission")
	}

	return response.Success(c, submission)
}

// GetAllSubmissions handles fetching all student submissions for a specific deadline
// GET /deadlines/:id/submissions
func (dr *DeadlineRoutes) GetAllSubmissions(c fiber.Ctx) error {
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to get user claims")
	}
	if claims.Role != "teacher" && claims.Role != "admin" {
		return lib.HandleServiceError(c, lib.ErrInsufficientPermissions, "only teachers or admins can fetch all submissions")
	}

	deadlineIDStr := c.Params("id")
	deadlineID, err := uuid.Parse(deadlineIDStr)
	if err != nil {
		return lib.HandleServiceError(c, err, "invalid deadline id")
	}

	submissions, err := dr.deadlineService.GetAllSubmissionsForDeadline(deadlineID)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to fetch submissions")
	}

	return response.Success(c, submissions)
}
