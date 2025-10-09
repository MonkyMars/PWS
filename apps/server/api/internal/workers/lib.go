package workers

// extractWorkerMetrics extracts relevant metrics from worker health status
func (wr *WorkerRoutes) ExtractWorkerMetrics(workerHealth any) map[string]any {
	if workerHealth == nil {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
		}
	}

	health, ok := workerHealth.(map[string]any)
	if !ok {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "invalid health data format",
		}
	}

	metrics := map[string]any{
		"enabled":        false,
		"worker_running": false,
		"is_healthy":     false,
	}

	if health != nil {
		if enabled, ok := health["enabled"]; ok {
			metrics["enabled"] = enabled
		}
		if running, ok := health["worker_running"]; ok {
			metrics["worker_running"] = running
		}
		if healthy, ok := health["is_healthy"]; ok {
			metrics["is_healthy"] = healthy
		}
	}

	// Add optional metrics if they exist
	if queueSize, ok := health["queue_size"]; ok {
		metrics["queue_size"] = queueSize
	}
	if queueCapacity, ok := health["queue_capacity"]; ok {
		metrics["queue_capacity"] = queueCapacity
	}
	if lastFlush, ok := health["last_flush_time"]; ok {
		metrics["last_flush_time"] = lastFlush
	}

	return metrics
}
