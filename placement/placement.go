package placement

import (
	"fmt"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/critsprinkler/bubble"
	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/config"
	"github.com/xackery/critsprinkler/library"
	"github.com/xackery/critsprinkler/tracker"
	"github.com/xackery/critsprinkler/util"
)

var (
	ui         *ebitenui.UI
	placements = make(map[common.PopupCategory]*common.Placement)
)

func New(eui *ebitenui.UI, cfg *config.CritSprinklerConfiguration) error {
	var err error
	ui = eui

	placements[common.PopupCategoryMeleeHitOut] = &cfg.MeleeHitOut
	placements[common.PopupCategoryMeleeHitIn] = &cfg.MeleeHitIn
	placements[common.PopupCategoryMeleeCritOut] = &cfg.MeleeCritOut
	placements[common.PopupCategoryMeleeCritIn] = &cfg.MeleeCritIn
	placements[common.PopupCategoryMeleeMissOut] = &cfg.MeleeMissOut
	placements[common.PopupCategoryMeleeMissIn] = &cfg.MeleeMissIn
	placements[common.PopupCategorySpellHitOut] = &cfg.SpellHitOut
	placements[common.PopupCategorySpellHitIn] = &cfg.SpellHitIn
	placements[common.PopupCategorySpellCritOut] = &cfg.SpellCritOut
	placements[common.PopupCategorySpellCritIn] = &cfg.SpellCritIn
	placements[common.PopupCategorySpellMissOut] = &cfg.SpellMissOut
	placements[common.PopupCategorySpellMissIn] = &cfg.SpellMissIn
	placements[common.PopupCategoryHealHitOut] = &cfg.HealHitOut
	placements[common.PopupCategoryHealHitIn] = &cfg.HealHitIn
	placements[common.PopupCategoryHealCritOut] = &cfg.HealCritOut
	placements[common.PopupCategoryHealCritIn] = &cfg.HealCritIn
	placements[common.PopupCategoryRuneHitOut] = &cfg.RuneHitOut
	placements[common.PopupCategoryRuneHitIn] = &cfg.RuneHitIn
	placements[common.PopupCategoryTotalDamageIn] = &cfg.TotalDamageIn
	placements[common.PopupCategoryTotalDamageOut] = &cfg.TotalDamageOut
	placements[common.PopupCategoryTotalHealIn] = &cfg.TotalHealIn
	placements[common.PopupCategoryTotalHealOut] = &cfg.TotalHealOut

	for category := range placements {
		placement := placements[category]
		placement.Category = category
		placement.FontFace, err = library.FontByKey(library.Font(library.FontPopupNotoSansBold42))
		if err != nil {
			return fmt.Errorf("fontByKey: %w", err)
		}
		placement.TitleFontFace, err = library.FontByKey(library.Font(library.FontSmall))
		if err != nil {
			return fmt.Errorf("fontByKey: %w", err)
		}

		if placement.IsVisible > 0 {
			err := Open(category)
			if err != nil {
				return fmt.Errorf("open %s: %w", category.String(), err)
			}
		}
	}

	return nil

}

func SetEditMode(ui *ebitenui.UI, editMode bool) {
	var err error
	state := widget.Visibility_Show
	if !editMode {
		state = widget.Visibility_Hide
		for _, placement := range placements {
			if placement.IsVisible == 0 {
				continue
			}
			placement.PanelContainer.BackgroundImage = nil
		}
	} else {
		for _, placement := range placements {
			if placement.IsVisible == 0 {
				continue
			}
			placement.PanelContainer.BackgroundImage, err = library.NinesliceByKey(library.NineSlicePanelIdle)
			if err != nil {
				fmt.Println("ninesliceByKey", err)
			}
		}
	}
	for _, placement := range placements {
		if placement.Window == nil {
			continue
		}

		placement.TitleBar.GetWidget().Visibility = state
	}
}

