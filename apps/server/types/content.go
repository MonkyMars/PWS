package types

import (
	"time"

	"github.com/google/uuid"
)

type DriveFile struct {
	FileID   string `json:"file_id"`
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
}

type UploadSingleFileRequest struct {
	File      DriveFile `json:"file"`
	SubjectID uuid.UUID `json:"subject_id"`
}

type UploadMultipleFilesRequest struct {
	Files     []DriveFile `json:"files"`
	SubjectID uuid.UUID   `json:"subject_id"`
}

type File struct {
	Id         uuid.UUID `json:"id"`
	FileID     string    `json:"file_id"`
	Name       string    `json:"name"`
	MimeType   string    `json:"mime_type"`
	SubjectID  uuid.UUID `json:"subject_id"`
	UploadedBy uuid.UUID `json:"uploaded_by"`
	Url        string    `json:"url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FolderId   uuid.UUID `json:"folder_id"`
	Active     bool      `json:"active"`
}

type Folder struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	ParentId  uuid.UUID `json:"parent_id"`
	SubjectId uuid.UUID `json:"subject_id"`
	Active    bool      `json:"active"`
}
