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
	subjectID, err := lib.GetParams(c, "subjectId")
	if err != nil {
		return lib.HandleServiceError(c, err)
	}

	subject, err := sr.subjectService.GetSubjectByID(subjectID["subjectId"])
	if err != nil {
		return lib.HandleServiceError(c, err)
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
	user := lib.GetUserFromContext(c)
	if user == nil {
		return response.Unauthorized(c, "You must be logged in to view your subjects")
	}

	var subjects []types.Subject
	if lib.HasPrivileges(c) {
		s, err := sr.subjectService.GetAllSubjects()
		if err != nil {
			sr.logger.Error("Failed to retrieve subjects", "error", err)
			return response.InternalServerError(c, "Failed to retrieve subjects")
		}
		subjects = s
	} else {
		s, err := sr.subjectService.GetUserSubjects(user.Id.String())
		if err != nil {
			sr.logger.Error("Failed to retrieve user subjects", "user_id", user.Id.String(), "error", err)
			return response.InternalServerError(c, "Failed to retrieve user subjects")
		}
		subjects = s
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

func (sr *SubjectRoutes) GetAllTeachers(c fiber.Ctx) error {
	teachers, err := sr.subjectService.GetAllTeachers()
	if err != nil {
		return lib.HandleServiceError(c, err)
	}

	return response.Success(c, teachers)
}
