package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	Title     string    `toml:"title"`
	URL       string    `toml:"url"`
	Token     string    `toml:"token"`
	Units     string    `toml:"units"`
	Interval  Duration  `toml:"interval"`
	Arrows    Arrows    `toml:"arrows"`
	LocalFile LocalFile `toml:"local-file"`
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
	Format  string `toml:"format"`
	Path    string `toml:"path"`
	Cleanup bool   `toml:"cleanup"`
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
