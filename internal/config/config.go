package config

import (
	"os"
	"path/filepath"
	"runtime"

	"gabe565.com/utils/colorx"
	"gabe565.com/utils/slogx"
	"github.com/spf13/pflag"
)

type Config struct {
	File      string         `toml:"-"`
	Flags     *pflag.FlagSet `toml:"-"`
	Version   string         `toml:"-"`
	callbacks []func()       `toml:"-"`

	Title       string      `toml:"title" comment:"Tray title."`
	URL         string      `toml:"url" comment:"Nightscout URL. (required)"`
	Token       string      `toml:"token" comment:"Nightscout token. Using an access token is recommended instead of the API secret."`
	Units       Unit        `toml:"units" comment:"Blood sugar unit. (one of: mg/dL, mmol/L)"`
	DynamicIcon DynamicIcon `toml:"dynamic-icon" comment:"Makes the tray icon show the current blood sugar reading."`
	Arrows      Arrows      `toml:"arrows" comment:"Customize the arrows."`
	Socket      Socket      `toml:"socket" comment:"Exposes the latest reading to other applications over a local socket."`
	Log         Log         `toml:"log" comment:"Log configuration"`
	Advanced    Advanced    `toml:"advanced" comment:"Advanced settings."`
}

type DynamicIcon struct {
	Enabled     bool       `toml:"enabled"`
	FontColor   colorx.Hex `toml:"font-color" comment:"Hex code used to render text."`
	FontFile    string     `toml:"font-file" comment:"Font path or filename of a system font. If left blank, an embedded font will be used."`
	MaxFontSize float64    `toml:"max-font-size" comment:"Maximum font size in points."`
}

type Arrows struct {
	DoubleUp      string `toml:"double-up"`
	SingleUp      string `toml:"single-up"`
	FortyFiveUp   string `toml:"forty-five-up"`
	Flat          string `toml:"flat"`
	FortyFiveDown string `toml:"forty-five-down"`
	SingleDown    string `toml:"single-down"`
	DoubleDown    string `toml:"double-down"`
	Unknown       string `toml:"unknown"`
}

type Socket struct {
	Enabled bool   `toml:"enabled"`
	Format  string `toml:"format" comment:"Local file format. (one of: csv)"`
	Path    string `toml:"path" comment:"File path. $TMPDIR will be replaced with the current temp directory."`
}

type Log struct {
	Level  slogx.Level  `toml:"level" comment:"Values: trace, debug, info, warn, error, fatal, panic"`
	Format slogx.Format `toml:"format" comment:"Values: auto, color, plain, json"`
}

type Advanced struct {
	FetchDelay       Duration `toml:"fetch-delay" comment:"Time to wait before the next reading should be ready.\nIn testing, this seems to be about 20s behind, so the default is 30s to be safe.\nYour results may vary."`
	FallbackInterval Duration `toml:"fallback-interval" comment:"Normally, readings will be fetched when ready (after ~5m).\nThis interval will be used if the next reading time cannot be estimated due to sensor warm-up, missed readings, errors, etc."`
	RoundAge         bool     `toml:"round-age" comment:"If enabled, the reading's age will be rounded up to the nearest minute.\nNightscout rounds the age, so enable this if you want the values to match."`
}

const configDir = "nightscout-menu-bar"

func GetDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			return filepath.Join(xdgConfigHome, configDir), nil
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, ".config", configDir), nil
	default:
		dir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}

		dir = filepath.Join(dir, configDir)
		return dir, nil
	}
}
