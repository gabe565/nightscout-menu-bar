package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gabe565.com/nightscout-menu-bar/internal/secret"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeConfig(t *testing.T, dir string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, configFile), nil, 0o644))
}

func TestResolveDir(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	preferred := filepath.Join(tmp, "preferred", configDir)
	legacy := filepath.Join(tmp, "legacy", configDir)
	dirs := []string{preferred, legacy}

	t.Run("none exist", func(t *testing.T) {
		assert.Equal(t, preferred, resolveDir(dirs))
	})

	t.Run("only legacy exists", func(t *testing.T) {
		writeConfig(t, legacy)
		assert.Equal(t, legacy, resolveDir(dirs))
	})

	t.Run("legacy dir exists without config", func(t *testing.T) {
		require.NoError(t, os.Remove(filepath.Join(legacy, configFile)))
		assert.Equal(t, preferred, resolveDir(dirs))
	})

	t.Run("both exist", func(t *testing.T) {
		writeConfig(t, legacy)
		writeConfig(t, preferred)
		assert.Equal(t, preferred, resolveDir(dirs))
	})
}

func TestGetDirs(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin-only")
	}

	home := t.TempDir()
	t.Setenv("HOME", home)

	t.Run("default", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		dirs, err := GetDirs()
		require.NoError(t, err)
		assert.Equal(t, []string{
			filepath.Join(home, "Library", "Application Support", configDir),
			filepath.Join(home, ".config", configDir),
		}, dirs)
	})

	t.Run("xdg", func(t *testing.T) {
		xdg := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", xdg)
		dirs, err := GetDirs()
		require.NoError(t, err)
		assert.Equal(t, []string{
			filepath.Join(home, "Library", "Application Support", configDir),
			filepath.Join(xdg, configDir),
		}, dirs)
	})
}

func TestGetDir(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin-only")
	}

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")

	preferred := filepath.Join(home, "Library", "Application Support", configDir)
	legacy := filepath.Join(home, ".config", configDir)

	t.Run("fresh install", func(t *testing.T) {
		dir, err := GetDir()
		require.NoError(t, err)
		assert.Equal(t, preferred, dir)
	})

	t.Run("existing install", func(t *testing.T) {
		writeConfig(t, legacy)
		dir, err := GetDir()
		require.NoError(t, err)
		assert.Equal(t, legacy, dir)
	})

	t.Run("migrated install", func(t *testing.T) {
		writeConfig(t, preferred)
		dir, err := GetDir()
		require.NoError(t, err)
		assert.Equal(t, preferred, dir)
	})
}

func TestConfig_Token(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	// Older configs kept the token in the file
	require.NoError(t, os.WriteFile(path, []byte("url = 'https://example.com'\ntoken = 'abc'\n"), 0o600))

	conf := New()
	conf.File = path
	conf.Secrets = secret.NewWithBackend(dir, &secret.Memory{})
	require.NoError(t, conf.Load())
	assert.Equal(t, "abc", conf.Data().Token)

	// Moved to the secret store
	token, err := conf.Secrets.Get(TokenKey)
	require.NoError(t, err)
	assert.Equal(t, "abc", token)
	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.NotContains(t, string(contents), "abc")

	// Loaded from the secret store
	reloaded := New()
	reloaded.File = path
	reloaded.Secrets = conf.Secrets
	require.NoError(t, reloaded.Load())
	assert.Equal(t, "abc", reloaded.Data().Token)

	// Cleared tokens are deleted
	data := reloaded.Data()
	data.Token = ""
	require.NoError(t, reloaded.Write(data))
	_, err = conf.Secrets.Get(TokenKey)
	require.ErrorIs(t, err, secret.ErrNotFound)
}
