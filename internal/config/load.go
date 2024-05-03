//go:build !wasm

package config

import (
	"bytes"
	"context"
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

func (conf *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVarP(&conf.File, "config", "c", "", "Config file")
}

func (conf *Config) Load() error {
	flag.Parse()
	k := koanf.New(".")

	// Load conf config
	if err := k.Load(structs.Provider(conf, "toml"), nil); err != nil {
		return err
	}

	// Find config file
	if conf.File == "" {
		cfgDir, err := GetDir()
		if err != nil {
			return err
		}

		conf.File = filepath.Join(cfgDir, "config.toml")
	}

	// Load config file if exists
	cfgContents, err := os.ReadFile(conf.File)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// Parse config file
	parser := TOMLParser{}
	if err := k.Load(rawbytes.Provider(cfgContents), parser); err != nil {
		return err
	}

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		return err
	}

	if err := conf.Write(); err != nil {
		return err
	}

	slog.Info("Loaded config", "file", conf.File)
	return nil
}

func (conf *Config) Write() error {
	// Find config file
	if conf.File == "" {
		cfgDir, err := GetDir()
		if err != nil {
			return err
		}

		conf.File = filepath.Join(cfgDir, "config.toml")
	}

	var cfgNotExists bool
	// Load config file if exists
	cfgContents, err := os.ReadFile(conf.File)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfgNotExists = true
		} else {
			return err
		}
	}

	newCfg, err := toml.Marshal(conf)
	if err != nil {
		return err
	}

	if !bytes.Equal(cfgContents, newCfg) {
		if cfgNotExists {
			slog.Info("Creating config", "file", conf.File)

			if err := os.MkdirAll(filepath.Dir(conf.File), 0o777); err != nil {
				return err
			}
		} else {
			slog.Info("Updating config", "file", conf.File)
		}

		if err := os.WriteFile(conf.File, newCfg, 0o666); err != nil {
			return err
		}
	}

	return nil
}

func (conf *Config) Watch(ctx context.Context) error {
	slog.Info("Watching config", "file", conf.File)
	f := file.Provider(conf.File)
	return f.Watch(func(_ any, err error) {
		if err != nil {
			slog.Error("Config watcher failed", "error", err.Error())
			if ctx.Err() != nil {
				conf.callbacks = nil
				return
			}
			time.Sleep(time.Second)
			defer func() {
				_ = conf.Watch(ctx)
			}()
		}

		if err := conf.Load(); err != nil {
			slog.Error("Failed to load config", "error", err.Error())
		}

		for _, fn := range conf.callbacks {
			fn()
		}
	})
}

func (conf *Config) AddCallback(fn func()) {
	conf.callbacks = append(conf.callbacks, fn)
}
