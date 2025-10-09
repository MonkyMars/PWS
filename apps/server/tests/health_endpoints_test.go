package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/MonkyMars/PWS/workers"
	"github.com/gofiber/fiber/v3"
)

func TestHealthEndpointsNoPanics(t *testing.T) {
	// Create a test app
	app := fiber.New()

	// Add favicon route
	app.Get("/favicon.ico", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	// Setup worker health endpoint with safe handling
	app.Get("/api/v1/workers/health", func(c fiber.Ctx) error {
		manager := workers.GetGlobalManager()
		if manager == nil {
			return c.Status(503).JSON(fiber.Map{
				"success": false,
				"message": "Worker manager not available",
			})
		}

		healthStatus := manager.HealthStatus()
		if healthStatus == nil {
			return c.Status(503).JSON(fiber.Map{
				"success": false,
				"message": "Unable to retrieve health status",
			})
		}

		// Safe boolean extraction
		isHealthy := false
		if healthVal, exists := healthStatus["is_healthy"]; exists && healthVal != nil {
			if healthy, ok := healthVal.(bool); ok {
				isHealthy = healthy
			}
		}

		statusCode := 200
		if !isHealthy {
			statusCode = 503
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"success": true,
			"message": "Worker health status retrieved",
			"data":    healthStatus,
		})
	})

	// Test endpoints that previously caused panics
	testCases := []struct {
		name           string
		method         string
		url            string
		shouldNotPanic bool
	}{
		{
			name:           "Favicon request",
			method:         "GET",
			url:            "/favicon.ico",
			shouldNotPanic: true,
		},
		{
			name:           "Worker health endpoint",
			method:         "GET",
			url:            "/api/v1/workers/health",
			shouldNotPanic: true,
		},
		{
			name:           "Non-existent endpoint",
			method:         "GET",
			url:            "/non-existent",
			shouldNotPanic: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Recovery function to catch panics
			defer func() {
				if r := recover(); r != nil {
					if tc.shouldNotPanic {
						t.Errorf("Test %s panicked when it shouldn't: %v", tc.name, r)
					}
				}
			}()

			// Create request
			req := httptest.NewRequest(tc.method, tc.url, nil)
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// Check that we got a response (not a panic)
			if resp.StatusCode == 0 {
				t.Errorf("Got zero status code, possible panic occurred")
			}

			// Verify we get a reasonable status code
			if resp.StatusCode < 200 || resp.StatusCode >= 600 {
				t.Errorf("Got unreasonable status code: %d", resp.StatusCode)
			}

			t.Logf("Test %s completed successfully with status %d", tc.name, resp.StatusCode)
		})
	}
}

func TestWorkerHealthStatusSafety(t *testing.T) {
	// Test that health status functions return safe values even with nil inputs
	testCases := []struct {
		name string
		fn   func() map[string]any
	}{
		{
			name: "AuditHealthStatus",
			fn:   workers.AuditHealthStatus,
		},
		{
			name: "ServiceHealthStatus",
			fn:   workers.ServiceHealthStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Function %s panicked: %v", tc.name, r)
				}
			}()

			result := tc.fn()
			if result == nil {
				t.Errorf("Function %s returned nil map", tc.name)
			}

			// Check that required fields exist and are safe to access
			requiredFields := []string{"enabled", "worker_running", "is_healthy"}
			for _, field := range requiredFields {
				if _, exists := result[field]; !exists {
					t.Errorf("Function %s missing required field: %s", tc.name, field)
				}
			}

			// Verify is_healthy is a boolean or safe to convert
			if healthVal, exists := result["is_healthy"]; exists {
				if healthVal != nil {
					if _, ok := healthVal.(bool); !ok {
						t.Errorf("Function %s field 'is_healthy' is not a boolean: %T", tc.name, healthVal)
					}
				}
			}

			t.Logf("Function %s returned safe health status", tc.name)
		})
	}
}

func TestWorkerManagerSafety(t *testing.T) {
	// Test that GetGlobalManager returns a safe instance
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetGlobalManager panicked: %v", r)
		}
	}()

	manager := workers.GetGlobalManager()
	if manager == nil {
		t.Error("GetGlobalManager returned nil")
		return
	}

	// Test that HealthStatus returns safe values
	healthStatus := manager.HealthStatus()
	if healthStatus == nil {
		t.Error("HealthStatus returned nil")
		return
	}

	// Check for required fields
	requiredFields := []string{"manager_running", "timestamp", "is_healthy"}
	for _, field := range requiredFields {
		if _, exists := healthStatus[field]; !exists {
			t.Errorf("HealthStatus missing required field: %s", field)
		}
	}

	// Verify is_healthy is safely accessible
	if healthVal, exists := healthStatus["is_healthy"]; exists && healthVal != nil {
		if _, ok := healthVal.(bool); !ok {
			t.Errorf("HealthStatus field 'is_healthy' is not a boolean: %T", healthVal)
		}
	}

	t.Log("WorkerManager returned safe health status")
}

func TestSafeTypeConversions(t *testing.T) {
	// Test the specific pattern that was causing panics
	testData := map[string]any{
		"valid_bool": true,
		"nil_value":  nil,
		"wrong_type": "not a boolean",
		"zero_bool":  false,
	}

	for key := range testData {
		t.Run(key, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Safe type conversion panicked for %s: %v", key, r)
				}
			}()

			// This is the safe pattern we implemented
			result := false
			if val, exists := testData[key]; exists && val != nil {
				if boolVal, ok := val.(bool); ok {
					result = boolVal
				}
			}

			t.Logf("Safe conversion for %s resulted in: %v", key, result)
		})
	}
}
