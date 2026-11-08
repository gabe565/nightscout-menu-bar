//go:build !darwin && !windows

package trayhost

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	watcher = "org.kde.StatusNotifierWatcher"
	// itemPath is where systray exports the tray icon.
	itemPath        = "/StatusNotifierItem"
	extension       = "appindicatorsupport@rgcjonas.gmail.com"
	shellExtensions = "org.gnome.Shell.Extensions"
)

// Available reports whether a system tray is running.
func Available() bool {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false
	}
	var found bool
	err = conn.BusObject().Call("org.freedesktop.DBus.NameHasOwner", 0, watcher).Store(&found)
	return err == nil && found
}

func IsGNOME() bool {
	return strings.Contains(strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP")), "gnome")
}

// EnableExtension enables GNOME's AppIndicator extension, asking GNOME Shell to
// install it if needed. GNOME shows its own confirmation dialog for installs.
func EnableExtension(ctx context.Context) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	shell := conn.Object("org.gnome.Shell", "/org/gnome/Shell")

	var info map[string]dbus.Variant
	err = shell.CallWithContext(ctx, shellExtensions+".GetExtensionInfo", 0, extension).Store(&info)
	if err != nil {
		return err
	}
	if len(info) != 0 {
		var enabled bool
		err := shell.CallWithContext(ctx, shellExtensions+".EnableExtension", 0, extension).Store(&enabled)
		if err != nil {
			return err
		}
		if !enabled {
			return ErrNotEnabled
		}
		return nil
	}

	var result string
	err = shell.CallWithContext(ctx, shellExtensions+".InstallRemoteExtension", 0, extension).Store(&result)
	if err != nil {
		return err
	}
	if result != "successful" {
		return ErrCanceled
	}
	return nil
}

// Register registers systray's icon once the tray accepts it, waiting until
// ctx is canceled. systray only tries when the tray's name appears, which can be
// before it's ready, e.g. right after the extension is enabled.
// It shares systray's D-Bus connection.
func Register(ctx context.Context) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	tray := conn.Object(watcher, "/StatusNotifierWatcher")

	delay := 500 * time.Millisecond
	for {
		err := tray.CallWithContext(ctx, watcher+".RegisterStatusNotifierItem", 0, itemPath).Err
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return err
		case <-time.After(delay):
		}
		delay = min(delay*2, 5*time.Second)
	}
}
