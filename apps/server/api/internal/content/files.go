package content

import (
	"fmt"
	"log"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

func (cr *ContentRoutes) GetSingleFile(c fiber.Ctx) error {
	// Get fileId from URL parameters (already validated by middleware)
	fileID := c.Params("fileId")

	// Retrieve file metadata using injected service
	file, err := cr.contentService.GetFileByID(fileID)
	if err != nil {
		return lib.HandleServiceError(c, err)
	}
	if file == nil {
		return lib.HandleServiceError(c, lib.ErrFileNotFound)
	}

	// Return file metadata
	return response.Success(c, file)
}

func (cr *ContentRoutes) GetFilesBySubject(c fiber.Ctx) error {
	// Get parameters from URL (already validated by middleware)
	subjectId := c.Params("subjectId")
	folderId := c.Params("folderId")

	// Retrieve files for the subject using injected service
	files, err := cr.contentService.GetFilesBySubjectID(subjectId, folderId)
	if err != nil {
		return lib.HandleServiceError(c, err)
	}

	items := []any{}
	for _, file := range files {
		items = append(items, file)
	}

	// Set pagination headers
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(items)))
	c.Set("X-Total-Pages", fmt.Sprintf("%s", "1"))
	c.Set("X-Page-Size", fmt.Sprintf("%d", len(items)))

	// Return list of files
	return response.Paginated(c, items, len(files), 1, len(files))
}

// /files/upload/single
func (cr *ContentRoutes) UploadSingleFile(c fiber.Ctx) error {
	claimsInterface := c.Locals("claims")

	if claimsInterface == nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	// Type assert claims
	claims, ok := claimsInterface.(*types.AuthClaims)
	if claims == nil || !ok {
		return response.Unauthorized(c, "Unauthorized")
	}

	if claims.Role != lib.RoleAdmin && claims.Role != lib.RoleTeacher {
		return response.Forbidden(c, "You do not have permission to upload files")
	}

	// Parse request body
	var req types.UploadSingleFileRequest
	if err := c.Bind().Body(&req); err != nil {
		log.Printf("UploadSingleFile: Failed to parse request body - %v", err)
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	if req.File.FileID == "" || req.File.Name == "" || req.File.MimeType == "" {
		log.Printf("UploadSingleFile: Missing required file fields")
		return response.BadRequest(c, "Missing required file fields")
	}

	// Upload meta data using injected service
	fileData := map[string]any{
		"file_id":     req.File.FileID,
		"name":        req.File.Name,
		"mime_type":   req.File.MimeType,
		"subject_id":  req.SubjectID,
		"uploaded_by": claims.Sub,
		"url":         fmt.Sprintf("https://drive.google.com/file/d/%s/preview", req.File.FileID),
	}

	file, err := cr.contentService.CreateFile(fileData)
	if err != nil {
		log.Printf("UploadSingleFile: Service failed - %v", err)
		return response.InternalServerError(c, "Failed to upload file: "+err.Error())
	}

	return response.Created(c, file)
}

func (cr *ContentRoutes) UploadMultipleFiles(c fiber.Ctx) error {
	logger := config.SetupLogger()
	claimsInterface := c.Locals("claims")

	if claimsInterface == nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	// Type assert claims
	claims, ok := claimsInterface.(*types.AuthClaims)
	if claims == nil || !ok {
		return response.Unauthorized(c, "Unauthorized")
	}

	if claims.Role != lib.RoleAdmin && claims.Role != lib.RoleTeacher {
		logger.AuditError("UploadMultipleFiles: User does not have permission to upload files")
		return response.Forbidden(c, "You do not have permission to upload files")
	}

	// Parse request body
	var req types.UploadMultipleFilesRequest
	if err := c.Bind().Body(&req); err != nil {
		logger.AuditError("UploadMultipleFiles: Failed to parse request body - %v", err)
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	if len(req.Files) == 0 {
		return response.BadRequest(c, "No files to upload")
	}

	// Validate and prepare file data
	filesData := make([]map[string]any, 0, len(req.Files))
	for _, file := range req.Files {
		if file.FileID == "" || file.Name == "" || file.MimeType == "" {
			logger.AuditError("UploadMultipleFiles: Missing required file fields for file: %s", file.Name)
			return response.BadRequest(c, "Missing required file fields")
		}

		fileData := map[string]any{
			"file_id":     file.FileID,
			"name":        file.Name,
			"mime_type":   file.MimeType,
			"subject_id":  req.SubjectID,
			"uploaded_by": claims.Sub,
			"url":         fmt.Sprintf("https://drive.google.com/file/d/%s/preview", file.FileID),
		}
		filesData = append(filesData, fileData)
	}

	// Upload metadata using injected service
	files, err := cr.contentService.CreateMultipleFiles(filesData)
	if err != nil {
		log.Printf("UploadMultipleFiles: Service failed - %v", err)
		return response.InternalServerError(c, "Failed to upload files: "+err.Error())
	}

	return response.Created(c, files)
}
