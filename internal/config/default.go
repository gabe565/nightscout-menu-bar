package config

import (
	"path/filepath"
	"runtime"
	"time"
)

var Default = NewDefault()

const LocalFileFormatCsv = "csv"

func NewDefault() Config {
	dynamicIconEnabled := true
	dynamicIconColor := White
	switch runtime.GOOS {
	case "darwin":
		dynamicIconEnabled = false
	case "windows":
		dynamicIconColor = Black
	}

	return Config{
		Title:    "Nightscout",
		Units:    UnitsMgdl,
		Interval: Duration{30 * time.Second},
		DynamicIcon: DynamicIcon{
			Enabled:   dynamicIconEnabled,
			FontColor: dynamicIconColor,
			FontSize:  19,
		},
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
