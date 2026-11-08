package tray

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
)

func TestMain(m *testing.M) {
	keyring.MockInit()
	os.Exit(m.Run())
}

// isolate keeps tests from loading or migrating the real config.
func isolate(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("AppData", filepath.Join(home, "AppData"))
}

func TestNew(t *testing.T) {
	isolate(t)
	tray := New("")
	assert.NotNil(t, tray)
	assert.NotNil(t, tray.config)
	assert.NotNil(t, tray.ticker)
}

func TestTray_onError(t *testing.T) {
	isolate(t)
	tray := New("")

	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	t.Cleanup(cancel)

	go func() {
		select {
		case msg := <-tray.bus:
			err, ok := msg.(error)
			assert.True(t, ok)
			assert.Error(t, err)
		case <-ctx.Done():
			return
		}
	}()

	select {
	case tray.bus <- context.DeadlineExceeded:
	case <-ctx.Done():
	}
	assert.NoError(t, ctx.Err())
}
