package content

import (
	"fmt"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

func (cr *ContentRoutes) GetSingleFile(c fiber.Ctx) error {
	// Get fileId from URL parameters
	params, err := lib.GetParams(c, "fileId")
	if err != nil {
		return lib.HandleServiceError(c, err)
	}

	// Retrieve file metadata using injected service
	file, err := cr.contentService.GetFileByID(params["fileId"])
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve file metadata for file ID %s: %v", fileID, err)
		return lib.HandleServiceError(c, err, msg)
	}
	if file == nil {
		msg := fmt.Sprintf("File not found for file ID %s", fileID)
		return lib.HandleServiceError(c, lib.ErrFileNotFound, msg)
	}

	// Return file metadata
	return response.Success(c, file)
}

func (cr *ContentRoutes) GetFilesBySubject(c fiber.Ctx) error {
	// Get parameters from URL
	params, err := lib.GetParams(c, "subjectId", "folderId")
	if err != nil {
		return lib.HandleServiceError(c, err)
	}

	// Retrieve files for the subject using injected service
	files, err := cr.contentService.GetFilesBySubjectID(params["subjectId"], params["folderId"], lib.HasPrivileges(c))
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve files for subject ID %s, folder ID %s: %v", subjectId, folderId, err)
		return lib.HandleServiceError(c, err, msg)
	}

	items := []any{}
	for _, file := range files {
		items = append(items, file)
	}

	page := lib.GetQueryParamAsInt(c, "page", 1, 1000)
	pageSize := lib.GetQueryParamAsInt(c, "pageSize", 20, 100)

	totalPages := (len(files) + pageSize - 1) / pageSize

	// Set pagination headers
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(items)))
	c.Set("X-Total-Pages", fmt.Sprintf("%d", totalPages))
	c.Set("X-Page-Size", fmt.Sprintf("%d", pageSize))

	// Return list of files
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(items) {
		start = len(items)
	}
	if end > len(items) {
		end = len(items)
	}
	paginatedItems := items[start:end]
	return response.Paginated(c, paginatedItems, page, pageSize, len(items))
}

// /files/upload/single
func (cr *ContentRoutes) UploadSingleFile(c fiber.Ctx) error {
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		msg := "Failed to get authenticated user claims for file upload"
		return lib.HandleServiceError(c, err, msg)
	}

	if claims.Role != lib.RoleAdmin && claims.Role != lib.RoleTeacher {
		msg := fmt.Sprintf("User ID %s with role %s attempted to upload file without permission", claims.Sub, claims.Role)
		return lib.HandleServiceError(c, lib.ErrForbidden, msg)
	}

	// Parse request body
	var req types.UploadSingleFileRequest
	if err := c.Bind().Body(&req); err != nil {
		msg := fmt.Sprintf("Failed to parse single file upload request body for user ID %s: %v", claims.Sub, err)
		return lib.HandleServiceError(c, lib.ErrInvalidRequest, msg)
	}

	if req.File.FileID == "" || req.File.Name == "" || req.File.MimeType == "" {
		msg := fmt.Sprintf("Missing required file fields in upload request for user ID %s", claims.Sub)
		return lib.HandleServiceError(c, lib.ErrMissingField, msg)
	}

	// Upload meta data using injected service
	fileData := map[string]any{
		"file_id":     req.File.FileID,
		"name":        req.File.Name,
		"mime_type":   req.File.MimeType,
		"subject_id":  req.SubjectID,
		"uploaded_by": user.Id,
		"url":         fmt.Sprintf("https://drive.google.com/file/d/%s/preview", req.File.FileID),
	}

	file, err := cr.contentService.CreateFile(fileData)
	if err != nil {
		msg := fmt.Sprintf("Failed to upload file %s (ID: %s) for user ID %s: %v", req.File.Name, req.File.FileID, claims.Sub, err)
		return lib.HandleServiceError(c, err, msg)
	}

	return response.Created(c, file)
}

func (cr *ContentRoutes) UploadMultipleFiles(c fiber.Ctx) error {
	claims, err := lib.GetValidatedClaims(c)
	if err != nil {
		msg := "Failed to get authenticated user claims for multiple file upload"
		return lib.HandleServiceError(c, err, msg)
	}

	if claims.Role != lib.RoleAdmin && claims.Role != lib.RoleTeacher {
		msg := fmt.Sprintf("User ID %s with role %s attempted to upload multiple files without permission", claims.Sub, claims.Role)
		return lib.HandleServiceError(c, lib.ErrForbidden, msg)
	}

	// Parse request body
	var req types.UploadMultipleFilesRequest
	if err := c.Bind().Body(&req); err != nil {
		msg := fmt.Sprintf("Failed to parse multiple files upload request body for user ID %s: %v", claims.Sub, err)
		return lib.HandleServiceError(c, lib.ErrInvalidRequest, msg)
	}

	if len(req.Files) == 0 {
		msg := fmt.Sprintf("No files provided in multiple upload request for user ID %s", claims.Sub)
		return lib.HandleServiceError(c, lib.ErrMissingField, msg)
	}

	// Validate and prepare file data
	filesData := make([]map[string]any, 0, len(req.Files))
	for _, file := range req.Files {
		if file.FileID == "" || file.Name == "" || file.MimeType == "" {
			msg := fmt.Sprintf("Missing required file fields for file %s in multiple upload for user ID %s", file.Name, claims.Sub)
			return lib.HandleServiceError(c, lib.ErrMissingField, msg)
		}

		fileData := map[string]any{
			"file_id":     file.FileID,
			"name":        file.Name,
			"mime_type":   file.MimeType,
			"subject_id":  req.SubjectID,
			"uploaded_by": user.Id,
			"url":         fmt.Sprintf("https://drive.google.com/file/d/%s/preview", file.FileID),
		}
		filesData = append(filesData, fileData)
	}

	// Upload metadata using injected service
	files, err := cr.contentService.CreateMultipleFiles(filesData)
	if err != nil {
		msg := fmt.Sprintf("Failed to upload %d files for user ID %s, subject ID %s: %v", len(req.Files), claims.Sub, req.SubjectID, err)
		return lib.HandleServiceError(c, err, msg)
	}

	return response.Created(c, files)
}
