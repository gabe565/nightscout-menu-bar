package systheme

import (
	"context"
	"log/slog"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const personalizeKey = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`

// Light reports whether Windows uses a light theme for the taskbar and tray menus.
func Light() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, personalizeKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer func() {
		_ = k.Close()
	}()
	v, _, err := k.GetIntegerValue("SystemUsesLightTheme")
	return err == nil && v == 1
}

// MenuLight reports whether tray menus are light. They follow the taskbar theme.
func MenuLight() bool { return Light() }

// Watch calls onChange when the theme changes, until ctx is canceled.
func Watch(ctx context.Context, onChange func()) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, personalizeKey, registry.NOTIFY)
	if err != nil {
		return err
	}
	event, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		_ = k.Close()
		return err
	}

	go func() {
		defer func() {
			_ = windows.CloseHandle(event)
			_ = k.Close()
		}()

		light := Light()
		for {
			err := windows.RegNotifyChangeKeyValue(windows.Handle(k), false, windows.REG_NOTIFY_CHANGE_LAST_SET, event, true)
			if err != nil {
				slog.Error("Failed to watch system theme", "error", err)
				return
			}
			// Wake up periodically to notice cancellation
			for {
				if ctx.Err() != nil {
					return
				}
				if res, _ := windows.WaitForSingleObject(event, 1000); res == windows.WAIT_OBJECT_0 {
					break
				}
			}
			if l := Light(); l != light {
				light = l
				onChange()
			}
		}
	}()
	return nil
}
