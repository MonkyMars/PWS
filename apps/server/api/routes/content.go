package routes

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) SetupContentRoutes(app *fiber.App) {
	files := app.Group("/files", middleware.AuthMiddleware())

	files.Post("/upload/single", r.ContentRoutes.UploadSingleFile)
	files.Post("/upload/multiple", r.ContentRoutes.UploadMultipleFiles)
	files.Get("/:fileId", r.ContentRoutes.GetSingleFile)
	files.Get("/subject/:subjectId/folder/:folderId", r.ContentRoutes.GetFilesBySubject)

	folders := app.Group("/folders")
	folders.Get("/subject/:subjectId/folder/:parentId", r.ContentRoutes.GetFoldersBySubjectParent)
}
