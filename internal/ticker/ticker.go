package ticker

import (
	"context"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/fetch"
	"gabe565.com/nightscout-menu-bar/internal/oauth"
	"gabe565.com/nightscout-menu-bar/internal/socket"
)

func New(conf *config.Config, store *oauth.Store, sock *socket.Socket, updateCh chan<- any) *Ticker {
	t := &Ticker{
		config: conf,
		fetch:  fetch.NewFetch(conf, store),
		socket: sock,
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

func (t *Ticker) Reload() {
	t.reloadConfig()
}

// Rerender redraws the last reading without fetching.
func (t *Ticker) Rerender() {
	if t.renderTicker != nil {
		t.renderTicker.Reset(time.Millisecond)
	}
}

func (t *Ticker) reloadConfig() {
	if t.renderTicker != nil {
		t.renderTicker.Reset(time.Millisecond)
	}
	if t.fetch != nil {
		t.fetch.Reset()
	}
	if t.fetchTicker != nil {
		t.fetchTicker.Reset(time.Millisecond)
	}
}

func (t *Ticker) Close() {
	if t.cancel != nil {
		t.cancel()
	}
}
