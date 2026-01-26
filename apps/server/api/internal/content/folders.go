package content

import (
	"fmt"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/gofiber/fiber/v3"
)

func (cr *ContentRoutes) GetFoldersBySubjectParent(c fiber.Ctx) error {
	subjectId := c.Params("subjectId")
	if subjectId == "" {
		msg := "Missing required subjectId parameter in request"
		return lib.HandleServiceError(c, lib.ErrMissingField, msg)
	}

	parentId := c.Params("parentId")
	if parentId == "" {
		msg := "Missing required parentId parameter in request"
		return lib.HandleServiceError(c, lib.ErrMissingField, msg)
	}

	// Retrieve folders for the subject using injected service
	folders, err := cr.contentService.GetFoldersByParentID(params["subjectId"], params["parentId"], lib.HasPrivileges(c))
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve folders for subject ID %s, parent ID %s: %v", subjectId, parentId, err)
		return lib.HandleServiceError(c, err, msg)
	}

	items := []any{}
	for _, folder := range folders {
		items = append(items, folder)
	}

	totalPages := (len(folders) + pageSize - 1) / pageSize

	// Set pagination headers
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(items)))
	c.Set("X-Total-Pages", fmt.Sprintf("%d", totalPages))
	c.Set("X-Page-Size", fmt.Sprintf("%d", pageSize))

	// Return list of folders
	return response.Paginated(c, items, page, pageSize, len(items))
}
