package dynamicicon

import (
	"bytes"
	_ "embed"
	"errors"
	"image"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/nightscout"
	"gabe565.com/nightscout-menu-bar/internal/util"
	"gabe565.com/utils/bytefmt"
	"github.com/flopp/go-findfont"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	width, height   = 64, 64
	widthF, heightF = fixed.Int26_6(width << 6), fixed.Int26_6(height << 6)
)

//go:embed RobotoCondensed-SemiBold.ttf
var defaultFont []byte

type DynamicIcon struct {
	config *config.Config
	mu     sync.Mutex

	font *truetype.Font
}

func New(conf *config.Config) *DynamicIcon {
	return &DynamicIcon{config: conf}
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
			path := util.ResolvePath(d.config.DynamicIcon.FontFile)

			if !filepath.IsAbs(path) {
				dir, err := config.GetDir()
				if err != nil {
					return nil, err
				}

				path = filepath.Join(dir, path)
			}

			var err error
			if b, err = os.ReadFile(path); err != nil {
				if !os.IsNotExist(err) {
					return nil, err
				}

				path, findErr := findfont.Find(d.config.DynamicIcon.FontFile)
				if findErr != nil {
					return nil, errors.Join(err, findErr)
				}

				if b, err = os.ReadFile(path); err != nil {
					return nil, err
				}
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

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawer := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(d.config.DynamicIcon.FontColor),
	}

	fontSize := d.config.DynamicIcon.MaxFontSize * 2
	for {
		face = truetype.NewFace(d.font, &truetype.Options{
			Size: fontSize,
		})
		drawer.Face = face

		if textWidth := drawer.MeasureString(bgnow); textWidth <= widthF+fixed.I(2) {
			break
		}

		_ = face.Close()
		if fontSize <= 1 {
			return nil, ErrFontSize
		}
		fontSize -= 1.0
	}

	metrics := face.Metrics()

	drawer.Dot.X = (widthF - drawer.MeasureString(bgnow)) / 2
	drawer.Dot.Y = (heightF + metrics.Ascent - metrics.Descent) / 2
	drawer.DrawString(bgnow)

	var buf bytes.Buffer
	buf.Grow(2 * bytefmt.KiB)
	if err := encode(&buf, img); err != nil {
		return nil, err
	}

	slog.Debug("Generated dynamic icon",
		"took", time.Since(start),
		"size", bytefmt.Encode(int64(buf.Len())),
		"font_size", fontSize,
		"value", bgnow,
		"size", buf.Len(),
	)
	return slices.Clip(buf.Bytes()), nil
}
