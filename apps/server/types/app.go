package types

import (
	"time"

	"github.com/google/uuid"
)

type HealthResponse struct {
	Status            string        `json:"status"`
	Message           string        `json:"message"`
	ApplicationUptime string        `json:"application_uptime"`
	DatabaseStatus    string        `json:"database_status"`
	Metrics           HealthMetrics `json:"metrics"`
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

type AuditLog struct {
	Id        uuid.UUID      `json:"id" pg:"id,pk,type:uuid,default:gen_random_uuid()"`
	Timestamp time.Time      `json:"timestamp"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Attrs     map[string]any `json:"attrs,omitempty"`
	EntryHash string         `json:"entry_hash,omitempty"`
}

type HealthLog struct {
	Timestamp      time.Time     `json:"timestamp"`
	Service        string        `json:"service"`
	StatusCode     int           `json:"status_code"`
	RequestCount   int64         `json:"request_count"`
	ErrorCount     int64         `json:"error_count"`
	AverageLatency time.Duration `json:"average_latency"`
	TimeSpan       time.Duration `json:"time_span"`
}
