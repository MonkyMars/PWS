package routes

import (
	"github.com/MonkyMars/PWS/api/internal/health"
	"github.com/MonkyMars/PWS/api/internal/auth"
	"github.com/MonkyMars/PWS/api/internal/content"
	"github.com/MonkyMars/PWS/api/internal/subjects"
	"github.com/MonkyMars/PWS/api/internal/workers"
)

type Router struct {
	AppRoutes     *health.HealthRoutes
	AuthRoutes    *auth.AuthRoutes
	ContentRoutes *content.ContentRoutes
	WorkerRoutes  *workers.WorkerRoutes
	SubjectRoutes *subjects.SubjectRoutes
}

func NewRouter() *Router {
	return &Router{
		AppRoutes:     health.NewHealthRoutes(),
		AuthRoutes:    auth.NewAuthRoutes(),
		ContentRoutes: content.NewContentRoutes(),
		WorkerRoutes:  workers.NewWorkerRoutes(),
		SubjectRoutes: subjects.NewSubjectRoutes(),
	}
}
