package library

import (
	"fmt"

	"github.com/ebitenui/ebitenui/image"
)

type Nineslice int

const (
	NineSlicePanelIdle Nineslice = iota
	NineSliceTitlebarIdle
	NineSliceButtonIdle
	NineSliceButtonHover
	NineSliceButtonSelectedHover
	NineSliceButtonPressed
	NineSliceButtonDisabled
	NineSliceButtonInvisibleIdle
	NineSliceButtonInvisibleHover
	NineSliceButtonInvisiblePressed
	NineSliceButtonInvisibleDisabled
	NineSliceButtonCloseIdle
	NineSliceButtonCloseHover
	NineSliceButtonCloseSelectedHover
	NineSliceButtonClosePressed
	NineSliceButtonCloseDisabled
)

// String returns the string representation of the nineslice
func (e Nineslice) String() string {
	switch e {
	case NineSlicePanelIdle:
		return "NineSlicePanelIdle"
	case NineSliceTitlebarIdle:
		return "NineSliceTitlebarIdle"
	case NineSliceButtonIdle:
		return "NineSliceButtonIdle"
	case NineSliceButtonHover:
		return "NineSliceButtonHover"
	case NineSliceButtonSelectedHover:
		return "NineSliceButtonSelectedHover"
	case NineSliceButtonPressed:
		return "NineSliceButtonPressed"
	case NineSliceButtonDisabled:
		return "NineSliceButtonDisabled"
	case NineSliceButtonInvisibleIdle:
		return "NineSliceButtonInvisibleIdle"
	case NineSliceButtonInvisibleHover:
		return "NineSliceButtonInvisibleHover"
	case NineSliceButtonInvisiblePressed:
		return "NineSliceButtonInvisiblePressed"
	case NineSliceButtonInvisibleDisabled:
		return "NineSliceButtonInvisibleDisabled"
	case NineSliceButtonCloseIdle:
		return "NineSliceButtonCloseIdle"
	case NineSliceButtonCloseHover:
		return "NineSliceButtonCloseHover"
	case NineSliceButtonCloseSelectedHover:
		return "NineSliceButtonCloseSelectedHover"
	case NineSliceButtonClosePressed:
		return "NineSliceButtonClosePressed"
	case NineSliceButtonCloseDisabled:
		return "NineSliceButtonCloseDisabled"

	default:
		return "Unknown"
	}
}

var (
	nineslices = map[Nineslice]*image.NineSlice{}
)

// Nineslice returns a nineslice from the library based on key
func NinesliceByKey(key Nineslice) (*image.NineSlice, error) {
	nineslice, ok := nineslices[key]
	if ok {
		return nineslice, nil
	}
	return nil, fmt.Errorf("nineslice not found: %s", key.String())
}

func ninesliceInit() error {
	assets := []struct {
		key          Nineslice
		image        Image
		centerWidth  int
		centerHeight int
		w            [3]int
		h            [3]int
	}{
		{NineSlicePanelIdle, ImagePanelIdle, 10, 10, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceTitlebarIdle, ImageTitlebarIdle, 10, 10, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonIdle, ButtonIdle, 12, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonHover, ButtonHover, 12, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonSelectedHover, ButtonSelectedHover, 12, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonPressed, ButtonPressed, 12, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonDisabled, ButtonDisabled, 12, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonInvisibleIdle, ButtonInvisibleIdle, 0, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonInvisibleHover, ButtonInvisibleIdle, 0, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonInvisiblePressed, ButtonInvisibleIdle, 0, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonInvisibleDisabled, ButtonInvisibleIdle, 0, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonCloseIdle, ButtonCloseIdle, 0, 0, [3]int{16, 0, 16}, [3]int{16, 0, 16}},
		{NineSliceButtonCloseHover, ButtonCloseHover, 0, 0, [3]int{16, 0, 16}, [3]int{16, 0, 16}},
		{NineSliceButtonCloseSelectedHover, ButtonCloseSelectedHover, 12, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonClosePressed, ButtonClosePressed, 0, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
		{NineSliceButtonCloseDisabled, ButtonCloseDisabled, 0, 0, [3]int{0, 0, 0}, [3]int{0, 0, 0}},
	}

	for _, asset := range assets {
		nineslice, err := ninesliceLoad(asset.image, asset.centerWidth, asset.centerHeight, asset.w, asset.h)
		if err != nil {
			return fmt.Errorf("nineslice load %s: %w", asset.key.String(), err)
		}
		nineslices[asset.key] = nineslice
	}

	return nil
}

func ninesliceLoad(key Image, centerWidth int, centerHeight int, w [3]int, h [3]int) (*image.NineSlice, error) {
	i, err := ImageByKey(key)
	if err != nil {
		return nil, err
	}

	if centerWidth != 0 && centerHeight != 0 {
		cW := i.Bounds().Dx()
		cH := i.Bounds().Dy()
		w = [3]int{(cW - centerWidth) / 2, centerWidth, cW - (cW-centerWidth)/2 - centerWidth}
		h = [3]int{(cH - centerHeight) / 2, centerHeight, cH - (cH-centerHeight)/2 - centerHeight}
	}

	return image.NewNineSlice(i, w, h), nil
}
