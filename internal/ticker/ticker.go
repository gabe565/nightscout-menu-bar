package ticker

import (
	"context"
	"log/slog"
	"time"

	"fyne.io/fyne/v2"
	"github.com/gabe565/nightscout-menu-bar/internal/fetch"
	"github.com/gabe565/nightscout-menu-bar/internal/localfile"
)

func New(app fyne.App, updateCh chan<- any, version string) *Ticker {
	t := &Ticker{
		app:       app,
		fetch:     fetch.NewFetch(app, version),
		localFile: localfile.New(app),
		bus:       updateCh,
	}
	app.Preferences().AddChangeListener(t.reloadConfig)
	return t
}

type Ticker struct {
	cancel context.CancelFunc

	app       fyne.App
	fetch     *fetch.Fetch
	localFile *localfile.LocalFile

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
	if err := t.localFile.Cleanup(); err != nil {
		slog.Error("Failed to cleanup local file", "error", err)
	}
}
