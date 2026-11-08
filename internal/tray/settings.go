package tray

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
)

// openSettings launches the settings window in a separate process so the
// webview's memory is released when it closes. If the window is already open,
// the new process focuses it and exits.
func (t *Tray) openSettings(ctx context.Context) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, executable, "settings", "--config", t.config.File)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	slog.Debug("Opening settings", "cmd", cmd.Args)
	if err := cmd.Start(); err != nil {
		return err
	}

	t.settingsMu.Lock()
	t.settingsCmds[cmd] = struct{}{}
	t.settingsMu.Unlock()

	go func() {
		if err := cmd.Wait(); err != nil {
			slog.Error("Settings window exited", "error", err)
		}
		t.settingsMu.Lock()
		delete(t.settingsCmds, cmd)
		t.settingsMu.Unlock()
	}()
	return nil
}

func (t *Tray) closeSettings() {
	t.settingsMu.Lock()
	defer t.settingsMu.Unlock()
	for cmd := range t.settingsCmds {
		_ = cmd.Process.Kill()
	}
}
