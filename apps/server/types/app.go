package types

type HealthResponse struct {
	Status            string            `json:"status"`
	Message           string            `json:"message"`
	ApplicationUptime string            `json:"application_uptime"`
	Services          map[string]string `json:"services"`
	Metrics           HealthMetrics     `json:"metrics"`
}

type HealthMetrics struct {
	MemoryUsageMB float64 `json:"memory_usage_mb"`
	GoRoutines    int     `json:"go_routines"`
	RequestCount  int64   `json:"request_count,omitempty"`
}

type DatabaseHealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Elapsed string `json:"elapsed,omitempty"`
}
