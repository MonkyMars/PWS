package services

import (
	"fmt"
	"time"

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

	fmt.Println("Filter Options:", filterOptions)

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
	// Submission-related
	CreateOrUpdateSubmission(deadlineID, studentID uuid.UUID, req types.CreateSubmissionRequest, now string) (*types.SubmissionResponse, error)
	GetSubmissionByStudent(deadlineID, studentID uuid.UUID) (*types.SubmissionResponse, error)
	GetAllSubmissionsForDeadline(deadlineID uuid.UUID) ([]*types.SubmissionResponse, error)
}

// CreateOrUpdateSubmission creates or updates a student's submission for a deadline
func (ds *DeadlineService) CreateOrUpdateSubmission(deadlineID, studentID uuid.UUID, req types.CreateSubmissionRequest, now string) (*types.SubmissionResponse, error) {
	// Fetch the deadline to get due_date
	deadline, err := ds.getDeadlineByID(deadlineID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deadline: %w", err)
	}
	if deadline == nil {
		return nil, fmt.Errorf("deadline not found")
	}

	// Check if a submission already exists
	query := Query().
		SetOperation("select").
		SetTable("submissions").
		SetLimit(1)
	query.Where = map[string]any{
		"deadline_id": deadlineID,
		"student_id":  studentID,
	}

	result, err := database.ExecuteQuery[types.Submission](query)
	if err != nil {
		return nil, fmt.Errorf("failed to query submission: %w", err)
	}

	var submission types.Submission
	isUpdate := false
	if len(result.Data) > 0 {
		// Update existing submission
		isUpdate = true
		submission = result.Data[0]
		updateQuery := Query().
			SetOperation("update").
			SetTable("submissions").
			SetData(map[string]any{
				"file_ids":   req.FileIDs,
				"message":    req.Message,
				"updated_at": now,
			})
		updateQuery.Where = map[string]any{
			"public.submissions.id": submission.ID,
		}
		_, err := database.ExecuteQuery[types.Submission](updateQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to update submission: %w", err)
		}
		// Update local struct for response
		submission.FileIDs = req.FileIDs
		submission.Message = req.Message
		submission.UpdatedAt = now
	} else {
		// Insert new submission
		newID := uuid.New()
		insertQuery := Query().
			SetOperation("insert").
			SetTable("submissions").
			SetData(map[string]any{
				"id":          newID,
				"deadline_id": deadlineID,
				"student_id":  studentID,
				"file_ids":    req.FileIDs,
				"message":     req.Message,
				"created_at":  now,
				"updated_at":  now,
			})
		_, err := database.ExecuteQuery[types.Submission](insertQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to insert submission: %w", err)
		}
		submission = types.Submission{
			ID:         newID,
			DeadlineID: deadlineID,
			StudentID:  studentID,
			FileIDs:    req.FileIDs,
			Message:    req.Message,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
	}

	// Calculate late/updated flags
	isLate := false
	isUpdated := false
	dueDate, err := parseTime(deadline.DueDate)
	if err == nil {
		createdAt, _ := parseTime(submission.CreatedAt)
		updatedAt, _ := parseTime(submission.UpdatedAt)
		if createdAt.After(dueDate) {
			isLate = true
		}
		if updatedAt.After(dueDate) && updatedAt != createdAt {
			isUpdated = true
		}
	}

	resp := &types.SubmissionResponse{
		ID:         submission.ID,
		DeadlineID: submission.DeadlineID,
		StudentID:  submission.StudentID,
		FileIDs:    submission.FileIDs,
		Message:    submission.Message,
		CreatedAt:  submission.CreatedAt,
		UpdatedAt:  submission.UpdatedAt,
		IsLate:     isLate,
		IsUpdated:  isUpdated,
	}

	// --- Notification logic for teachers/admins ---
	// Find all teachers/admins for the subject of this deadline
	// For now, just log the notification; replace with your actual notification system as needed

	// Fetch subject teachers
	subjectTeachers, err := ds.getTeachersForSubject(deadline.SubjectID)
	if err == nil {
		for _, teacher := range subjectTeachers {
			ds.Logger.Info("Notify teacher of new/updated submission",
				"teacher_id", teacher.Id,
				"student_id", studentID,
				"deadline_id", deadlineID,
				"is_update", isUpdate,
			)
			// TODO: Integrate with actual notification system (email, in-app, etc.)
		}
	}
	// Optionally, notify admins as well (not implemented here, but can be added similarly)

	return resp, nil
}

// GetAllSubmissionsForDeadline fetches all student submissions for a specific deadline
func (ds *DeadlineService) GetAllSubmissionsForDeadline(deadlineID uuid.UUID) ([]*types.SubmissionResponse, error) {
	// Fetch the deadline to get due_date
	deadline, err := ds.getDeadlineByID(deadlineID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deadline: %w", err)
	}
	if deadline == nil {
		return []*types.SubmissionResponse{}, nil
	}

	query := Query().
		SetOperation("select").
		SetTable("submissions")
	query.Where = map[string]any{
		"submissions.deadline_id": deadlineID,
	}
	result, err := database.ExecuteQuery[types.Submission](query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch submissions: %w", err)
	}

	dueDate, _ := parseTime(deadline.DueDate)
	var responses []*types.SubmissionResponse
	for _, s := range result.Data {
		isLate := false
		isUpdated := false
		createdAt, _ := parseTime(s.CreatedAt)
		updatedAt, _ := parseTime(s.UpdatedAt)
		if createdAt.After(dueDate) {
			isLate = true
		}
		if updatedAt.After(dueDate) && updatedAt != createdAt {
			isUpdated = true
		}
		responses = append(responses, &types.SubmissionResponse{
			ID:         s.ID,
			DeadlineID: s.DeadlineID,
			StudentID:  s.StudentID,
			FileIDs:    s.FileIDs,
			Message:    s.Message,
			CreatedAt:  s.CreatedAt,
			UpdatedAt:  s.UpdatedAt,
			IsLate:     isLate,
			IsUpdated:  isUpdated,
		})
	}
	return responses, nil
}

// GetSubmissionByStudent fetches a student's submission for a specific deadline
func (ds *DeadlineService) GetSubmissionByStudent(deadlineID, studentID uuid.UUID) (*types.SubmissionResponse, error) {
	// Fetch the deadline to get due_date
	deadline, err := ds.getDeadlineByID(deadlineID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deadline: %w", err)
	}
	if deadline == nil {
		return nil, nil
	}

	query := Query().
		SetOperation("select").
		SetTable("submissions").
		SetLimit(1)

	query.Where = map[string]any{
		"submissions.deadline_id": deadlineID,
		"student_id":              studentID,
	}
	result, err := database.ExecuteQuery[types.Submission](query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch submission: %w", err)
	}
	if len(result.Data) == 0 {
		return nil, nil
	}
	s := result.Data[0]

	dueDate, _ := parseTime(deadline.DueDate)
	isLate := false
	isUpdated := false
	createdAt, _ := parseTime(s.CreatedAt)
	updatedAt, _ := parseTime(s.UpdatedAt)
	if createdAt.After(dueDate) {
		isLate = true
	}
	if updatedAt.After(dueDate) && updatedAt != createdAt {
		isUpdated = true
	}

	resp := &types.SubmissionResponse{
		ID:         s.ID,
		DeadlineID: s.DeadlineID,
		StudentID:  s.StudentID,
		FileIDs:    s.FileIDs,
		Message:    s.Message,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
		IsLate:     isLate,
		IsUpdated:  isUpdated,
	}
	return resp, nil
}

func (ds *DeadlineService) getDeadlineByID(deadlineID uuid.UUID) (*types.Deadline, error) {
	query := Query().
		SetOperation("select").
		SetTable("deadlines").
		SetLimit(1)
	query.Where = map[string]any{
		"public.deadlines.id": deadlineID,
	}
	result, err := database.ExecuteQuery[types.Deadline](query)
	if err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, nil
	}
	return &result.Data[0], nil
}

func (ds *DeadlineService) getTeachersForSubject(subjectID uuid.UUID) ([]types.User, error) {
	query := Query().
		SetOperation("select").
		SetTable("users")
	query.Where = map[string]any{
		"role": "teacher",
	}
	// Assuming there's a subject_teachers table mapping subjects to their teachers
	subjectTeacherQuery := Query().
		SetOperation("select").
		SetTable("subject_teachers")
	subjectTeacherQuery.Where = map[string]any{
		"subject_id": subjectID,
	}
	subjectTeachersResult, err := database.ExecuteQuery[types.Teacher](subjectTeacherQuery)
	if err != nil {
		return nil, err
	}
	var teacherIDs []uuid.UUID
	for _, st := range subjectTeachersResult.Data {
		teacherIDs = append(teacherIDs, st.Id)
	}
	if len(teacherIDs) == 0 {
		return []types.User{}, nil
	}
	query.Where["id"] = teacherIDs

	result, err := database.ExecuteQuery[types.User](query)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}
