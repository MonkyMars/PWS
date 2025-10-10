package content

import (
	"github.com/MonkyMars/PWS/api/middleware"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// ContentRoutes handles HTTP routing for content-related endpoints.
// It follows clean architecture principles by depending on interfaces rather than concrete implementations.
// This makes the code more testable and maintainable.
type ContentRoutes struct {
	contentService services.ContentServiceInterface
}

// NewContentRoutesWithDefaults creates a ContentRoutes instance with default dependencies.
// This is a convenience constructor for production use where you want to use
// the default implementations of all services.
func NewContentRoutesWithDefaults() *ContentRoutes {
	return &ContentRoutes{
		contentService: services.NewContentService(),
	}
}

// RegisterRoutes registers all content-related routes with the Fiber application.
// This method organizes routes logically and follows RESTful conventions.
// It groups related functionality (files, folders) and applies appropriate middleware.
func (cr *ContentRoutes) RegisterRoutes(app *fiber.App) {
	// Files API group - handles file upload, retrieval, and management
	files := app.Group("/files")
	cr.registerFileRoutes(files)

	// Folders API group - handles folder creation and hierarchy management
	folders := app.Group("/folders")
	cr.registerFolderRoutes(folders)
}

// registerFileRoutes sets up all file-related endpoints with proper middleware and handlers
func (cr *ContentRoutes) registerFileRoutes(router fiber.Router) {
	// Apply authentication middleware to all file routes
	router.Use(middleware.AuthMiddleware())

	// File upload endpoints - organized under /upload for clear separation
	upload := router.Group("/upload")
	upload.Post("/single",
		middleware.ValidateRequest[types.UploadSingleFileRequest](middleware.FileUploadValidation),
		cr.UploadSingleFile,
	)
	upload.Post("/multiple",
		middleware.ValidateRequest[types.UploadMultipleFilesRequest](middleware.FileUploadValidation),
		cr.UploadMultipleFiles,
	)

	// File retrieval endpoints - RESTful resource access patterns
	router.Get("/:fileId",
		middleware.RequiredIDParam("fileId"),
		cr.GetSingleFile,
	)
	router.Get("/subject/:subjectId/folder/:folderId",
		middleware.RequiredIDParam("subjectId"),
		middleware.RequiredIDParam("folderId"),
		cr.GetFilesBySubject,
	)

	// Future file management endpoints can be added here following RESTful conventions:
	// router.Put("/:fileId", cr.UpdateFile)          // Update file metadata
	// router.Delete("/:fileId", cr.DeleteFile)       // Delete file
	// router.Get("/", cr.ListFiles)                  // List files with query params
}

// registerFolderRoutes sets up all folder-related endpoints
func (cr *ContentRoutes) registerFolderRoutes(router fiber.Router) {
	// Folder hierarchy and retrieval - readable URL structure
	router.Get("/subject/:subjectId/folder/:parentId",
		middleware.RequiredIDParam("subjectId"),
		middleware.RequiredIDParam("parentId"),
		cr.GetFoldersBySubjectParent,
	)

	// Future folder management endpoints can be added here following RESTful conventions:
	// router.Post("/", cr.CreateFolder)              // Create new folder
	// router.Put("/:folderId", cr.UpdateFolder)      // Update folder metadata
	// router.Delete("/:folderId", cr.DeleteFolder)   // Delete folder
	// router.Get("/", cr.ListFolders)                // List folders with query params
}
