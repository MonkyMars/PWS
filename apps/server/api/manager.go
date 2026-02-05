package api

import (
	"github.com/MonkyMars/PWS/api/internal/auth"
	"github.com/MonkyMars/PWS/api/internal/content"
	"github.com/MonkyMars/PWS/api/internal/deadlines"
	"github.com/MonkyMars/PWS/api/internal/health"
	"github.com/MonkyMars/PWS/api/internal/subjects"
	"github.com/MonkyMars/PWS/api/internal/workers"
)

// router aggregates all route handlers for the application
// Following clean architecture principles, each route handler manages its own dependencies
type router struct {
	HealthRoutes   *health.HealthRoutes
	AuthRoutes     *auth.AuthRoutes
	ContentRoutes  *content.ContentRoutes
	WorkerRoutes   *workers.WorkerRoutes
	SubjectRoutes  *subjects.SubjectRoutes
	DeadlineRoutes *deadlines.DeadlineRoutes
}

// NewRouter creates a new Router instance with default dependencies
// This uses the WithDefaults constructors to initialize all route handlers
// with their default service implementations
func newRouter() *router {
	return &router{
		HealthRoutes:   health.NewHealthRoutesWithDefaults(),
		AuthRoutes:     auth.NewAuthRoutesWithDefaults(),
		ContentRoutes:  content.NewContentRoutesWithDefaults(),
		WorkerRoutes:   workers.NewWorkerRoutesWithDefaults(),
		SubjectRoutes:  subjects.NewSubjectRoutesWithDefaults(),
		DeadlineRoutes: deadlines.NewDeadlineRoutesWithDefaults(),
	}
}

// NewRouterWithDependencies creates a Router with explicit dependency injection
// This is useful for testing where you want to inject mock implementations
func NewRouterWithDependencies(
	healthRoutes *health.HealthRoutes,
	authRoutes *auth.AuthRoutes,
	contentRoutes *content.ContentRoutes,
	workerRoutes *workers.WorkerRoutes,
	subjectRoutes *subjects.SubjectRoutes,
	deadlineRoutes *deadlines.DeadlineRoutes,
) *router {
	return &router{
		HealthRoutes:   healthRoutes,
		AuthRoutes:     authRoutes,
		ContentRoutes:  contentRoutes,
		WorkerRoutes:   workerRoutes,
		SubjectRoutes:  subjectRoutes,
		DeadlineRoutes: deadlineRoutes,
	}
}
