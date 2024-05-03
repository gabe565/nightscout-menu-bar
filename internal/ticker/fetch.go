package ticker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/fetch"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
)

func (t *Ticker) beginFetch(render chan<- *nightscout.Properties) {
	go func() {
		t.Fetch(render)
		t.fetchTicker = time.NewTicker(t.config.Interval.Duration)

		for {
			select {
			case <-t.ctx.Done():
				return
			case <-t.fetchTicker.C:
				t.Fetch(render)
				t.fetchTicker.Reset(t.config.Interval.Duration)
			}
		}
	}()
}

func (t *Ticker) Fetch(render chan<- *nightscout.Properties) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	properties, err := t.fetch.Do(ctx)
	if err != nil && !errors.Is(err, fetch.ErrNotModified) {
		t.bus <- err
		return
	}
	if properties != nil {
		if render != nil {
			render <- properties
		}
		if t.config.LocalFile.Enabled {
			if err := t.localFile.Write(properties); err != nil {
				slog.Error("Failed to write local file", "error", err.Error())
			}
		}
	}
}
