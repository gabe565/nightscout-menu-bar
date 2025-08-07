//go:build !wasm

package config

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/pelletier/go-toml/v2"
)

func (conf *Config) RegisterFlags() {
	conf.Flags.StringVarP(&conf.File, "config", "c", "", "Config file")
}

func (conf *Config) Load() error {
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

	if err := migrateConfig(k); err != nil {
		return err
	}

	data := conf.Data()

	if err := k.UnmarshalWithConf("", &data, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		return err
	}

	if err := conf.Write(data); err != nil {
		return err
	}

	conf.InitLog(os.Stderr)

	slog.Info("Loaded config", "file", conf.File)
	return nil
}

func (conf *Config) Write(data Data) error {
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

	newCfg, err := toml.Marshal(&data)
	if err != nil {
		return err
	}

	if !bytes.Equal(cfgContents, newCfg) {
		logger := slog.With("file", conf.File)
		if cfgNotExists {
			logger.Info("Creating config")

			if err := os.MkdirAll(filepath.Dir(conf.File), 0o777); err != nil {
				return err
			}
		} else {
			logger.Info("Updating config")
		}

		if err := os.WriteFile(conf.File, newCfg, 0o666); err != nil {
			return err
		}
	}

	conf.data.Store(&data)
	return nil
}

func (conf *Config) Watch(ctx context.Context) error {
	logger := slog.With("file", conf.File)
	logger.Info("Watching config")
	f := file.Provider(conf.File)
	return f.Watch(func(_ any, err error) {
		if err != nil {
			logger.Error("Config watcher failed", "error", err)
			if ctx.Err() != nil {
				clear(conf.callbacks)
				conf.callbacks = conf.callbacks[:0]
				return
			}
			time.Sleep(time.Second)
			defer func() {
				_ = conf.Watch(ctx)
			}()
		}

		logger.Debug("Config watcher triggered")
		if err := conf.Load(); err != nil {
			logger.Error("Failed to load config", "error", err)
		}

		for _, fn := range conf.callbacks {
			fn()
		}
	})
}

func (conf *Config) AddCallback(fn func()) int {
	conf.callbacks = append(conf.callbacks, fn)
	return len(conf.callbacks) - 1
}

func (conf *Config) RemoveCallback(idx int) {
	conf.callbacks = slices.Delete(conf.callbacks, idx, idx+1)
}

func migrateConfig(k *koanf.Koanf) error {
	if k.Exists("interval") {
		slog.Info("Migrating config: interval to advanced.fallback-interval")
		if err := k.Set("advanced.fallback-interval", k.Get("interval")); err != nil {
			return err
		}
	}

	return nil
}
