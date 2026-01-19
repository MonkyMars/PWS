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
