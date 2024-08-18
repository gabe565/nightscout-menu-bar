package dynamicicon

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	"os"
	"sync"

	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

//go:embed Inconsolata_Condensed-Black.ttf
var defaultFont []byte

type DynamicIcon struct {
	config         *config.Config
	mu             sync.Mutex
	callbackOffset int

	face   font.Face
	drawer *font.Drawer
	img    *image.RGBA
}

func New(conf *config.Config) *DynamicIcon {
	d := &DynamicIcon{config: conf}
	d.callbackOffset = conf.AddCallback(d.reloadConfig)
	return d
}

func (d *DynamicIcon) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	var err error
	if d.face != nil {
		err = d.face.Close()
	}
	d.config.RemoveCallback(d.callbackOffset)
	systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
	return err
}

func (d *DynamicIcon) reloadConfig() {
	_ = d.Close()
	d.config.AddCallback(d.reloadConfig)
}

func (d *DynamicIcon) Generate(p *nightscout.Properties) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	const (
		width, height   = 32, 32
		widthF, heightF = fixed.Int26_6(width << 6), fixed.Int26_6(height << 6)
	)

	if d.face == nil {
		var b []byte
		if d.config.DynamicIcon.FontFile == "" {
			b = defaultFont
		} else {
			var err error
			if b, err = os.ReadFile(d.config.DynamicIcon.FontFile); err != nil {
				return nil, err
			}
		}

		f, err := truetype.Parse(b)
		if err != nil {
			return nil, err
		}

		d.face = truetype.NewFace(f, &truetype.Options{
			Size: d.config.DynamicIcon.FontSize,
		})

		d.img = image.NewRGBA(image.Rectangle{Max: image.Point{X: width, Y: height}})

		m := d.face.Metrics()
		src := image.NewUniform(d.config.DynamicIcon.FontColor.RGBA())
		d.drawer = &font.Drawer{
			Dst:  d.img,
			Src:  src,
			Face: d.face,
			Dot:  fixed.Point26_6{Y: (heightF + m.Ascent - m.Descent) / 2},
		}
	} else {
		draw.Draw(d.img, d.img.Bounds(), image.Transparent, image.Point{}, draw.Src)
	}

	bgnow := p.Bgnow.DisplayBg(d.config.Units)
	d.drawer.Dot.X = (widthF - d.drawer.MeasureString(bgnow)) / 2
	d.drawer.DrawString(bgnow)

	var buf bytes.Buffer
	if err := encode(&buf, d.img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
