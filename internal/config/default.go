package config

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
)

const LocalFileFormatCsv = "csv"

func New() *Config {
	conf := &Config{
		Title: "Nightscout",
		Units: UnitsMgdl,
		DynamicIcon: DynamicIcon{
			Enabled:   true,
			FontColor: White(),
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
			Format: LocalFileFormatCsv,
			Path:   filepath.Join("$TMPDIR", "nightscout.csv"),
		},
		Log: Log{
			Level: zerolog.InfoLevel.String(),
		},
		Advanced: Advanced{
			FetchDelay:       Duration{30 * time.Second},
			FallbackInterval: Duration{30 * time.Second},
			RoundAge:         true,
		},
	}

	switch runtime.GOOS {
	case "darwin":
		conf.DynamicIcon.Enabled = false
	case "windows":
		conf.DynamicIcon.FontColor = Black()
	}

	conf.Flags = flag.NewFlagSet("", flag.ContinueOnError)
	conf.RegisterFlags()

	return conf
}
