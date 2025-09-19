package internal

import (
	"runtime"
	"time"

	"github.com/MonkyMars/PWS/api/response"
	"github.com/MonkyMars/PWS/lib"
	"github.com/MonkyMars/PWS/services"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

var (
	appStartTime = time.Now()
	requestCount int64
)

func GetSystemHealth(c fiber.Ctx) error {
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
		Services: map[string]string{
			"database": dbStatus,
			"api":      "ok",
		},
		Metrics: types.HealthMetrics{
			MemoryUsageMB: float64(memStats.Alloc) / 1024 / 1024,
			GoRoutines:    runtime.NumGoroutine(),
			RequestCount:  requestCount,
		},
	})
}

func GetDatabaseHealth(c fiber.Ctx) error {
	now := time.Now()
	if err := services.Ping(); err != nil {
		return response.ServiceUnavailable(c, "Database connection error: "+err.Error())
	}

	return response.Success(c, types.DatabaseHealthResponse{
		Status:  "ok",
		Message: "Database connection is healthy",
		Elapsed: time.Since(now).String(),
	})
}

func NotFoundHandler(c fiber.Ctx) error {
	return response.NotFound(c, "The requested resource was not found.")
}
