package routes

import (
	internal_files "github.com/MonkyMars/PWS/api/internal/content"
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/gofiber/fiber/v3"
)

func SetupFileRoutes(app *fiber.App) {
	files := app.Group("/files", middleware.AuthMiddleware())

	files.Post("/upload/single", internal_files.UploadSingleFile)
	files.Post("/upload/multiple", internal_files.UploadMultipleFiles)
	files.Get("/:fileId", internal_files.GetSingleFile)
	files.Get("/subject/:subjectId/folder/:folderId", internal_files.GetFilesBySubject)

	folders := app.Group("/folders")
	folders.Get("/subject/:subjectId/folder/:parentId", internal_files.GetFoldersBySubjectParent)
}
