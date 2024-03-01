package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Title       string      `toml:"title" comment:"Tray title."`
	URL         string      `toml:"url" comment:"Nightscout URL. (required)"`
	Token       string      `toml:"token" comment:"Nightscout token. Using an access token is recommended instead of the API secret."`
	Units       string      `toml:"units" comment:"Blood sugar unit. (one of: mg/dL, mmol/L)"`
	DynamicIcon DynamicIcon `toml:"dynamic_icon" comment:"Makes the tray icon show the current blood sugar reading."`
	Interval    Duration    `toml:"interval" comment:"Update interval."`
	Arrows      Arrows      `toml:"arrows" comment:"Customize the arrows."`
	LocalFile   LocalFile   `toml:"local-file" comment:"Enables writing the latest blood sugar to a local temporary file."`
}

type DynamicIcon struct {
	Enabled   bool     `toml:"enabled"`
	FontColor HexColor `toml:"font_color" comment:"Hex code used to render text."`
	FontFile  string   `toml:"font_file" comment:"If left blank, an embedded font will be used."`
	FontSize  float64  `toml:"font_size" comment:"Font size in points."`
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

type LocalFile struct {
	Enabled bool   `toml:"enabled"`
	Format  string `toml:"format" comment:"Local file format. (one of: csv)"`
	Path    string `toml:"path" comment:"Local file path. $TMPDIR will be replaced with the current temp directory."`
	Cleanup bool   `toml:"cleanup" comment:"If enabled, the local file will be cleaned up when Nightscout Menu Bar is closed."`
}

var configDir = "nightscout-menu-bar"

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
