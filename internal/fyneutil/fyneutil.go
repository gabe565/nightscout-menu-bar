package fyneutil

import (
	"bytes"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	svg "github.com/ajstarks/svgo"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

type DesktopApp interface {
	fyne.App
	desktop.App
}

func ColorImage(c util.HexColor) *fyne.StaticResource {
	var buf bytes.Buffer
	canvas := svg.New(&buf)
	canvas.Start(32, 32)
	canvas.Roundrect(0, 0, 32, 32, 4, 4, "fill:"+c.String())
	canvas.End()
	return fyne.NewStaticResource(c.String()+".svg", buf.Bytes())
}

func ValidateDuration(required bool) fyne.StringValidator {
	return func(s string) error {
		if !required && s == "" {
			return nil
		}
		_, err := time.ParseDuration(s)
		return err
	}
}
