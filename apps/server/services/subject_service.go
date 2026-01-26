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
	query := Query().SetOperation("select").SetTable("subjects").SetLimit(1).SetSelect([]string{
		"id", "name", "code", "color", "created_at", "updated_at", "teacher_id", "teacher_name",
	})
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
	query := Query().SetOperation("select").SetTable(lib.TableSubjects).SetSelect([]string{
		"id", "name", "code", "color", "created_at", "updated_at",
	})

	data, err := database.ExecuteQuery[types.Subject](query)
	if err != nil {
		ss.Logger.Error("Failed to retrieve subjects", "error", err)
		return nil, err
	}

	return data.Data, nil
}

func (ss *SubjectService) GetUserSubjects(userID string) ([]types.Subject, error) {
	// Raw SQL query to join subjects and user_subjects tables
	query := Query().SetRawSQL(`
		SELECT s.id, s.name, s.code, s.color, s.created_at, s.updated_at
		FROM subjects s
		JOIN user_subjects us ON s.id = us.subject_id
		WHERE us.user_id = ? AND s.is_active = true
		ORDER BY s.name ASC
	`, userID)

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

func (ss *SubjectService) GetSubjectTeachers(subjectID string) ([]types.User, error) {
	query := Query().SetRawSQL(`
			SELECT u.id, u.username, u.email, u.role, u.created_at
			FROM users u
			JOIN subject_teachers st ON u.id = st.user_id
			WHERE st.subject_id = ?
		`, subjectID)

	data, err := database.ExecuteQuery[types.User](query)
	if err != nil {
		ss.Logger.Error("Failed to retrieve subject teachers", "subject_id", subjectID, "error", err)
		return nil, err
	}

	return data.Data, nil
}

type SubjectServiceInterface interface {
	GetSubjectByID(subjectID string) (any, error)
	GetAllSubjects() ([]types.Subject, error)
	GetUserSubjects(userID string) ([]types.Subject, error)
	GetSubjectTeachers(subjectID string) ([]types.User, error)
}
