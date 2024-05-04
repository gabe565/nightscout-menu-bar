package config

import (
	"path/filepath"
	"time"
)

const LocalFileFormatCsv = "csv"

func NewDefault() *Config {
	return &Config{
		Title:    "Nightscout",
		Units:    UnitsMgdl,
		Interval: Duration{30 * time.Second},
		Arrows: Arrows{
			DoubleUp:      "⇈",
			SingleUp:      "↑",
			FortyFiveUp:   "↗",
			Flat:          "→",
			FortyFiveDown: "↘",
			SingleDown:    "↓",
			DoubleDown:    "⇊",
			Unknown:       "-",
		},
		LocalFile: LocalFile{
			Format: LocalFileFormatCsv,
			Path:   filepath.Join("$TMPDIR", "nightscout.csv"),
		},
		Advanced: Advanced{
			FetchDelay: Duration{30 * time.Second},
			RoundAge:   true,
		},
	}
}
