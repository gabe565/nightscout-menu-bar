package config

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"strconv"
)

var (
	ErrMissingPrefix = errors.New(`hex code missing "#" prefix`)
	ErrInvalidLength = errors.New("hex code should be 4 or 7 characters")
)

type HexColor color.NRGBA

func (h HexColor) MarshalText() ([]byte, error) {
	shorthand := h.R>>4 == h.R&0xF && h.G>>4 == h.G&0xF && h.B>>4 == h.B&0xF
	if shorthand {
		return []byte(fmt.Sprintf("#%x%x%x", h.R&0xF, h.G&0xF, h.B&0xF)), nil
	}
	return []byte(fmt.Sprintf("#%02x%02x%02x", h.R, h.G, h.B)), nil
}

func (h *HexColor) UnmarshalText(text []byte) error {
	if !bytes.HasPrefix(text, []byte("#")) {
		return ErrMissingPrefix
	}
	switch len(text) {
	case 4, 7:
	default:
		return ErrInvalidLength
	}

	parsed, err := strconv.ParseUint(string(text[1:]), 16, 32)
	if err != nil {
		return err
	}

	//nolint:gosec
	if parsed > 0xFFF {
		h.R = uint8(parsed >> 16 & 0xFF)
		h.G = uint8(parsed >> 8 & 0xFF)
		h.B = uint8(parsed & 0xFF)
	} else {
		h.R = uint8(parsed >> 8 & 0xF)
		h.R |= h.R << 4
		h.G = uint8(parsed >> 4 & 0xF)
		h.G |= h.G << 4
		h.B = uint8(parsed & 0xF)
		h.B |= h.B << 4
	}
	h.A = 0xFF
	return nil
}

func (h HexColor) RGBA() color.RGBA {
	return color.RGBA(h)
}

func White() HexColor {
	return HexColor{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
}

func Black() HexColor {
	return HexColor{A: 0xFF}
}
