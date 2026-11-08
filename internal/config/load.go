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

	"gabe565.com/nightscout-menu-bar/internal/secret"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/pelletier/go-toml/v2"
)

const TokenKey = "api-token"

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

		conf.File = filepath.Join(cfgDir, configFile)
	}

	if conf.Secrets == nil {
		conf.Secrets = secret.New(filepath.Dir(conf.File))
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

	// Older configs kept the token in the file. It's moved to the secret store by Write.
	data.Token = k.String("token")
	if data.Token == "" {
		token, err := conf.Secrets.Get(TokenKey)
		if err != nil && !errors.Is(err, secret.ErrNotFound) {
			slog.Error("Failed to load token", "error", err)
		}
		data.Token = token
		conf.tokenMu.Lock()
		conf.storedToken = token
		conf.tokenMu.Unlock()
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

		conf.File = filepath.Join(cfgDir, configFile)
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

	// Saved before the file so migrating the token out of the file can't lose it
	if err := conf.saveToken(data.Token); err != nil {
		return err
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
		conf.reload()
	})
}

// WatchSecrets reloads the config when another process changes the token.
func (conf *Config) WatchSecrets(ctx context.Context) error {
	if conf.Secrets == nil {
		return nil
	}
	return conf.Secrets.Watch(ctx, func() {
		token, err := conf.Secrets.Get(TokenKey)
		if err != nil && !errors.Is(err, secret.ErrNotFound) {
			slog.Error("Failed to load token", "error", err)
			return
		}
		conf.tokenMu.Lock()
		changed := token != conf.storedToken
		conf.tokenMu.Unlock()
		if changed {
			slog.Debug("Token changed")
			conf.reload()
		}
	})
}

func (conf *Config) reload() {
	if err := conf.Load(); err != nil {
		slog.Error("Failed to load config", "file", conf.File, "error", err)
	}
	for _, fn := range conf.callbacks {
		fn()
	}
}

func (conf *Config) saveToken(token string) error {
	if conf.Secrets == nil {
		return nil
	}
	conf.tokenMu.Lock()
	defer conf.tokenMu.Unlock()
	if token == conf.storedToken {
		return nil
	}

	var err error
	if token == "" {
		err = conf.Secrets.Delete(TokenKey)
	} else {
		err = conf.Secrets.Set(TokenKey, token)
	}
	if err != nil {
		return err
	}
	conf.storedToken = token
	return nil
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
