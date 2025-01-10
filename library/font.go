package library

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Font int

const (
	FontDefault Font = iota
	FontSmall
	FontPopupNotoSansBold42
	FontPopupNotoSansBold36
	FontPopupNotoSansBold24
)

// String returns the string representation of the font
func (f Font) String() string {
	switch f {
	case FontDefault:
		return "FontDefault"
	case FontSmall:
		return "FontSmall"
	case FontPopupNotoSansBold42:
		return "FontPopupNotoSansBold42"
	case FontPopupNotoSansBold36:
		return "FontPopupNotoSansBold36"
	case FontPopupNotoSansBold24:
		return "FontPopupNotoSansBold24"
	default:
		return "Unknown"
	}
}

var (
	fonts = map[Font]text.Face{}
)

// Font returns a font from the library based on key
func FontByKey(key Font) (text.Face, error) {
	font, ok := fonts[key]
	if ok {
		return font, nil
	}
	return nil, fmt.Errorf("font not found: %s", key.String())
}

func fontInit() error {
	elements := []struct {
		key  Font
		path string
		size float64
	}{
		{FontDefault, "assets/fonts/notosans-regular.ttf", 20},
		{FontSmall, "assets/fonts/notosans-regular.ttf", 14},
		{FontPopupNotoSansBold42, "assets/fonts/notosans-bold.ttf", 42},
		{FontPopupNotoSansBold36, "assets/fonts/notosans-bold.ttf", 36},
		{FontPopupNotoSansBold24, "assets/fonts/notosans-bold.ttf", 24},
	}

	for _, element := range elements {
		font, err := fontLoad(element.path, element.size)
		if err != nil {
			return fmt.Errorf("font load %s: %w", element.key.String(), err)
		}
		fonts[element.key] = font
	}

	return nil
}

func fontLoad(path string, size float64) (text.Face, error) {
	fontData, err := embeddedAssets.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone, // Disable hinting for a crisp look
	})
	if err != nil {
		return nil, err
	}
	return text.NewGoXFace(face), nil
}
