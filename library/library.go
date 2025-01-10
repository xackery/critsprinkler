package library

import (
	"embed"
	"fmt"
	"image"
	"image/png"
)

//go:embed assets
var embeddedAssets embed.FS

// New instantiates the library
func New() error {
	err := fontInit()
	if err != nil {
		return fmt.Errorf("font init: %w", err)
	}
	err = imageInit()
	if err != nil {
		return fmt.Errorf("image init: %w", err)
	}
	err = ninesliceInit()
	if err != nil {
		return fmt.Errorf("nineslice init: %w", err)
	}
	err = buttonImageInit()
	if err != nil {
		return fmt.Errorf("button image init: %w", err)
	}

	return nil
}

func AppIcons() ([]image.Image, error) {
	icons := []image.Image{}
	iconPaths := []string{
		"assets/icon_16.png",
		"assets/icon_24.png",
		"assets/icon_32.png",
	}
	for _, iconPath := range iconPaths {
		r, err := embeddedAssets.Open(iconPath)
		if err != nil {
			return icons, fmt.Errorf("open icon: %w", err)
		}
		defer r.Close()
		img, err := png.Decode(r)
		if err != nil {
			return icons, fmt.Errorf("decode icon: %w", err)
		}
		icons = append(icons, img)
	}
	return icons, nil
}
