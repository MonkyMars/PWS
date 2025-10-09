package types

import "github.com/google/uuid"

type Subject struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Color       string    `json:"color"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	TeacherId   uuid.UUID `json:"teacher_id"`
	TeacherName string    `json:"teacher_name"`
	IsActive    bool      `json:"is_active"`
}
