package config

import (
	"image/color"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gabe565.com/utils/colorx"
	flag "github.com/spf13/pflag"
)

const LocalFileFormatCsv = "csv"

func New(opts ...Option) *Config {
	conf := &Config{
		Title: "Nightscout",
		Units: UnitMgdl,
		DynamicIcon: DynamicIcon{
			Enabled:     true,
			FontColor:   colorx.Hex{Color: color.White},
			MaxFontSize: 40,
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
			Level:  strings.ToLower(slog.LevelInfo.String()),
			Format: FormatAuto.String(),
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
		conf.DynamicIcon.FontColor = colorx.Hex{Color: color.Black}
	}

	conf.Flags = flag.NewFlagSet("", flag.ContinueOnError)
	conf.RegisterFlags()

	for _, opt := range opts {
		opt(conf)
	}

	return conf
}
