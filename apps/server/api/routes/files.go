package routes

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/gofiber/fiber/v3"
)

type FileRoutes struct{}

func NewFileRoutes() *FileRoutes {
	return &FileRoutes{}
}

func (r *Router) SetupFileRoutes(app *fiber.App) {
	files := app.Group("/files", middleware.AuthMiddleware())

	files.Post("/upload/single", r.FilesRoutes.UploadSingleFile)
	files.Post("/upload/multiple", r.FilesRoutes.UploadMultipleFiles)
	files.Get("/:fileId", r.FilesRoutes.GetSingleFile)
	files.Get("/subject/:subjectId", r.FilesRoutes.GetFilesBySubject)
}
