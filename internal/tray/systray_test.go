package tray

import (
	"context"
	"testing"
	"time"

	fyneapp "fyne.io/fyne/v2/app"
	"github.com/gabe565/nightscout-menu-bar/internal/fyneutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	app, ok := fyneapp.New().(fyneutil.DesktopApp)
	require.True(t, ok)
	tray := New(app, "")
	assert.NotNil(t, tray)
	assert.NotNil(t, tray.app)
	assert.NotNil(t, tray.ticker)
}

func TestTray_onError(t *testing.T) {
	t.Parallel()
	app, ok := fyneapp.New().(fyneutil.DesktopApp)
	require.True(t, ok)
	tray := New(app, "")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

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
