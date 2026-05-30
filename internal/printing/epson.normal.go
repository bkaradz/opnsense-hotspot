//go:build !tinygo

// This file contains all code that is not compatible with TinyGo

package printing

import (
	"image"
)

// Image prints a raster image
//
// The image must be narrower than the printer's pixel width
func Image(img image.Image) string {
	xL, xH, yL, yH, imgData := printImage(img)
	return ("\x1dv\x30\x00" + string(append([]byte{xL, xH, yL, yH}, imgData...)))
}
