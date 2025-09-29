package types

import (
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
	FileID     string    `json:"file_id"`
	Name       string    `json:"name"`
	MimeType   string    `json:"mime_type"`
	SubjectID  uuid.UUID `json:"subject_id"`
	UploadedBy uuid.UUID `json:"uploaded_by"`
	Url        string    `json:"url,omitempty"`
}
