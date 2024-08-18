//go:build !wasm

package config

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/mattn/go-isatty"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (conf *Config) RegisterFlags() {
	conf.Flags.StringVarP(&conf.File, "config", "c", "", "Config file")
}

func (conf *Config) Load() error {
	initLogFormat()
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

	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "toml"}); err != nil {
		return err
	}

	if err := conf.Write(); err != nil {
		return err
	}

	initLogLevel(conf)
	log.Info().Str("file", conf.File).Msg("Loaded config")
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
			log.Info().Str("file", conf.File).Msg("Creating config")

			if err := os.MkdirAll(filepath.Dir(conf.File), 0o777); err != nil {
				return err
			}
		} else {
			log.Info().Str("file", conf.File).Msg("Updating config")
		}

		if err := os.WriteFile(conf.File, newCfg, 0o666); err != nil {
			return err
		}
	}

	return nil
}

func (conf *Config) Watch(ctx context.Context) error {
	log.Info().Str("file", conf.File).Msg("Watching config")
	f := file.Provider(conf.File)
	return f.Watch(func(_ any, err error) {
		if err != nil {
			log.Err(err).Msg("Config watcher failed")
			if ctx.Err() != nil {
				conf.callbacks = nil
				return
			}
			time.Sleep(time.Second)
			defer func() {
				_ = conf.Watch(ctx)
			}()
		}

		log.Trace().Msg("Config watcher triggered")
		if err := conf.Load(); err != nil {
			log.Err(err).Msg("Failed to load config")
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

func initLogLevel(conf *Config) {
	level, err := zerolog.ParseLevel(conf.Log.Level)
	if err != nil {
		log.Warn().Msg("Invalid log level. Defaulting to info.")
	}
	zerolog.SetGlobalLevel(level)
}

func initLogFormat() {
	useColor := isatty.IsTerminal(os.Stderr.Fd())
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    !useColor,
		TimeFormat: time.DateTime,
	})
}

func migrateConfig(k *koanf.Koanf) error {
	if k.Exists("interval") {
		log.Info().Msg("Migrating config: interval to advanced.fallback-interval")
		if err := k.Set("advanced.fallback-interval", k.Get("interval")); err != nil {
			return err
		}
	}

	return nil
}
