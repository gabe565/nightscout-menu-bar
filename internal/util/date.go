package util

import "time"

func GetNextMinChange(t time.Time) time.Duration {
	// Time since last update
	duration := time.Since(t)
	// Only keep seconds
	duration %= time.Minute
	// Time until rounded output changes
	duration = time.Minute - duration
	return duration
}
