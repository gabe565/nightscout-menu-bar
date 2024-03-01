package dynamic_icon

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	"io"
)

type iconDir struct {
	reserved  uint16
	imageType uint16
	numImages uint16
}

type iconDirEntry struct {
	imageWidth   uint8
	imageHeight  uint8
	numColors    uint8
	reserved     uint8
	colorPlanes  uint16
	bitsPerPixel uint16
	sizeInBytes  uint32
	offset       uint32
}

func newIcondir() iconDir {
	return iconDir{
		imageType: 1,
		numImages: 1,
	}
}

func newIcondirentry() iconDirEntry {
	return iconDirEntry{
		colorPlanes:  1,  // windows is supposed to not mind 0 or 1, but other icon files seem to have 1 here
		bitsPerPixel: 32, // can be 24 for bitmap or 24/32 for png. Set to 32 for now
		offset:       22, // 6 iconDir + 16 iconDirEntry, next image will be this image size + 16 iconDirEntry, etc
	}
}

func encode(w io.Writer, img image.Image) error {
	dir := newIcondir()
	entry := newIcondirentry()

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}
	entry.sizeInBytes = uint32(buf.Len())

	bounds := img.Bounds()
	entry.imageWidth = uint8(bounds.Dx())
	entry.imageHeight = uint8(bounds.Dy())

	if err := binary.Write(w, binary.LittleEndian, dir); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, entry); err != nil {
		return err
	}

	if _, err := buf.WriteTo(w); err != nil {
		return err
	}

	return nil
}
