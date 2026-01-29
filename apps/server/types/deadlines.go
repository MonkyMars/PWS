package types

import (
	"github.com/google/uuid"
)

type CreateDeadlineRequest struct {
	SubjectID   uuid.UUID `json:"subject_id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	CreatedAt   string    `json:"created_at"`
}

type Deadline struct {
	ID          uuid.UUID `json:"id"`
	SubjectID   uuid.UUID `json:"subject_id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type Submission struct {
	ID         uuid.UUID `json:"id"`
	DeadlineID uuid.UUID `json:"deadline_id"`
	StudentID  uuid.UUID `json:"student_id"`
	FileIDs    []string  `json:"file_ids" pg:"file_ids,type:text[]"` // Google Drive file IDs
	Message    string    `json:"message"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

// Used for creating/updating a submission
type CreateSubmissionRequest struct {
	FileIDs []string `json:"file_ids" pg:"file_ids,type:text[]"` // Google Drive file IDs
	Message string   `json:"message"`
}

// Used for returning a submission to the client
type SubmissionResponse struct {
	ID         uuid.UUID `json:"id"`
	DeadlineID uuid.UUID `json:"deadline_id"`
	StudentID  uuid.UUID `json:"student_id"`
	FileIDs    []string  `json:"file_ids" pg:"file_ids,type:text[]"`
	Message    string    `json:"message"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
	IsLate     bool      `json:"is_late"`
	IsUpdated  bool      `json:"is_updated"`
}

type DeadlineWithSubject struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	Subject     Subject   `json:"subject"`
}
