//go:build !windows

package dynamicicon

import (
	"image"
	"image/png"
	"io"
)

func encode(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}
