//go:build !windows

package dynamic_icon

import (
	"image"
	"image/png"
	"io"
)

func encode(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}
