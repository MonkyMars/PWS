package lib

import "time"

func GetUptimeString(startTime time.Time) string {
	duration := time.Since(startTime)
	return duration.Truncate(time.Second).String()
}
