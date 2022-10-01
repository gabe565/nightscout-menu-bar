package util

import (
	"strings"
	"time"
)

func MinAgo(t time.Time) string {
	// Drop resolution to minutes
	duration := time.Since(t).Truncate(time.Minute)
	str := duration.String()
	str = strings.TrimSuffix(str, "0s")
	if str == "" {
		str = "0m"
	}
	return str
}

func GetNextMinChange(t time.Time) time.Duration {
	duration := time.Since(t)
	return time.Minute - duration%time.Minute
}
