package main

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	fontFaceRegular = "assets/fonts/notosans-regular.ttf"
	fontFaceBold    = "assets/fonts/notosans-bold.ttf"
)

type fonts struct {
	face         text.Face
	titleFace    text.Face
	bigTitleFace text.Face
	toolTipFace  text.Face
	popupBold42  text.Face
	popupBold36  text.Face
}

func loadFonts() (*fonts, error) {
	fontFace, err := loadFont(fontFaceRegular, 20)
	if err != nil {
		return nil, err
	}

	titleFontFace, err := loadFont(fontFaceBold, 24)
	if err != nil {
		return nil, err
	}

	bigTitleFontFace, err := loadFont(fontFaceBold, 28)
	if err != nil {
		return nil, err
	}

	toolTipFace, err := loadFont(fontFaceRegular, 15)
	if err != nil {
		return nil, err
	}

	popupBold42, err := loadFont(fontFaceBold, 42)
	if err != nil {
		return nil, err
	}
	popupBold36, err := loadFont(fontFaceBold, 36)
	if err != nil {
		return nil, err
	}

	return &fonts{
		face:         fontFace,
		titleFace:    titleFontFace,
		bigTitleFace: bigTitleFontFace,
		toolTipFace:  toolTipFace,
		popupBold42:  popupBold42,
		popupBold36:  popupBold36,
	}, nil
}
