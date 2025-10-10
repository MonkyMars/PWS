package subjects

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// GetSubjectByID retrieves a subject by its ID
func (sr *SubjectRoutes) GetSubjectByID(c fiber.Ctx) error {
	logger := config.SetupLogger()
	subjectID := c.Params("subjectId")

	if subjectID == "" {
		return response.BadRequest(c, "Subject ID is required")
	}

	subject, err := sr.subjectService.GetSubjectByID(subjectID)
	if err != nil {
		logger.Error("Failed to retrieve subject", "subject_id", subjectID, "error", err)
		return response.InternalServerError(c, "Failed to retrieve subject")
	}

	if subject == nil {
		return response.NotFound(c, "Subject not found")
	}

	return response.Success(c, subject)
}

func (sr *SubjectRoutes) GetAllSubjects(c fiber.Ctx) error {
	logger := config.SetupLogger()

	subjects, err := sr.subjectService.GetAllSubjects()
	if err != nil {
		logger.Error("Failed to retrieve subjects", "error", err)
		return response.InternalServerError(c, "Failed to retrieve subjects")
	}

	return response.Success(c, subjects)
}

func (sr *SubjectRoutes) GetUserSubjects(c fiber.Ctx) error {
	logger := config.SetupLogger()
	claimsInterface := c.Locals("claims")

	if claimsInterface == nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	// Type assert claims
	claims, ok := claimsInterface.(*types.AuthClaims)
	if claims == nil || !ok {
		return response.Unauthorized(c, "Unauthorized")
	}

	var subjects []types.Subject
	switch claims.Role {
	case lib.RoleAdmin, lib.RoleTeacher:
		s, err := sr.subjectService.GetAllSubjects()
		if err != nil {
			logger.Error("Failed to retrieve subjects", "error", err)
			return response.InternalServerError(c, "Failed to retrieve subjects")
		}
		subjects = s
	case lib.RoleStudent:
		s, err := sr.subjectService.GetUserSubjects(claims.Sub.String())
		if err != nil {
			logger.Error("Failed to retrieve user subjects", "user_id", claims.Sub.String(), "error", err)
			return response.InternalServerError(c, "Failed to retrieve user subjects")
		}
		subjects = s
	default:
		return response.Forbidden(c, "You do not have permission to view subjects")
	}

	return response.Success(c, subjects)
}
