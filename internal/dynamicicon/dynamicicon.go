package dynamicicon

import (
	"bytes"
	_ "embed"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/nightscout"
	"gabe565.com/nightscout-menu-bar/internal/systheme"
	"gabe565.com/nightscout-menu-bar/internal/util"
	"gabe565.com/utils/bytefmt"
	"github.com/flopp/go-findfont"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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

	font *opentype.Font
}

func New(conf *config.Config) *DynamicIcon {
	return &DynamicIcon{config: conf}
}

var ErrFontSize = errors.New("unable to determine the correct font size")

// Generate renders the icon and encodes it for the system tray.
func (d *DynamicIcon) Generate(p *nightscout.Properties) ([]byte, error) {
	start := time.Now()
	img, err := d.Render(p)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Grow(2 * bytefmt.KiB)
	if err := encode(&buf, img); err != nil {
		return nil, err
	}

	slog.Debug("Generated dynamic icon",
		"took", time.Since(start),
		"size", bytefmt.Encode(int64(buf.Len())),
	)
	return slices.Clip(buf.Bytes()), nil
}

// Color is the icon's text color, which can contrast with the taskbar or menu bar.
func Color(data config.Data) color.Color {
	if data.DynamicIcon.AutoColor {
		return systheme.IconColor()
	}
	return data.DynamicIcon.FontColor
}

// LoadFont reads a font by path, path relative to the config directory, or the
// file name of an installed font. An empty name returns the embedded font.
func LoadFont(name string) ([]byte, error) {
	if name == "" {
		return defaultFont, nil
	}

	path := util.ResolvePath(name)
	if !filepath.IsAbs(path) {
		dir, err := config.GetDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(dir, path)
	}

	b, err := os.ReadFile(path)
	if err == nil || !os.IsNotExist(err) {
		return b, err
	}

	path, findErr := findfont.Find(name)
	if findErr != nil {
		return nil, errors.Join(err, findErr)
	}
	return os.ReadFile(path)
}

// Render draws the icon for the given reading.
func (d *DynamicIcon) Render(p *nightscout.Properties) (*image.RGBA, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	data := d.config.Data()

	if d.font == nil {
		b, err := LoadFont(data.DynamicIcon.FontFile)
		if err != nil {
			return nil, err
		}

		f, err := opentype.Parse(b)
		if err != nil {
			return nil, err
		}

		d.font = f
	}

	bgnow := p.Bgnow.DisplayBg(data.Units)

	var face font.Face
	defer func() {
		if face != nil {
			_ = face.Close()
		}
	}()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawer := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(Color(data)),
	}

	fontSize := data.DynamicIcon.MaxFontSize * 2
	for {
		var err error
		if face, err = opentype.NewFace(d.font, &opentype.FaceOptions{
			Size: fontSize,
			DPI:  72,
		}); err != nil {
			return nil, err
		}

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

	readingX := (widthF - drawer.MeasureString(bgnow)) / 2
	readingY := (heightF + metrics.Ascent - metrics.Descent) / 2
	drawer.Dot = fixed.Point26_6{X: readingX, Y: readingY}
	drawer.DrawString(bgnow)

	if time.Since(p.Bgnow.Mills.Time) > data.LastReading.StaleThreshold.Duration {
		// Strikethrough
		const thickness = 4
		y := readingY.Round() - int(float64(metrics.XHeight)/64/2) - thickness/2
		rect := image.Rect(readingX.Round(), y, int(drawer.Dot.X/64), y+thickness)
		draw.Draw(img, rect, image.NewUniform(color.Black), image.Point{}, draw.Over)
	}

	slog.Debug("Rendered dynamic icon", "font_size", fontSize, "value", bgnow)
	return img, nil
}
