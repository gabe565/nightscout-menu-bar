package ticker

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/localfile"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/tray"
)

var fetchTimer = time.NewTimer(0)

func BeginFetch() {
	go func() {
		for range fetchTimer.C {
			Fetch()
			fetchTimer.Reset(config.Default.Interval.Duration)
		}
	}()
}

func Fetch() {
	properties, err := nightscout.Fetch()
	if err != nil && !errors.Is(err, nightscout.ErrNotModified) {
		tray.Error <- err
		return
	}
	if properties != nil {
		RenderCh <- properties
		if config.Default.LocalFile.Enabled {
			if err := localfile.Write(properties); err != nil {
				slog.Error("Failed to write local file", "error", err.Error())
			}
		}
	}
}
