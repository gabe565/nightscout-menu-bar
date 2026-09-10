package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

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
