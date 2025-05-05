package ticker

import (
	"context"
	"log/slog"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/fetch"
	"gabe565.com/nightscout-menu-bar/internal/socket"
)

func New(conf *config.Config, updateCh chan<- any) *Ticker {
	t := &Ticker{
		config: conf,
		fetch:  fetch.NewFetch(conf),
		socket: socket.New(conf),
		bus:    updateCh,
	}

	conf.AddCallback(t.reloadConfig)
	return t
}

type Ticker struct {
	cancel context.CancelFunc

	config *config.Config
	fetch  *fetch.Fetch
	socket *socket.Socket

	fetchTicker  *time.Ticker
	renderTicker *time.Ticker
	bus          chan<- any
}

func (t *Ticker) Start(ctx context.Context) {
	ctx, t.cancel = context.WithCancel(ctx)
	renderCh := t.beginRender(ctx)
	t.beginFetch(ctx, renderCh)
	t.beginSleepNotifier(ctx)
}

func (t *Ticker) reloadConfig() {
	if t.renderTicker != nil {
		t.renderTicker.Reset(time.Millisecond)
	}
	t.fetch.Reset()
	if t.fetchTicker != nil {
		t.fetchTicker.Reset(time.Millisecond)
	}
}

func (t *Ticker) Close() {
	if t.cancel != nil {
		t.cancel()
	}
	if err := t.socket.Close(); err != nil {
		slog.Error("Failed to cleanup socket", "error", err)
	}
}