func Open(category common.PopupCategory) error {
	placement := placements[category]
	placement.IsVisible = 1

	fmt.Println("opening", category.String(), "with", placement.Category)
	panelNineSlice, err := library.NinesliceByKey(library.NineSlicePanelIdle)
	if err != nil {
		return fmt.Errorf("ninesliceByKey: %w", err)
	}
	titleNineSlice, err := library.NinesliceByKey(library.NineSliceTitlebarIdle)
	if err != nil {
		return fmt.Errorf("ninesliceByKey: %w", err)
	}

	buttonInvisibleImage, err := library.ButtonImageByKey(library.ButtonImageInvisible)
	if err != nil {
		return fmt.Errorf("buttonImageByKey: %w", err)
	}

	buttonCloseImage, err := library.ButtonImageByKey(library.ButtonImageClose)
	if err != nil {
		return fmt.Errorf("buttonImageByKey: %w", err)
	}

	placement.TitleBar = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(titleNineSlice),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{false, true, false}, []bool{true}),
			widget.GridLayoutOpts.Padding(widget.Insets{Left: 10, Right: 5, Top: 0, Bottom: 0}),
		)))
	placement.TitleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(buttonCloseImage),
		widget.ButtonOpts.TextPadding(widget.Insets{Left: 16, Right: 16}),
		//widget.ButtonOpts.TextPadding(widget.Insets{Left: 30, Right: 30}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			bubble.DamageEvent(&common.DamageEvent{
				Category:  category,
				Source:    tracker.PlayerName(),
				Target:    "Test",
				SpellName: "Ice Comet",
				Damage:    "100",
			})
		}),
		widget.ButtonOpts.TabOrder(99),
	))
	placement.TitleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(buttonInvisibleImage),
		//widget.ButtonOpts.TextPadding(widget.Insets{Left: 30, Right: 30}),
		widget.ButtonOpts.Text(placement.Category.String(), placement.TitleFontFace, &widget.ButtonTextColor{
			Idle:     util.HexToColor("dFF4FFFF"),
			Disabled: util.HexToColor("5A7A91FF"),
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {

		}),
		widget.ButtonOpts.TabOrder(99),
	))
	placement.TitleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(buttonCloseImage),
		widget.ButtonOpts.TextPadding(widget.Insets{Left: 16, Right: 16}),
		//widget.ButtonOpts.TextPadding(widget.Insets{Left: 30, Right: 30}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			Close(placement.Category)
		}),
		widget.ButtonOpts.TabOrder(99),
	))

	// sc := widget.NewContainer(
	// 	widget.ContainerOpts.Layout(widget.NewRowLayout(
	// 		widget.RowLayoutOpts.Direction(widget.DirectionVertical),
	// 		widget.RowLayoutOpts.Padding(widget.Insets{
	// 			Top:    20,
	// 			Bottom: 20,
	// 		}))),
	// )

	// titleBar.AddChild(sc)

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(panelNineSlice),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			),
		),
	)
	placement.PanelContainer = c

	placement.Window = widget.NewWindow(
		//widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.TitleBar(placement.TitleBar, 16),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(100, 80),
		//widget.WindowOpts.MaxSize(300, 1000),
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			OnResize()
		}),
	)
	//windowSize := input.GetWindowSize()

	//r := image.Rect(500, 50, 100, 100)
	//r = r.Add(image.Point{windowSize.X / 4 / 2, windowSize.Y * 2 / 3 / 5})

	//window.SetLocation(r)
	placement.Window.SetLocation(*placement.WindowRect)

	_ = ui.AddWindow(placement.Window)

	return nil
}

func Close(category common.PopupCategory) error {
	fmt.Println("closing", category.String())
	placement := placements[category]
	if placement.Window == nil {
		fmt.Println("window is nil")
		return nil
	}
	placement.Window.Close()
	placement.Window = nil
	placement.IsVisible = 0
	return nil
}

// Toggle opens or closes the window
func Toggle(category common.PopupCategory) error {
	window := placements[category].Window
	if window != nil {
		return Close(category)
	}
	return Open(category)
}

func Update(isEditMode bool) {
}

func OnResize() {
	for _, placement := range placements {
		w, h := ebiten.WindowSize()
		if placement.Window == nil {
			continue
		}

		fmt.Println("onresize", placement.Category)
		rect := placement.Window.GetContainer().GetWidget().Rect
		originalRect := rect

		newMinX := util.ClampInt(rect.Min.X, 0, w-rect.Dx())
		newMaxX := newMinX + rect.Dx()

		newMinY := util.ClampInt(rect.Min.Y, 0, h-rect.Dy())
		newMaxY := newMinY + rect.Dy()

		newRect := rect
		newRect.Min.X = newMinX
		newRect.Max.X = newMaxX
		newRect.Min.Y = newMinY
		newRect.Max.Y = newMaxY

		fmt.Println("going from", originalRect, "to", newRect)
		placement.Window.SetLocation(newRect)
		*placement.WindowRect = newRect
	}
}

func buttonImage(spellIconID int) *widget.ButtonImage {
	nineSlice := library.SpellByIDNineSlice(spellIconID)
	if nineSlice == nil {
		return nil
	}

	return &widget.ButtonImage{
		Idle:    nineSlice,
		Pressed: nineSlice,
	}
}

func ByCategory(category common.PopupCategory) *common.Placement {
	return placements[category]
}
