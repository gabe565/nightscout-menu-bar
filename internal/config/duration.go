package config

import (
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalText() ([]byte, error) {
	s := d.String()
	if before, found := strings.CutSuffix(s, "m0s"); found {
		s = before + "m"
	}
	return []byte(s), nil
}

func (d *Duration) UnmarshalText(text []byte) error {
	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}

	d.Duration = duration
	return nil
}
