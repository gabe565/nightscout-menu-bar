package ticker

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/fetch"
	"gabe565.com/nightscout-menu-bar/internal/nightscout"
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
	data := t.config.Data()
	if properties != nil {
		if render != nil {
			render <- properties
		}
		if data.Socket.Enabled {
			t.socket.Write(properties)
		}
		if len(properties.Buckets) != 0 {
			bucket := properties.Buckets[0]
			lastDiff := bucket.ToMills.Sub(bucket.FromMills.Time)
			nextRead := properties.Bgnow.Mills.Add(lastDiff + data.Advanced.FetchDelay.Duration)
			if until := time.Until(nextRead); until > 0 {
				return until
			}
		}
	}
	return data.Advanced.FallbackInterval.Duration
}
