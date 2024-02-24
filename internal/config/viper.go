package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gabe565/nightscout-menu-bar/internal/local_file"
	"github.com/gabe565/nightscout-menu-bar/internal/ticker"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string

func init() {
	flag.StringVarP(&cfgFile, "config", "c", "", "Config file (default $HOME/.config/nightscout-menu-bar/config.yaml)")
}

func InitViper() error {
	flag.Parse()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		var configDir string
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			configDir = filepath.Join(xdgConfigHome, "nightscout-menu-bar")
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			configDir = filepath.Join(home, ".config", "nightscout-menu-bar")
		}

		viper.AddConfigPath(configDir)

		if err := os.MkdirAll(configDir, 0o700); err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}

		oldConfigPath := filepath.Join(filepath.Dir(configDir), "nightscout-menu-bar.yaml")
		if _, err := os.Stat(oldConfigPath); err == nil {
			log.Println("Moving config")
			if err := os.Rename(oldConfigPath, filepath.Join(configDir, "config.yaml")); err != nil {
				return err
			}
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("nightscout")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetConfigType("toml")
			if err := viper.SafeWriteConfig(); err != nil {
				if err, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
					return err
				}
			} else {
				log.Println("Created config file")
			}
			_ = viper.ReadInConfig()
		} else {
			return err
		}
	}
	log.Println("Loaded config:", viper.ConfigFileUsed())

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		local_file.ReloadConfig()
		ticker.ReloadConfig()
		tray.ReloadConfig <- struct{}{}
	})

	return nil
}
