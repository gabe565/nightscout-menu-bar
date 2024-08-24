package dynamicicon

import (
	"bytes"
	_ "embed"
	"errors"
	"image"
	"image/draw"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	width, height   = 32, 32
	widthF, heightF = fixed.Int26_6(width << 6), fixed.Int26_6(height << 6)
)

//go:embed Inconsolata_Condensed-Black.ttf
var defaultFont []byte

type DynamicIcon struct {
	config *config.Config
	mu     sync.Mutex

	font *truetype.Font
	img  *image.RGBA
}

func New(conf *config.Config) *DynamicIcon {
	d := &DynamicIcon{
		config: conf,
		img:    image.NewRGBA(image.Rectangle{Max: image.Point{X: width, Y: height}}),
	}
	return d
}

var ErrFontSize = errors.New("unable to determine the correct font size")

func (d *DynamicIcon) Generate(p *nightscout.Properties) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.font == nil {
		var b []byte
		if d.config.DynamicIcon.FontFile == "" {
			b = defaultFont
		} else {
			path := d.config.DynamicIcon.FontFile

			if !filepath.IsAbs(path) {
				dir, err := config.GetDir()
				if err != nil {
					return nil, err
				}

				path = filepath.Join(dir, path)
			}

			var err error
			if b, err = os.ReadFile(path); err != nil {
				return nil, err
			}
		}

		f, err := truetype.Parse(b)
		if err != nil {
			return nil, err
		}

		d.font = f
	}

	start := time.Now()
	bgnow := p.Bgnow.DisplayBg(d.config.Units)

	var face font.Face
	defer func() {
		if face != nil {
			_ = face.Close()
		}
	}()

	drawer := &font.Drawer{
		Dst: d.img,
		Src: image.NewUniform(d.config.DynamicIcon.FontColor.RGBA()),
	}

	fontSize := d.config.DynamicIcon.MaxFontSize
	for {
		face = truetype.NewFace(d.font, &truetype.Options{
			Size: fontSize,
		})
		drawer.Face = face

		if textWidth := drawer.MeasureString(bgnow); textWidth <= widthF {
			break
		}

		_ = face.Close()
		if fontSize <= 1 {
			return nil, ErrFontSize
		}
		fontSize -= 0.5
	}

	metrics := face.Metrics()

	draw.Draw(d.img, d.img.Bounds(), image.Transparent, image.Point{}, draw.Src)
	drawer.Dot.X = (widthF - drawer.MeasureString(bgnow)) / 2
	drawer.Dot.Y = (heightF + metrics.Ascent - metrics.Descent) / 2
	drawer.DrawString(bgnow)

	var buf bytes.Buffer
	if err := encode(&buf, d.img); err != nil {
		return nil, err
	}

	slog.Debug("Generated dynamic icon",
		"took", time.Since(start),
		"font_size", fontSize,
		"value", bgnow,
		"size", buf.Len(),
	)
	return buf.Bytes(), nil
}
