package ticker

import (
	"context"
	"errors"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/fetch"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/rs/zerolog/log"
)

func (t *Ticker) beginFetch(render chan<- *nightscout.Properties) {
	go func() {
		t.fetchTicker = time.NewTicker(t.Fetch(render))

		for {
			select {
			case <-t.ctx.Done():
				return
			case <-t.fetchTicker.C:
				next := t.Fetch(render)
				t.fetchTicker.Reset(next)
			}
		}
	}()
}

func (t *Ticker) Fetch(render chan<- *nightscout.Properties) time.Duration {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	properties, err := t.fetch.Do(ctx)
	if err != nil && !errors.Is(err, fetch.ErrNotModified) {
		t.bus <- err
	}
	if properties != nil {
		if render != nil {
			render <- properties
		}
		if t.config.LocalFile.Enabled {
			if err := t.localFile.Write(properties); err != nil {
				log.Err(err).Msg("Failed to write local file")
			}
		}
		if len(properties.Buckets) != 0 {
			bucket := properties.Buckets[0]
			lastDiff := bucket.ToMills.Sub(bucket.FromMills.Time)
			nextRead := properties.Bgnow.Mills.Add(lastDiff + t.config.Advanced.FetchDelay.Duration)
			if until := time.Until(nextRead); until > 0 {
				return until
			}
		}
	}
	return t.config.Interval.Duration
}
