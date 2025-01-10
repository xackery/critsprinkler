package util

import (
	"image/color"
	"strconv"
)

// HexToColor converts a hex string to a color
func HexToColor(rgba string) color.Color {
	if len(rgba) == 6 {
		rgba = rgba + "FF"
	}

	u, err := strconv.ParseUint(rgba, 16, 0)
	if err != nil {
		panic(err)
	}

	return color.NRGBA{
		R: uint8(u >> 24),
		G: uint8(u >> 16),
		B: uint8(u >> 8),
		A: uint8(u),
	}
}
