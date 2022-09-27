package main

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string

func init() {
	flag.StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.config/nightscout-menu-bar.yaml)")
}

func initViper() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("nightscout-menu-bar")
		viper.SetConfigType("yaml")

		viper.AddConfigPath("$HOME/.config")
		viper.AddConfigPath(".")
	}

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
