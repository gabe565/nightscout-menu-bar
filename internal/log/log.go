package log

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

//go:generate enumer -type Format -trimprefix Format -transform lower -text

type Format uint8

const (
	FormatAuto Format = iota
	FormatColor
	FormatPlain
	FormatJSON
)

func Init(w io.Writer, level slog.Level, format Format) {
	switch format {
	case FormatJSON:
		slog.SetDefault(slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{
				Level: level,
			}),
		))
	default:
		var color bool
		switch format {
		case FormatAuto:
			if f, ok := w.(*os.File); ok {
				color = isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
			}
		case FormatColor:
			color = true
		}

		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      level,
				TimeFormat: time.DateTime,
				NoColor:    !color,
			}),
		))
	}
}
