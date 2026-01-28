package deadlines

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/gofiber/fiber/v3"
)

// FetchDeadlines handles fetching all deadlines
// GET /deadlines/me
func (dr *DeadlineRoutes) FetchDeadlinesForUser(c fiber.Ctx) error {
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to get user claims")
	}

	if claims.Role == "user" {
		deadlines, err := dr.deadlineService.FetchDeadlinesByUser(claims.Sub)
		if err != nil {
			return lib.HandleServiceError(c, err, "failed to fetch deadlines for user")
		}

		return response.Success(c, deadlines)
	}

	deadlines, err := dr.deadlineService.FetchAllDeadlines()
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to fetch deadlines")
	}

	return response.Success(c, deadlines)
}
