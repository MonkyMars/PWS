package content

import (
	"fmt"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/services"
	"github.com/gofiber/fiber/v3"
)

// /files/:fileId
func (cr *ContentRoutes) GetSingleFile(c fiber.Ctx) error {
	// Get fileId from URL parameters
	fileID := c.Params("fileId")
	if fileID == "" {
		return response.BadRequest(c, "fileId parameter is required")
	}

	// Retrieve file metadata from database
	contentService := services.NewContentService()
	file, err := contentService.GetFileByID(fileID)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve file: "+err.Error())
	}
	if file == nil {
		return response.NotFound(c, "File not found")
	}

	// Return file metadata
	return response.Success(c, file)
}

// /files/subject/:subjectId/folder/:folderId
func (cr *ContentRoutes) GetFilesBySubject(c fiber.Ctx) error {
	// Get subjectId from URL parameters
	subjectId := c.Params("subjectId")
	if subjectId == "" {
		return response.BadRequest(c, "subjectId parameter is required")
	}

	folderId := c.Params("folderId")
	if folderId == "" {
		return response.BadRequest(c, "folderId parameter is required")
	}

	// Retrieve files for the subject from database
	contentService := &services.ContentService{}
	files, err := contentService.GetFilesBySubjectID(subjectId, folderId)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve files: "+err.Error())
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

func (cr *ContentRoutes) GetFoldersBySubjectParent(c fiber.Ctx) error {
	subjectId := c.Params("subjectId")
	if subjectId == "" {
		return response.BadRequest(c, "subjectId parameter is required")
	}

	parentId := c.Params("parentId")
	if parentId == "" {
		return response.BadRequest(c, "parentId parameter is required")
	}

	// Retrieve folders for the subject from database
	contentService := &services.ContentService{}
	folders, err := contentService.GetFoldersByParentID(subjectId, parentId)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve folders: "+err.Error())
	}

	items := []any{}
	for _, folder := range folders {
		items = append(items, folder)
	}

	// Set pagination headers
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(items)))
	c.Set("X-Total-Pages", fmt.Sprintf("%s", "1"))
	c.Set("X-Page-Size", fmt.Sprintf("%d", len(items)))

	// Return list of folders
	return response.Paginated(c, items, len(folders), 1, len(folders))
}
