package health

import "time"

type HealthRoutes struct{}

func NewHealthRoutes() *HealthRoutes {
	return &HealthRoutes{}
}

var (
	appStartTime = time.Now()
	requestCount int64
)
