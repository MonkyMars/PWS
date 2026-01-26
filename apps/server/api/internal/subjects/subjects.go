package subjects

import (
	"fmt"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// GetSubjectByID retrieves a subject by its ID
func (sr *SubjectRoutes) GetSubjectByID(c fiber.Ctx) error {
	subjectID := c.Params("subjectId")

	if subjectID == "" {
		msg := "Missing required subjectId parameter in request"
		return lib.HandleServiceError(c, lib.ErrMissingField, msg)
	}

	subject, err := sr.subjectService.GetSubjectByID(subjectID["subjectId"])
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve subject for subject ID %s: %v", subjectID, err)
		return lib.HandleServiceError(c, err, msg)
	}

	if subject == nil {
		msg := fmt.Sprintf("Subject not found for subject ID %s", subjectID)
		return lib.HandleServiceError(c, lib.ErrSubjectNotFound, msg)
	}

	return response.Success(c, subject)
}

func (sr *SubjectRoutes) GetAllSubjects(c fiber.Ctx) error {
	subjects, err := sr.subjectService.GetAllSubjects()
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve all subjects: %v", err)
		return lib.HandleServiceError(c, err, msg)
	}

	return response.Success(c, subjects)
}

func (sr *SubjectRoutes) GetUserSubjects(c fiber.Ctx) error {
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		msg := "Failed to get authenticated user claims for user subjects retrieval"
		return lib.HandleServiceError(c, err, msg)
	}

	var subjects []types.Subject
	if lib.HasPrivileges(c) {
		s, err := sr.subjectService.GetAllSubjects()
		if err != nil {
			msg := fmt.Sprintf("Failed to retrieve all subjects for user ID %s with role %s: %v", claims.Sub.String(), claims.Role, err)
			return lib.HandleServiceError(c, err, msg)
		}
		subjects = s
	} else {
		s, err := sr.subjectService.GetUserSubjects(user.Id.String())
		if err != nil {
			msg := fmt.Sprintf("Failed to retrieve subjects for student user ID %s: %v", claims.Sub.String(), err)
			return lib.HandleServiceError(c, err, msg)
		}
		subjects = s
	default:
		msg := fmt.Sprintf("User ID %s with role %s does not have permission to view subjects", claims.Sub.String(), claims.Role)
		return lib.HandleServiceError(c, lib.ErrForbidden, msg)
	}

	return response.Success(c, subjects)
}

func (sr *SubjectRoutes) GetSubjectTeachers(c fiber.Ctx) error {
	subjectId, err := lib.GetParams(c, "subjectId")
	if err != nil {
		return lib.HandleServiceError(c, err)
	}

	teachers, err := sr.subjectService.GetSubjectTeachers(subjectId["subjectId"])
	if err != nil {
		return lib.HandleServiceError(c, err)
	}

	return response.Success(c, teachers)
}
