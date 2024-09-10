package ticker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/app/settings"
	"github.com/gabe565/nightscout-menu-bar/internal/fetch"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
)

func (t *Ticker) beginFetch(ctx context.Context, render chan<- *nightscout.Properties) {
	go func() {
		t.fetchTicker = time.NewTicker(time.Millisecond)
		defer t.fetchTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.fetchTicker.C:
				next := t.Fetch(render)
				t.fetchTicker.Reset(next)
				slog.Debug("Scheduled next fetch", "in", next)
			}
		}
	}()
}

func (t *Ticker) Fetch(render chan<- *nightscout.Properties) time.Duration {
	properties, err := t.fetch.Do(context.Background())
	if err != nil && !errors.Is(err, fetch.ErrNotModified) {
		t.bus <- err
	}

	prefs := t.app.Preferences()
	if properties != nil {
		if render != nil {
			render <- properties
		}
		if prefs.Bool(settings.LocalEnabledKey) {
			if err := t.localFile.Write(properties); err != nil {
				slog.Error("Failed to write local file", "error", err)
			}
		}
		if len(properties.Buckets) != 0 {
			bucket := properties.Buckets[0]
			lastDiff := bucket.ToMills.Sub(bucket.FromMills.Time)
			delay, err := time.ParseDuration(prefs.String(settings.FetchDelayKey))
			if err != nil {
				delay = settings.FetchDelayDefault
			}
			nextRead := properties.Bgnow.Mills.Add(lastDiff + delay)
			if until := time.Until(nextRead); until > 0 {
				return until
			}
		}
	}

	interval, err := time.ParseDuration(prefs.String(settings.FallbackIntervalKey))
	if err != nil {
		interval = settings.FallbackIntervalDefault
	}
	return interval
}
