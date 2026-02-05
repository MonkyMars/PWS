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

	filterOptions, err := lib.GetQueryParams(c, map[string]bool{
		"due_date_from": false,
		"due_date_to":   false,
		"subject_id":    false,
	})
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to get filter options")
	}

	dr.logger.Info("Fetching deadlines for user", "userID", claims.Sub, "role", claims.Role)

	if claims.Role == "student" {
		deadlines, err := dr.deadlineService.FetchDeadlinesByUser(claims.Sub, filterOptions)
		if err != nil {
			return lib.HandleServiceError(c, err, "failed to fetch deadlines for user")
		}

		return response.Success(c, deadlines)
	}

	deadlines, err := dr.deadlineService.FetchAllDeadlines(filterOptions)
	if err != nil {
		return lib.HandleServiceError(c, err, "failed to fetch deadlines")
	}

	return response.Success(c, deadlines)
}
