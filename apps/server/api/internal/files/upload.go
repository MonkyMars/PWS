package files

import (
	"fmt"
	"log"
	"sync"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/database"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// /files/upload/single
func UploadSingleFile(c fiber.Ctx) error {
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

	// Upload meta data to database
	fileData := map[string]any{
		"file_id":     req.File.FileID,
		"name":        req.File.Name,
		"mime_type":   req.File.MimeType,
		"subject_id":  req.SubjectID,
		"uploaded_by": claims.Sub,
		"url":         fmt.Sprintf("https://drive.google.com/file/d/%s/preview", req.File.FileID),
	}

	query := services.Query().SetOperation("insert").SetTable("files").SetData(fileData)
	query.Returning = []string{"file_id", "name", "mime_type", "subject_id", "uploaded_by", "url"}

	data, err := database.ExecuteQuery[types.File](query)
	if err != nil {
		log.Printf("UploadSingleFile: Database query failed - %v", err)
		return response.InternalServerError(c, "Failed to upload file: "+err.Error())
	}

	return response.Created(c, data.Single)
}

func UploadMultipleFiles(c fiber.Ctx) error {
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
	var req types.UploadMultipleFilesRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, "Invalid request body: "+err.Error())
	}

	if len(req.Files) == 0 {
		return response.BadRequest(c, "No files to upload")
	}

	// Validate all files upfront before starting goroutines
	for _, file := range req.Files {
		if file.FileID == "" || file.Name == "" || file.MimeType == "" {
			return response.BadRequest(c, "Missing required file fields")
		}
	}

	// Result structure for collecting goroutine results
	type uploadResult struct {
		file  *types.File
		err   error
		index int
	}

	// Create channels and sync primitives
	results := make([]uploadResult, len(req.Files))
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Start goroutines for concurrent uploads
	for i, file := range req.Files {
		wg.Add(1)
		go func(file types.DriveFile, index int) {
			defer wg.Done()

			fileData := map[string]any{
				"file_id":     file.FileID,
				"name":        file.Name,
				"mime_type":   file.MimeType,
				"subject_id":  req.SubjectID,
				"uploaded_by": claims.Sub,
				"url":         fmt.Sprintf("https://drive.google.com/file/d/%s/preview", file.FileID),
			}

			query := services.Query().SetOperation("insert").SetTable("files").SetData(fileData)
			query.Returning = []string{"file_id", "name", "mime_type", "subject_id", "uploaded_by", "url"}
			result, err := database.ExecuteQuery[types.File](query)

			// Safely store the result
			mu.Lock()
			if err != nil {
				results[index] = uploadResult{nil, fmt.Errorf("failed to upload %s: %w", file.Name, err), index}
			} else {
				results[index] = uploadResult{result.Single, nil, index}
			}
			mu.Unlock()
		}(file, i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check for any errors and collect successful uploads
	var uploadedFiles []*types.File
	var firstError error

	for _, result := range results {
		if result.err != nil && firstError == nil {
			firstError = result.err
		} else if result.err == nil {
			uploadedFiles = append(uploadedFiles, result.file)
		}
	}

	// If any errors occurred, return the first one
	if firstError != nil {
		// TODO: Rollback logic can be added here to prevent dead files in DB
		return response.InternalServerError(c, firstError.Error())
	}

	return response.Created(c, uploadedFiles)
}
