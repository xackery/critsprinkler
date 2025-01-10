package library

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Image int

const (
	ImagePanelIdle Image = iota
	ImageTitlebarIdle
	ButtonIdle
	ButtonHover
	ButtonSelectedHover
	ButtonPressed
	ButtonDisabled
	ButtonInvisibleIdle
	ButtonCloseIdle
	ButtonCloseHover
	ButtonCloseSelectedHover
	ButtonClosePressed
	ButtonCloseDisabled
)

// String returns the string representation of the image
func (e Image) String() string {
	switch e {
	case ImagePanelIdle:
		return "ImagePanelIdle"
	case ImageTitlebarIdle:
		return "ImageTitlebarIdle"
	case ButtonIdle:
		return "ButtonIdle"
	case ButtonHover:
		return "ButtonHover"
	case ButtonSelectedHover:
		return "ButtonSelectedHover"
	case ButtonPressed:
		return "ButtonPressed"
	case ButtonDisabled:
		return "ButtonDisabled"
	case ButtonInvisibleIdle:
		return "ButtonInvisibleIdle"
	case ButtonCloseIdle:
		return "ButtonCloseIdle"
	case ButtonCloseHover:
		return "ButtonCloseHover"
	case ButtonCloseSelectedHover:
		return "ButtonCloseSelectedHover"
	case ButtonClosePressed:
		return "ButtonClosePressed"
	case ButtonCloseDisabled:
		return "ButtonCloseDisabled"

	default:
		return "Unknown"
	}
}

var (
	images = map[Image]*ebiten.Image{}
)

// Image returns a image from the library based on key
func ImageByKey(key Image) (*ebiten.Image, error) {
	image, ok := images[key]
	if ok {
		return image, nil
	}
	return nil, fmt.Errorf("image not found: %s", key.String())
}

func imageInit() error {

	assets := []struct {
		key  Image
		path string
	}{
		{ImagePanelIdle, "assets/graphics/panel-idle.png"},
		{ImageTitlebarIdle, "assets/graphics/titlebar-idle.png"},
		{ButtonIdle, "assets/graphics/button-idle.png"},
		{ButtonHover, "assets/graphics/button-hover.png"},
		{ButtonSelectedHover, "assets/graphics/button-selected-hover.png"},
		{ButtonPressed, "assets/graphics/button-pressed.png"},
		{ButtonDisabled, "assets/graphics/button-disabled.png"},
		{ButtonInvisibleIdle, "assets/graphics/button-invisible-idle.png"},
		{ButtonCloseIdle, "assets/graphics/button-close-idle.png"},
		{ButtonCloseHover, "assets/graphics/button-close-hover.png"},
		{ButtonCloseSelectedHover, "assets/graphics/button-close-selected-hover.png"},
		{ButtonClosePressed, "assets/graphics/button-close-pressed.png"},
		{ButtonCloseDisabled, "assets/graphics/button-close-disabled.png"},
	}

	for _, asset := range assets {
		image, err := imageLoad(asset.path)
		if err != nil {
			return fmt.Errorf("image load %s: %w", asset.key.String(), err)
		}
		images[asset.key] = image
	}

	return nil
}

func imageLoad(path string) (*ebiten.Image, error) {
	f, err := embeddedAssets.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := ebitenutil.NewImageFromReader(f)
	return i, err
}
