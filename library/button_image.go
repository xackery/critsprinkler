package library

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
)

type ButtonImage int

const (
	ButtonImageDefault ButtonImage = iota
	ButtonImageInvisible
	ButtonImageClose
)

// String returns the string representation of the Buttonimage
func (e ButtonImage) String() string {
	switch e {
	case ButtonImageDefault:
		return "ButtonImageDefault"
	case ButtonImageInvisible:
		return "ButtonImageInvisible"
	case ButtonImageClose:
		return "ButtonImageClose"
	default:
		return "Unknown"
	}
}

var (
	Buttonimages = map[ButtonImage]*widget.ButtonImage{}
)

// ButtonImage returns a Buttonimage from the library based on key
func ButtonImageByKey(key ButtonImage) (*widget.ButtonImage, error) {
	Buttonimage, ok := Buttonimages[key]
	if ok {
		return Buttonimage, nil
	}
	return nil, fmt.Errorf("buttonimage not found: %s", key.String())
}

func buttonImageInit() error {

	assets := []struct {
		key                  ButtonImage
		idleNineSliceKey     Nineslice
		hoverNineSliceKey    Nineslice
		pressedNineSliceKey  Nineslice
		disabledNineSliceKey Nineslice
	}{
		{ButtonImageDefault, NineSliceButtonIdle, NineSliceButtonHover, NineSliceButtonPressed, NineSliceButtonDisabled},
		{ButtonImageInvisible, NineSliceButtonInvisibleIdle, NineSliceButtonInvisibleIdle, NineSliceButtonInvisibleIdle, NineSliceButtonInvisibleIdle},
		{ButtonImageClose, NineSliceButtonCloseIdle, NineSliceButtonCloseHover, NineSliceButtonClosePressed, NineSliceButtonCloseDisabled},
	}

	for _, asset := range assets {
		Buttonimage, err := buttonImageLoad(asset.idleNineSliceKey, asset.hoverNineSliceKey, asset.pressedNineSliceKey, asset.disabledNineSliceKey)
		if err != nil {
			return fmt.Errorf("buttonimage load %s: %w", asset.key.String(), err)
		}
		Buttonimages[asset.key] = Buttonimage
	}

	return nil
}

func buttonImageLoad(idleNineSliceKey, hoverNineSliceKey, pressedNineSliceKey, disabledNineSliceKey Nineslice) (*widget.ButtonImage, error) {
	idleNineSlice, err := NinesliceByKey(idleNineSliceKey)
	if err != nil {
		return nil, fmt.Errorf("idle nineslice %s: %w", idleNineSliceKey.String(), err)
	}

	hoverNineSlice, err := NinesliceByKey(hoverNineSliceKey)
	if err != nil {
		return nil, fmt.Errorf("hover nineslice %s: %w", hoverNineSliceKey.String(), err)
	}

	pressedNineSlice, err := NinesliceByKey(pressedNineSliceKey)
	if err != nil {
		return nil, fmt.Errorf("pressed nineslice %s: %w", pressedNineSliceKey.String(), err)
	}

	disabledNineSlice, err := NinesliceByKey(disabledNineSliceKey)
	if err != nil {
		return nil, fmt.Errorf("disabled nineslice %s: %w", disabledNineSliceKey.String(), err)
	}

	return &widget.ButtonImage{Idle: idleNineSlice, Hover: hoverNineSlice, Pressed: pressedNineSlice, Disabled: disabledNineSlice}, nil
}
