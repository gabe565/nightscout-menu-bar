//go:build !windows

package dynamicicon

import (
	"image"
	"image/png"
	"io"
)

func encode(w io.Writer, img image.Image) error {
	encoder := png.Encoder{CompressionLevel: png.BestSpeed}
	return encoder.Encode(w, img)
}
