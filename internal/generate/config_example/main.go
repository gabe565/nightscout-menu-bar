package main

import (
	"log/slog"
	"os"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	if err := createConfig(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func createConfig() error {
	conf := config.New()
	conf.InitLog(os.Stderr)

	f, err := os.Create("config_example.toml")
	if err != nil {
		return err
	}

	if err := toml.NewEncoder(f).Encode(conf); err != nil {
		return err
	}

	return f.Close()
}
