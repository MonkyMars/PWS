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
	FileIDs    []string  `json:"file_ids"` // Google Drive file IDs
	Message    string    `json:"message"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
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
