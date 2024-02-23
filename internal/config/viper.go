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
	flag.StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.config/nightscout-menu-bar.yaml)")
}

func InitViper() error {
	flag.Parse()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configDir := filepath.Join(home, ".config")

		viper.SetConfigName("nightscout-menu-bar")
		viper.SetConfigType("yaml")

		viper.AddConfigPath(configDir)

		if err := viper.SafeWriteConfig(); err != nil {
			if !errors.Is(err, err.(viper.ConfigFileAlreadyExistsError)) {
				return err
			}
		} else {
			log.Println("Created config file")
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("nightscout")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return err
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
