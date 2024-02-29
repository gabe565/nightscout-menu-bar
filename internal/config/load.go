//go:build !wasm

package config

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/pelletier/go-toml/v2"
	flag "github.com/spf13/pflag"
)

var cfgFile string

func init() {
	flag.StringVarP(&cfgFile, "config", "c", "", "Config file")
}

func Load() error {
	flag.Parse()
	k := koanf.New(".")

	// Load default config
	if err := k.Load(structs.Provider(Default, "toml"), nil); err != nil {
		return err
	}

	// Find config file
	if cfgFile == "" {
		cfgDir, err := GetDir()
		if err != nil {
			return err
		}

		cfgFile = filepath.Join(cfgDir, "config.toml")
	}

	// Load config file if exists
	cfgContents, err := os.ReadFile(cfgFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// Parse config file
	parser := TOMLParser{}
	if err := k.Load(rawbytes.Provider(cfgContents), parser); err != nil {
		return err
	}

	if err := k.UnmarshalWithConf("", &Default, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		return err
	}

	if err := Write(); err != nil {
		return err
	}

	slog.Info("Loaded config", "file", cfgFile)
	return err
}

func Write() error {
	// Find config file
	if cfgFile == "" {
		cfgDir, err := GetDir()
		if err != nil {
			return err
		}

		cfgFile = filepath.Join(cfgDir, "config.toml")
	}

	var cfgNotExists bool
	// Load config file if exists
	cfgContents, err := os.ReadFile(cfgFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfgNotExists = true
		} else {
			return err
		}
	}

	newCfg, err := toml.Marshal(Default)
	if err != nil {
		return err
	}

	if !bytes.Equal(cfgContents, newCfg) {
		if cfgNotExists {
			slog.Info("Creating config", "file", cfgFile)

			if err := os.MkdirAll(filepath.Dir(cfgFile), 0o777); err != nil {
				return err
			}
		} else {
			slog.Info("Updating config", "file", cfgFile)
		}

		if err := os.WriteFile(cfgFile, newCfg, 0o666); err != nil {
			return err
		}
	}

	return err
}

func Watch() error {
	slog.Info("Watching config", "file", cfgFile)
	f := file.Provider(cfgFile)
	return f.Watch(func(event interface{}, err error) {
		if err != nil {
			slog.Error("Config watcher failed", "error", err.Error())
			time.Sleep(time.Second)
			defer func() {
				_ = Watch()
			}()
		}

		if err := Load(); err != nil {
			slog.Error("Failed to load config", "error", err.Error())
		}
		for _, reloader := range reloaders {
			reloader()
		}
	})
}
