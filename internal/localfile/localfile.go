package localfile

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"github.com/gabe565/nightscout-menu-bar/internal/app/settings"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

func New(app fyne.App) *LocalFile {
	l := &LocalFile{
		app: app,
	}
	l.reloadConfig()
	app.Preferences().AddChangeListener(l.reloadConfig)
	return l
}

type LocalFile struct {
	app  fyne.App
	path string
}

func (l *LocalFile) Format(last *nightscout.Properties) string {
	prefs := l.app.Preferences()
	var units settings.Unit
	if err := units.UnmarshalText([]byte(prefs.String(settings.UnitsKey))); err != nil {
		units = settings.UnitMgdl
	}

	return last.Bgnow.DisplayBg(units) + "," +
		last.Bgnow.Arrow() + "," +
		last.Delta.Display(units) + "," +
		strconv.Itoa(int(last.Bgnow.Mills.Time.Unix())) +
		"\n"
}

func (l *LocalFile) reloadConfig() {
	var path string
	prefs := l.app.Preferences()
	if prefs.Bool(settings.LocalEnabledKey) {
		path = util.ResolvePath(prefs.String(settings.LocalPathKey))
	}
	if l.path != "" && path != l.path {
		if err := l.Cleanup(); err != nil {
			slog.Error("Failed to cleanup local file", "error", err)
		}
	}
	l.path = path
}

func (l *LocalFile) Write(last *nightscout.Properties) error {
	if l.path != "" {
		data := l.Format(last)
		slog.Debug("Writing local file", "data", data)
		err := os.WriteFile(l.path, []byte(data), 0o600)
		return err
	}
	return nil
}

func (l *LocalFile) Cleanup() error {
	if l.path != "" {
		slog.Debug("Removing local file", "path", l.path)
		if err := os.Remove(l.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}
