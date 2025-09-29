package files

import (
	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
)

// /files/:fileId
func GetSingleFile(c fiber.Ctx) error {
	// Get fileId from URL parameters
	fileID := c.Params("fileId")
	if fileID == "" {
		return response.BadRequest(c, "fileId parameter is required")
	}

	// Retrieve file metadata from database
	fileService := &services.FileService{}
	file, err := fileService.GetFileByID(fileID)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve file: "+err.Error())
	}
	if file == nil {
		return response.NotFound(c, "File not found")
	}

	// Return file metadata
	return response.Success(c, file)
}

// /files/subject/:subjectId
func GetFilesBySubject(c fiber.Ctx) error {
	// Get subjectId from URL parameters
	subjectID := c.Params("subjectId")
	if subjectID == "" {
		return response.BadRequest(c, "subjectId parameter is required")
	}

	// Retrieve files for the subject from database
	fileService := &services.FileService{}
	files, err := fileService.GetFilesBySubjectID(subjectID)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve files: "+err.Error())
	}

	// Return list of files
	return response.Success(c, files)
}
