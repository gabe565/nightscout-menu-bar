//go:build !darwin && !windows

package systheme

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
)

const (
	portalDest     = "org.freedesktop.portal.Desktop"
	portalPath     = "/org/freedesktop/portal/desktop"
	portalSettings = "org.freedesktop.portal.Settings"
	appearance     = "org.freedesktop.appearance"
	colorSchemeKey = "color-scheme"
)

// Color scheme values from the desktop portal.
const (
	noPreference uint32 = iota
	preferDark
	preferLight
)

// Light is false since panels are usually dark, even in light themes,
// and no desktop reports the panel's theme.
func Light() bool { return false }

// MenuLight reports whether tray menus are light, based on the desktop's color scheme.
func MenuLight() bool {
	return menuLight(colorScheme(), os.Getenv("XDG_CURRENT_DESKTOP"))
}

func menuLight(scheme uint32, desktop string) bool {
	switch scheme {
	case preferDark:
		return false
	case preferLight:
		return true
	}
	// GNOME's light style has no preference, but its shell menus stay dark.
	// Ubuntu's menus follow the style, and most other desktops default to light.
	desktop = strings.ToLower(desktop)
	return !strings.Contains(desktop, "gnome") || strings.Contains(desktop, "ubuntu")
}

func colorScheme() uint32 {
	conn, err := dbus.SessionBus()
	if err != nil {
		return noPreference
	}
	var v dbus.Variant
	obj := conn.Object(portalDest, portalPath)
	if err := obj.Call(portalSettings+".ReadOne", 0, appearance, colorSchemeKey).Store(&v); err != nil {
		// Older portals only have Read, which wraps the value in another variant
		if err := obj.Call(portalSettings+".Read", 0, appearance, colorSchemeKey).Store(&v); err != nil {
			return noPreference
		}
		if inner, ok := v.Value().(dbus.Variant); ok {
			v = inner
		}
	}
	scheme, _ := v.Value().(uint32)
	return scheme
}

// Watch calls onChange when the color scheme changes, until ctx is canceled.
func Watch(ctx context.Context, onChange func()) error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return err
	}
	if err := conn.AddMatchSignal(
		dbus.WithMatchObjectPath(portalPath),
		dbus.WithMatchInterface(portalSettings),
		dbus.WithMatchMember("SettingChanged"),
	); err != nil {
		_ = conn.Close()
		return err
	}

	signals := make(chan *dbus.Signal, 8)
	conn.Signal(signals)
	go func() {
		defer func() {
			_ = conn.Close()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case sig, ok := <-signals:
				if !ok {
					return
				}
				if len(sig.Body) >= 2 && sig.Body[0] == appearance && sig.Body[1] == colorSchemeKey {
					slog.Debug("Color scheme changed")
					onChange()
				}
			}
		}
	}()
	return nil
}
