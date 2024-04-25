package dynamic_icon

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

var (
	//go:embed Inconsolata_Condensed-Black.ttf
	robotoBold []byte

	mu     sync.Mutex
	face   font.Face
	drawer *font.Drawer
	img    *image.RGBA
)

func init() {
	config.AddReloader(func() {
		mu.Lock()
		defer mu.Unlock()
		if face != nil {
			_ = face.Close()
			systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
		}
		face, drawer, img = nil, nil, nil
	})
}

func Generate(p *nightscout.Properties) ([]byte, error) {
	mu.Lock()
	defer mu.Unlock()

	const width, height = 32, 32

	if face == nil {
		var b []byte
		if config.Default.DynamicIcon.FontFile == "" {
			b = robotoBold
		} else {
			var err error
			if b, err = os.ReadFile(config.Default.DynamicIcon.FontFile); err != nil {
				return nil, err
			}
		}

		f, err := truetype.Parse(b)
		if err != nil {
			return nil, err
		}

		face = truetype.NewFace(f, &truetype.Options{
			Size: config.Default.DynamicIcon.FontSize,
		})

		img = image.NewRGBA(image.Rectangle{Max: image.Point{X: width, Y: height}})

		m := face.Metrics()
		src := image.NewUniform(config.Default.DynamicIcon.FontColor.RGBA())
		drawer = &font.Drawer{
			Dst:  img,
			Src:  src,
			Face: face,
			Dot:  fixed.Point26_6{Y: fixed.I(height) - m.Height/2 + m.Descent},
		}
	} else {
		draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)
	}

	bgnow := p.Bgnow.DisplayBg()
	drawer.Dot.X = (fixed.I(width) - drawer.MeasureString(bgnow)) / 2
	drawer.DrawString(bgnow)

	var buf bytes.Buffer
	if err := encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
