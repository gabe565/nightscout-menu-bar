package config

import (
	"image/color"
	"path/filepath"
	"runtime"
	"time"

	"gabe565.com/utils/colorx"
	"gabe565.com/utils/slogx"
	flag "github.com/spf13/pflag"
)

const SocketFormatCSV = "csv"

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
		Socket: Socket{
			Format: SocketFormatCSV,
			Path:   filepath.Join("$TMPDIR", "nightscout.sock"),
		},
		Log: Log{
			Level:  slogx.LevelInfo,
			Format: slogx.FormatAuto,
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
