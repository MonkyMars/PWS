package health

import (
	"runtime"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

func (hr *HealthRoutes) GetSystemHealth(c fiber.Ctx) error {
	// Memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Check database health
	dbStatus := "ok"
	if err := services.Ping(); err != nil {
		dbStatus = "error"
	}

	status := "ok"
	message := "All systems operational"
	if dbStatus != "ok" {
		status = "degraded"
		message = "Database connection issues"
	}

	return response.Success(c, types.HealthResponse{
		Status:            status,
		Message:           message,
		ApplicationUptime: lib.GetUptimeString(appStartTime),
		DatabaseStatus:    dbStatus,
		Metrics: types.HealthMetrics{
			MemoryUsageMB: float64(memStats.Alloc) / 1024 / 1024,
			GoRoutines:    runtime.NumGoroutine(),
			RequestCount:  requestCount,
		},
	})
}
