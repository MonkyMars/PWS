package content

import (
	"fmt"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/gofiber/fiber/v3"
)

func (cr *ContentRoutes) GetFoldersBySubjectParent(c fiber.Ctx) error {
	params, err := lib.GetParams(c, "subjectId", "parentId")
	if err != nil {
		msg := "Missing required parameters in request"
		return lib.HandleServiceError(c, lib.ErrMissingField, msg)
	}

	// Retrieve folders for the subject using injected service
	folders, err := cr.contentService.GetFoldersByParentID(params["subjectId"], params["parentId"], lib.HasPrivileges(c))
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve folders for subject ID %s, parent ID %s: %v", params["subjectId"], params["parentId"], err)
		return lib.HandleServiceError(c, err, msg)
	}

	items := []any{}
	for _, folder := range folders {
		items = append(items, folder)
	}

	pageSize := lib.GetQueryParamAsInt(c, "pageSize", 20, 100)
	page := lib.GetQueryParamAsInt(c, "page", 1, 1000)

	totalPages := (len(folders) + pageSize - 1) / pageSize

	// Set pagination headers
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(items)))
	c.Set("X-Total-Pages", fmt.Sprintf("%d", totalPages))
	c.Set("X-Page-Size", fmt.Sprintf("%d", pageSize))

	// Return list of folders
	return response.Paginated(c, items, page, pageSize, len(items))
}
