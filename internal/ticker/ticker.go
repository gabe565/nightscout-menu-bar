package ticker

import (
	"context"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/fetch"
	"github.com/gabe565/nightscout-menu-bar/internal/localfile"
	"github.com/rs/zerolog/log"
)

func New(conf *config.Config, updateCh chan<- any) *Ticker {
	t := &Ticker{
		config:    conf,
		fetch:     fetch.NewFetch(conf),
		localFile: localfile.New(conf),
		bus:       updateCh,
	}

	conf.AddCallback(t.reloadConfig)
	return t
}

type Ticker struct {
	ctx    context.Context
	cancel context.CancelFunc

	config    *config.Config
	fetch     *fetch.Fetch
	localFile *localfile.LocalFile

	fetchTicker  *time.Ticker
	renderTicker *time.Ticker
	bus          chan<- any
}

func (t *Ticker) Start() {
	t.ctx, t.cancel = context.WithCancel(context.Background())
	renderCh := t.beginRender()
	t.beginFetch(renderCh)
	t.beginSleepNotifier()
}

func (t *Ticker) reloadConfig() {
	if t.renderTicker != nil {
		t.renderTicker.Reset(time.Second)
	}
	t.fetch.Reset()
	if t.fetchTicker != nil {
		t.fetchTicker.Reset(time.Second)
	}
}

func (t *Ticker) Close() {
	t.fetchTicker.Stop()
	t.renderTicker.Stop()
	if t.cancel != nil {
		t.cancel()
	}
	if err := t.localFile.Cleanup(); err != nil {
		log.Err(err).Msg("Failed to cleanup local file")
	}
}
