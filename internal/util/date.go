package util

import (
	"time"
)

func GetNextMinChange(t time.Time, round bool) time.Duration {
	if round {
		// Offset time by 30s since output is rounded
		t = t.Add(-30 * time.Second)
	}
	// Time since last update
	duration := time.Since(t)
	// Only keep seconds
	duration %= time.Minute
	// Time until rounded output changes
	duration = time.Minute - duration
	return duration
}
