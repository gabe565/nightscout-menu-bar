package config

import (
	"path/filepath"
	"time"
)

var Default = NewDefault()

const LocalFileFormatCsv = "csv"

func NewDefault() Config {
	return Config{
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
			Format:  LocalFileFormatCsv,
			Path:    filepath.Join("$TMPDIR", "nightscout.csv"),
			Cleanup: true,
		},
	}
}
