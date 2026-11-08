package tray

import (
	"context"
	"errors"
	"log/slog"

	"gabe565.com/nightscout-menu-bar/internal/trayhost"
)

// checkTrayHost offers to enable GNOME's AppIndicator extension when there's
// no system tray. If the extension is already enabled but hasn't started yet,
// e.g. during login, it just waits for it.
func (t *Tray) checkTrayHost(ctx context.Context) {
	if trayhost.Available() {
		return
	}
	if !trayhost.IsGNOME() {
		slog.Warn("No system tray found; the tray icon won't be visible")
		return
	}
	go func() {
		if err := trayhost.EnableExtension(ctx); err != nil {
			if errors.Is(err, trayhost.ErrCanceled) {
				slog.Warn("The AppIndicator extension wasn't installed; the tray icon won't be visible")
				return
			}
			// The install dialog can outlast the D-Bus call, so keep waiting for the tray
			slog.Warn("Failed to enable the AppIndicator extension", "error", err)
		}
		if err := trayhost.Register(ctx); err != nil {
			slog.Error("Failed to register with the system tray", "error", err)
		}
	}()
}
