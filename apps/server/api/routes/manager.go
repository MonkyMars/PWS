package routes

import (
	"github.com/MonkyMars/PWS/api/internal"
	"github.com/MonkyMars/PWS/api/internal/auth"
	"github.com/MonkyMars/PWS/api/internal/files"
)

type Router struct {
	AppRoutes    *internal.AppRoutes
	AuthRoutes   *auth.AuthRoutes
	FilesRoutes  *files.FileRoutes
	WorkerRoutes *internal.WorkerRoutes
}

func NewRouter() *Router {
	return &Router{
		AppRoutes:    internal.NewAppRoutes(),
		AuthRoutes:   auth.NewAuthRoutes(),
		FilesRoutes:  files.NewFileRoutes(),
		WorkerRoutes: internal.NewWorkerRoutes(),
	}
}
