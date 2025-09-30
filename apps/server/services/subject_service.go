package services

import (
	"fmt"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
)

type SubjectService struct {
	Logger *config.Logger
}

func NewSubjectService() *SubjectService {
	return &SubjectService{
		Logger: config.SetupLogger(),
	}
}

func (ss *SubjectService) GetSubjectByID(subjectID string) (any, error) {
	query := Query().SetOperation("select").SetTable("subjects").SetLimit(1)
	query.Where[fmt.Sprintf("public.%s.id", lib.TableSubjects)] = subjectID

	data, err := database.ExecuteQuery[types.Subject](query)
	if err != nil {
		ss.Logger.Error("Failed to retrieve subject", "subject_id", subjectID, "error", err)
		return nil, err
	}

	if len(data.Data) == 0 {
		return nil, nil
	}

	return data.Single, nil
}

func (ss *SubjectService) GetAllSubjects() ([]types.Subject, error) {
	query := Query().SetOperation("select").SetTable("subjects")

	data, err := database.ExecuteQuery[types.Subject](query)
	if err != nil {
		ss.Logger.Error("Failed to retrieve subjects", "error", err)
		return nil, err
	}

	return data.Data, nil
}

func (ss *SubjectService) GetSubjectsByIDs(subjectIDs []string) ([]types.Subject, error) {
	if len(subjectIDs) == 0 {
		return []types.Subject{}, nil
	}

	query := Query().SetOperation("select").SetTable("subjects")
	query.Where[fmt.Sprintf("public.%s.subject_id", lib.TableSubjects)] = subjectIDs

	data, err := database.ExecuteQuery[types.Subject](query)
	if err != nil {
		ss.Logger.Error("Failed to retrieve subjects", "subject_ids", subjectIDs, "error", err)
		return nil, err
	}

	return data.Data, nil
}

func (ss *SubjectService) GetSubjectsByTeacherID(teacherID string) ([]types.Subject, error) {
	query := Query().SetOperation("select").SetTable(lib.TableSubjects)
	query.Where[fmt.Sprintf("public.%s.teacher_id", lib.TableSubjects)] = teacherID

	data, err := database.ExecuteQuery[types.Subject](query)
	if err != nil {
		ss.Logger.Error("Failed to retrieve subjects", "teacher_id", teacherID, "error", err)
		return nil, err
	}

	return data.Data, nil
}

func (ss *SubjectService) GetUserSubjects(userID string) ([]types.Subject, error) {
	query := Query().SetRawSQL(`
			SELECT s.id, s.name, s.code, s.color, s.created_at, s.updated_at, s.teacher_id, s.teacher_name
			FROM subjects s
			JOIN user_subjects us ON s.id = us.subject_id
			WHERE us.user_id = ? AND s.is_active = true
			ORDER BY s.name ASC
		`, userID)
	query.Where[fmt.Sprintf("public.%s.user_id", lib.TableUserSubjects)] = userID

	userSubjects, err := database.ExecuteQuery[types.Subject](query)
	if err != nil {
		ss.Logger.Error("Failed to retrieve user subjects", "user_id", userID, "error", err)
		return nil, err
	}

	if len(userSubjects.Data) == 0 {
		return []types.Subject{}, nil
	}

	return userSubjects.Data, nil
}
