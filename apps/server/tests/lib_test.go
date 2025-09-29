package tests

import (
	"testing"
	"time"
	"github.com/MonkyMars/PWS/lib"
)

func TestTimeHandling(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"1 second", 1 * time.Second, "1s"},
		{"2 seconds", 2 * time.Second, "2s"},
		{"1 minute", 1 * time.Minute, "1m0s"},
		{"1 hour", 1 * time.Hour, "1h0m0s"},
		{"complex time", 1*time.Hour + 30*time.Minute + 45*time.Second, "1h30m45s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now().Add(-tt.duration)
			uptime := lib.GetUptimeString(startTime)
			if uptime != tt.expected {
				t.Errorf("Expected uptime %s, but got %s", tt.expected, uptime)
			}
		})
	}
}

func TestTimeHandlingEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"zero duration", 0 * time.Second, "0s"},
		{"less than a second", 500 * time.Millisecond, "0s"},
		{"negative duration", -1 * time.Second, "0s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now().Add(-tt.duration)
			uptime := lib.GetUptimeString(startTime)
			if uptime != tt.expected {
				t.Errorf("Expected uptime %s, but got %s", tt.expected, uptime)
			}
		})
	}
}