package services

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/types"
)

type DeadlineService struct {
	Logger *config.Logger
}

func NewDeadlineService() *DeadlineService {
	return &DeadlineService{
		Logger: config.SetupLogger(),
	}
}

func (ds *DeadlineService) CreateDeadline(req *types.CreateDeadlineRequest) error {
	if req.SubjectID == uuid.Nil {
		return fmt.Errorf("subject_id is required")
	}
	if req.OwnerID == uuid.Nil {
		return fmt.Errorf("owner_id is required")
	}
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if req.Description == "" {
		return fmt.Errorf("description is required")
	}
	if req.DueDate == "" {
		return fmt.Errorf("due_date is required")
	}
	if req.CreatedAt == "" {
		return fmt.Errorf("created_at is required")
	}

	query := Query().SetOperation("insert").SetTable("deadlines")
	query.Data = map[string]any{
		"subject_id":  req.SubjectID,
		"owner_id":    req.OwnerID,
		"title":       req.Title,
		"description": req.Description,
		"due_date":    req.DueDate,
		"created_at":  req.CreatedAt,
	}

	_, err := database.ExecuteQuery[any](query)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DeadlineService) FetchDeadlinesByUser(userId uuid.UUID, filterOptions map[string]string) ([]types.DeadlineWithSubject, error) {
	var (
		query = `
			SELECT
				d.id, d.owner_id, d.title, d.description, d.due_date, d.created_at, d.updated_at,
				s.id AS subject__id, s.name AS subject__name, s.code AS subject__code, s.color AS subject__color,
				s.created_at AS subject__created_at, s.updated_at AS subject__updated_at,
				s.teacher_id AS subject__teacher_id, s.teacher_name AS subject__teacher_name, s.is_active AS subject__is_active
			FROM deadlines d
			LEFT JOIN subjects s ON d.subject_id = s.id
		`
		conditions []string
		args       []any
	)

	// Always filter by owner_id
	conditions = append(conditions, "d.owner_id = ?")
	args = append(args, userId)

	if subjectID, ok := filterOptions["subject_id"]; ok && subjectID != "" {
		conditions = append(conditions, "s.id = ?")
		args = append(args, subjectID)
	}
	if dueDateFrom, ok := filterOptions["due_date_from"]; ok && dueDateFrom != "" {
		conditions = append(conditions, "d.due_date >= ?")
		args = append(args, dueDateFrom)
	}
	if dueDateTo, ok := filterOptions["due_date_to"]; ok && dueDateTo != "" {
		conditions = append(conditions, "d.due_date <= ?")
		args = append(args, dueDateTo)
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " LIMIT 50;"

	deadlines, err := database.Raw[types.DeadlineWithSubject](query, args...)
	if err != nil {
		return nil, err
	}

	if deadlines.Count == 0 || deadlines.Data == nil {
		return []types.DeadlineWithSubject{}, nil
	}

	return deadlines.Data, nil
}

func (ds *DeadlineService) FetchAllDeadlines(filterOptions map[string]string) ([]types.DeadlineWithSubject, error) {
	var (
		query = `
			SELECT
				d.id, d.owner_id, d.title, d.description, d.due_date, d.created_at, d.updated_at,
				s.id AS subject__id, s.name AS subject__name, s.code AS subject__code, s.color AS subject__color,
				s.created_at AS subject__created_at, s.updated_at AS subject__updated_at,
				s.teacher_id AS subject__teacher_id, s.teacher_name AS subject__teacher_name, s.is_active AS subject__is_active
			FROM deadlines d
			LEFT JOIN subjects s ON d.subject_id = s.id
		`
		conditions []string
		args       []any
	)

	if subjectID, ok := filterOptions["subject_id"]; ok && subjectID != "" {
		conditions = append(conditions, "s.id = ?")
		args = append(args, subjectID)
	}
	if dueDateFrom, ok := filterOptions["due_date_from"]; ok && dueDateFrom != "" {
		conditions = append(conditions, "d.due_date >= ?")
		args = append(args, dueDateFrom)
	}
	if dueDateTo, ok := filterOptions["due_date_to"]; ok && dueDateTo != "" {
		conditions = append(conditions, "d.due_date <= ?")
		args = append(args, dueDateTo)
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " LIMIT 100;"

	deadlines, err := database.Raw[types.DeadlineWithSubject](query, args...)
	if err != nil {
		return nil, err
	}

	if deadlines.Count == 0 || deadlines.Data == nil {
		return []types.DeadlineWithSubject{}, nil
	}

	return deadlines.Data, nil
}

func (ds *DeadlineService) DeleteDeadlineById(deadlineId string) error {
	query := Query().SetOperation("delete").SetTable("deadlines")
	query.Where = map[string]any{
		"id": deadlineId,
	}

	_, err := database.ExecuteQuery[any](query)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DeadlineService) DeleteDeadlinesFromUser(userId uuid.UUID) error {
	query := Query().SetOperation("delete").SetTable("deadlines")
	query.Where = map[string]any{
		"owner_id": userId,
	}

	_, err := database.ExecuteQuery[any](query)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DeadlineService) UpdateDeadlineById(deadlineId string, updateData types.Deadline) error {
	query := Query().SetOperation("update").SetTable("deadlines")
	query.Where = map[string]any{
		"id": deadlineId,
	}
	data := map[string]any{}

	if updateData.Title != "" {
		data["title"] = updateData.Title
	}
	if updateData.Description != "" {
		data["description"] = updateData.Description
	}
	if updateData.DueDate != "" {
		data["due_date"] = updateData.DueDate
	}

	_, err := database.ExecuteQuery[any](query.SetData(data))
	if err != nil {
		return err
	}

	return nil
}

// DeadlineServiceInterface defines the methods that the DeadlineService must implement.
// This interface is used for dependency injection and to facilitate testing.
type DeadlineServiceInterface interface {
	CreateDeadline(req *types.CreateDeadlineRequest) error
	FetchDeadlinesByUser(userId uuid.UUID, filterOptions map[string]string) ([]types.DeadlineWithSubject, error)
	DeleteDeadlineById(deadlineId string) error
	DeleteDeadlinesFromUser(userId uuid.UUID) error
	FetchAllDeadlines(filterOptions map[string]string) ([]types.DeadlineWithSubject, error)
	UpdateDeadlineById(deadlineId string, updateData types.Deadline) error
}
