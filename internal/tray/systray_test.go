package tray

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()
	tray := New()
	assert.NotNil(t, tray)
	assert.NotNil(t, tray.config)
	assert.NotNil(t, tray.ticker)
}

func TestTray_onError(t *testing.T) {
	t.Parallel()
	tray := New()

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
