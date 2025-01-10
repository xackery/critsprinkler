package common

import (
	"image"
	"image/color"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Placement struct {
	IsVisible      int
	IsTallyEnabled int
	WindowRect     *image.Rectangle
	FontColor      color.RGBA
	Direction      Direction
	Font           Font
	// entries below are not config saved
	Category       PopupCategory
	FontFace       text.Face
	TitleFontFace  text.Face
	Window         *widget.Window
	TitleBar       *widget.Container
	PanelContainer *widget.Container
	LastSpawnX     int
	LastSpawnY     int
}
