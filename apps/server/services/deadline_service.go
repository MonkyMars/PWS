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

func (ds *DeadlineService) FetchDeadlinesByUser(userId uuid.UUID) ([]types.Deadline, error) {
	query := Query().SetOperation("select").SetTable("deadlines")
	query.Where = map[string]any{
		"owner_id": userId,
	}

	deadlines, err := database.ExecuteQuery[types.Deadline](query)
	if err != nil {
		return nil, err
	}

	data := deadlines.Data

	return data, nil
}

func (ds *DeadlineService) FetchAllDeadlines() ([]types.Deadline, error) {
	query := Query().SetOperation("select").SetTable("deadlines")

	deadlines, err := database.ExecuteQuery[types.Deadline](query)
	if err != nil {
		return nil, err
	}

	data := deadlines.Data

	return data, nil
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

	query.Data = data

	_, err := database.ExecuteQuery[any](query)
	if err != nil {
		return err
	}

	return nil
}

// DeadlineServiceInterface defines the methods that the DeadlineService must implement.
// This interface is used for dependency injection and to facilitate testing.
type DeadlineServiceInterface interface {
	CreateDeadline(req *types.CreateDeadlineRequest) error
	FetchDeadlinesByUser(userId uuid.UUID) ([]types.Deadline, error)
	DeleteDeadlineById(deadlineId string) error
	DeleteDeadlinesFromUser(userId uuid.UUID) error
	FetchAllDeadlines() ([]types.Deadline, error)
	UpdateDeadlineById(deadlineId string, updateData types.Deadline) error
}
