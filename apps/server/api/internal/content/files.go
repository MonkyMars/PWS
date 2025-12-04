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
		return lib.HandleServiceError(c, err)
	}
	if file == nil {
		return lib.HandleServiceError(c, lib.ErrFileNotFound)
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
		return lib.HandleServiceError(c, err)
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
	return response.Paginated(c, items, page, pageSize, len(items))
}

// /files/upload/single
func (cr *ContentRoutes) UploadSingleFile(c fiber.Ctx) error {
	user := lib.GetUserFromContext(c)
	if user == nil {
		cr.logger.AuditWarn("UploadSingleFile: Unauthorized access attempt")
		return response.Unauthorized(c, "Unauthorized")
	}

	if !lib.HasPrivileges(c) {
		return response.Forbidden(c, "You do not have permission to upload files")
	}

	// Parse request body
	var req types.UploadSingleFileRequest
	if err := c.Bind().Body(&req); err != nil {
		cr.logger.AuditWarn("UploadSingleFile: Failed to parse request body - %v", err)
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	if req.File.FileID == "" || req.File.Name == "" || req.File.MimeType == "" {
		cr.logger.AuditWarn("UploadSingleFile: Missing required file fields")
		return response.BadRequest(c, "Missing required file fields")
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
		return lib.HandleServiceError(c, err)
	}

	return response.Created(c, file)
}

func (cr *ContentRoutes) UploadMultipleFiles(c fiber.Ctx) error {
	if !lib.HasPrivileges(c) {
		return lib.HandleServiceError(c, lib.ErrInsufficientPermissions)
	}

	user := lib.GetUserFromContext(c)
	if user == nil {
		return lib.HandleServiceError(c, lib.ErrUnauthorized)
	}

	// Parse request body
	var req types.UploadMultipleFilesRequest
	if err := c.Bind().Body(&req); err != nil {
		cr.logger.AuditError("UploadMultipleFiles: Failed to parse request body - %v", err)
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	if len(req.Files) == 0 {
		return lib.HandleServiceError(c, lib.ErrInvalidInput)
	}

	// Validate and prepare file data
	filesData := make([]map[string]any, 0, len(req.Files))
	for _, file := range req.Files {
		if file.FileID == "" || file.Name == "" || file.MimeType == "" {
			cr.logger.AuditError("UploadMultipleFiles: Missing required file fields for file: %s", file.Name)
			return response.BadRequest(c, "Missing required file fields")
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
		return lib.HandleServiceError(c, err)
	}

	return response.Created(c, files)
}
