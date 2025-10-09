package routes

import (
	"github.com/MonkyMars/PWS/api/internal"
	"github.com/MonkyMars/PWS/api/internal/auth"
	"github.com/MonkyMars/PWS/api/internal/content"
)

type Router struct {
	AppRoutes     *internal.AppRoutes
	AuthRoutes    *auth.AuthRoutes
	ContentRoutes *content.ContentRoutes
	WorkerRoutes  *internal.WorkerRoutes
	SubjectRoutes *internal.SubjectRoutes
}

func NewRouter() *Router {
	return &Router{
		AppRoutes:     internal.NewAppRoutes(),
		AuthRoutes:    auth.NewAuthRoutes(),
		ContentRoutes: content.NewContentRoutes(),
		WorkerRoutes:  internal.NewWorkerRoutes(),
	}
}
