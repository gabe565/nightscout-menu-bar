package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gabe565/nightscout-systray/internal/tray"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
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
		viper.SetConfigName("nightscout-menu-bar")
		viper.SetConfigType("yaml")

		viper.AddConfigPath("$HOME/.config")
		viper.AddConfigPath(".")
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		tray.ResetTicker()
	})

	viper.AutomaticEnv()
	viper.SetEnvPrefix("nightscout")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			return err
		}
	}
	return nil
}
