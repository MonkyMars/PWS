package content

import (
	"fmt"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/gofiber/fiber/v3"
)

func (cr *ContentRoutes) GetFoldersBySubjectParent(c fiber.Ctx) error {
	subjectId := c.Params("subjectId")
	if subjectId == "" {
		return response.BadRequest(c, "subjectId parameter is required")
	}

	parentId := c.Params("parentId")
	if parentId == "" {
		return response.BadRequest(c, "parentId parameter is required")
	}

	// Retrieve folders for the subject using injected service
	folders, err := cr.contentService.GetFoldersByParentID(subjectId, parentId)
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
