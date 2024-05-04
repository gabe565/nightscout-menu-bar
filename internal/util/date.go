package util

import (
	"strings"
	"time"
)

func MinAgo(t time.Time, round bool) string {
	// Drop resolution to minutes
	duration := time.Since(t)
	if round {
		duration = duration.Round(time.Minute)
	} else {
		duration = duration.Truncate(time.Minute)
	}
	str := duration.String()
	str = strings.TrimSuffix(str, "0s")
	if str == "" {
		str = "0m"
	}
	return str
}

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
