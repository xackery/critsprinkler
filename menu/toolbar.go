package menu

import (
	"fmt"
	goimage "image"
	"image/color"
	"path/filepath"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/config"
	"github.com/xackery/critsprinkler/dialog"
	"github.com/xackery/critsprinkler/library"
	"github.com/xackery/critsprinkler/money"
	"github.com/xackery/critsprinkler/placement"
	"github.com/xackery/critsprinkler/status"
	"golang.org/x/image/colornames"
)

var (
	toolbar     *toolbarStruct
	ui          *ebitenui.UI
	defaultFont text.Face
)

// NOTE: It's not strictly necessary to store references to all the buttons in the toolbarStruct struct, but this example does
// so for completeness' sake. When you keep a reference to buttons in the struct, you can later configure them to respond
// to certain events in your application, and keep your program's logic outside the toolbarStruct.
type toolbarStruct struct {
	container               *widget.Container
	mnuFile                 *widget.Button
	btnFileLoadEQLog        *widget.Button
	btnFileSave             *widget.Button
	btnFileQuit             *widget.Button
	menuSettings            *widget.Button
	btnFullscreenBorderless *widget.Button
	mnuMelee                *widget.Button
	btnMeleeHitOut          *widget.Button
	btnMeleeHitIn           *widget.Button
	btnMeleeCritOut         *widget.Button
	btnMeleeCritIn          *widget.Button
	btnMeleeMissOut         *widget.Button
	btnMeleeMissIn          *widget.Button
	mnuSpell                *widget.Button
	btnSpellHitOut          *widget.Button
	btnSpellHitIn           *widget.Button
	btnSpellCritOut         *widget.Button
	btnSpellCritIn          *widget.Button
	btnSpellMissOut         *widget.Button
	btnSpellMissIn          *widget.Button
	mnuHeal                 *widget.Button
	btnHealHitOut           *widget.Button
	btnHealHitIn            *widget.Button
	btnHealCritOut          *widget.Button
	btnHealCritIn           *widget.Button
	mnuRune                 *widget.Button
	btnRuneHitOut           *widget.Button
	btnRuneHitIn            *widget.Button
	mnuTotal                *widget.Button
	btnTotalDamageOut       *widget.Button
	btnTotalDamageIn        *widget.Button
	btnTotalHealOut         *widget.Button
	btnTotalHealIn          *widget.Button
	mnuExtra                *widget.Button
	btnMoney                *widget.Button
}

func toolbarNew(cfg *config.CritSprinklerConfiguration, eui *ebitenui.UI) (*toolbarStruct, error) {
	var err error
	toolbar = &toolbarStruct{}
	ui = eui

	// Create a root container for the toolbar.
	toolbar.container = widget.NewContainer(
		// Use black background for the toolbar.
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.Black)),

		// Toolbar components must be aligned horizontally.
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			),
		),

		widget.ContainerOpts.WidgetOpts(
			// Make the toolbar fill the whole horizontal space of the screen.
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{StretchHorizontal: true}),
		),
	)

	defaultFont, err = library.FontByKey(library.FontDefault)
	if err != nil {
		return nil, err
	}

	toolbar.mnuFile = toolbarButtonNew("File", defaultFont)
	toolbar.btnFileLoadEQLog = toolbarButtonNew("Load EQ Log", defaultFont)
	toolbar.btnFileSave = toolbarButtonNew("Save", defaultFont)
	toolbar.btnFileQuit = toolbarButtonNew("Quit", defaultFont)
	toolbar.btnFileLoadEQLog.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			path, err := dialog.FileDialogBox(cfg.LogPath)
			if err != nil && err.Error() == "cancelled" {
				dialog.MsgBox("Error", fmt.Sprintf("Error loading log: %v", err))
				return
			}
			cfg.LogPath = path
			cfg.EQPath = filepath.Dir(fmt.Sprintf("%s/../", filepath.Dir(cfg.LogPath)))
			err = eqPathLoad(cfg)
			if err != nil {
				dialog.MsgBox("Error", fmt.Sprintf("Error loading eq assets: %v", err))
				return
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
			if cfg.LogPath == "" {
				status.Set("Load EQ Log file (Currently not set)")
			} else {
				logPath := filepath.Base(cfg.LogPath)
				status.Setf("Load EQ Log file (Currently %s)", logPath)
			}
		}),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { status.Set("") }),
	)
	// Make the toolbar entry open a menu with our "save" and "load" entries  when the user clicks it.
	toolbar.mnuFile.Configure(
		// Make the toolbar entry open a menu with our "save" and "load" entries  when the user clicks it.
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			toolbarMenuOpen(args.Button.GetWidget(), ui, toolbar.btnFileLoadEQLog, toolbar.btnFileSave, toolbar.btnFileQuit)
		}),
	)
	toolbar.container.AddChild(toolbar.mnuFile)

	toolbar.menuSettings = toolbarButtonNew("Settings", defaultFont)
	toolbar.btnFullscreenBorderless = toolbarButtonNew("Fullscreen Borderless", defaultFont)
	toolbar.btnFullscreenBorderless.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			toggleFullscreenBorderless(cfg)
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { status.Set("Toggle Fullscreen Borderless") }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { status.Set("") }),
	)

	toolbar.menuSettings.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			toolbarMenuOpen(args.Button.GetWidget(), ui, toolbar.btnFullscreenBorderless)
		}),
	)
	toolbar.container.AddChild(toolbar.menuSettings)

	elements := []struct {
		category common.PopupCategory
		button   **widget.Button
	}{
		{common.PopupCategoryMeleeHitOut, &toolbar.btnMeleeHitOut},
		{common.PopupCategoryMeleeHitIn, &toolbar.btnMeleeHitIn},
		{common.PopupCategoryMeleeCritOut, &toolbar.btnMeleeCritOut},
		{common.PopupCategoryMeleeCritIn, &toolbar.btnMeleeCritIn},
		{common.PopupCategoryMeleeMissOut, &toolbar.btnMeleeMissOut},
		{common.PopupCategoryMeleeMissIn, &toolbar.btnMeleeMissIn},
		{common.PopupCategorySpellHitOut, &toolbar.btnSpellHitOut},
		{common.PopupCategorySpellHitIn, &toolbar.btnSpellHitIn},
		{common.PopupCategorySpellCritOut, &toolbar.btnSpellCritOut},
		{common.PopupCategorySpellCritIn, &toolbar.btnSpellCritIn},
		{common.PopupCategorySpellMissOut, &toolbar.btnSpellMissOut},
		{common.PopupCategorySpellMissIn, &toolbar.btnSpellMissIn},
		{common.PopupCategoryHealHitOut, &toolbar.btnHealHitOut},
		{common.PopupCategoryHealHitIn, &toolbar.btnHealHitIn},
		{common.PopupCategoryHealCritOut, &toolbar.btnHealCritOut},
		{common.PopupCategoryHealCritIn, &toolbar.btnHealCritIn},
		{common.PopupCategoryRuneHitOut, &toolbar.btnRuneHitOut},
		{common.PopupCategoryRuneHitIn, &toolbar.btnRuneHitIn},
		{common.PopupCategoryTotalDamageOut, &toolbar.btnTotalDamageOut},
		{common.PopupCategoryTotalDamageIn, &toolbar.btnTotalDamageIn},
		{common.PopupCategoryTotalHealOut, &toolbar.btnTotalHealOut},
		{common.PopupCategoryTotalHealIn, &toolbar.btnTotalHealIn},
	}

	for _, element := range elements {
		*element.button = toolbarButtonNew(element.category.String(), defaultFont)
		button := *element.button
		button.Configure(
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				placement.Toggle(element.category)
			}),
			widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { status.Setf("Toggle %s", element.category.String()) }),
			widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { status.Set("") }),
		)
	}

	toolbar.mnuMelee = toolbarButtonNew("Melee", defaultFont)
	toolbar.container.AddChild(toolbar.mnuMelee)
	toolbar.mnuMelee.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			toolbarMenuOpen(args.Button.GetWidget(), ui, toolbar.btnMeleeHitOut, toolbar.btnMeleeHitIn, toolbar.btnMeleeCritOut, toolbar.btnMeleeCritIn, toolbar.btnMeleeMissOut, toolbar.btnMeleeMissIn)
		}))

	toolbar.mnuSpell = toolbarButtonNew("Spell", defaultFont)
	toolbar.container.AddChild(toolbar.mnuSpell)
	toolbar.mnuSpell.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			toolbarMenuOpen(args.Button.GetWidget(), ui, toolbar.btnSpellHitOut, toolbar.btnSpellHitIn, toolbar.btnSpellCritOut, toolbar.btnSpellCritIn, toolbar.btnSpellMissOut, toolbar.btnSpellMissIn)
		}))

	toolbar.mnuHeal = toolbarButtonNew("Heal", defaultFont)
	toolbar.container.AddChild(toolbar.mnuHeal)
	toolbar.mnuHeal.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			toolbarMenuOpen(args.Button.GetWidget(), ui, toolbar.btnHealHitOut, toolbar.btnHealHitIn, toolbar.btnHealCritOut, toolbar.btnHealCritIn)
		}))

	toolbar.mnuRune = toolbarButtonNew("Rune", defaultFont)
	toolbar.container.AddChild(toolbar.mnuRune)
	toolbar.mnuRune.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			toolbarMenuOpen(args.Button.GetWidget(), ui, toolbar.btnRuneHitOut, toolbar.btnRuneHitIn)
		}))
	toolbar.mnuTotal = toolbarButtonNew("Total", defaultFont)
	toolbar.container.AddChild(toolbar.mnuTotal)
	toolbar.mnuTotal.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			toolbarMenuOpen(args.Button.GetWidget(), ui, toolbar.btnTotalDamageOut, toolbar.btnTotalDamageIn)
		}))

	toolbar.mnuExtra = toolbarButtonNew("Extra", defaultFont)
	toolbar.container.AddChild(toolbar.mnuExtra)
	toolbar.mnuExtra.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			toolbarMenuOpen(args.Button.GetWidget(), ui, toolbar.btnMoney)
		}))

	toolbar.btnMoney = toolbarButtonNew("Money", defaultFont)
	toolbar.btnMoney.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			money.Toggle()
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { status.Set("Toggle Money") }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { status.Set("") }),
	)

	return toolbar, nil
}

func toolbarButtonNew(label string, face text.Face) *widget.Button {
	// Create a button for the toolbar.
	return widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.Text(label, face, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    4,
			Left:   4,
			Right:  32,
			Bottom: 4,
		}),
	)
}

func toolbarMenuEntryNew(label string, face text.Face) *widget.Button {
	// Create a button for a menu entry.
	return widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.Text(label, face, &widget.ButtonTextColor{
			Idle:     color.White,
			Disabled: colornames.Gray,
			Hover:    color.White,
			Pressed:  color.Black,
		}),
		widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		widget.ButtonOpts.TextPadding(widget.Insets{Left: 16, Right: 64}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
	)
}

func toolbarMenuOpen(opener *widget.Widget, ui *ebitenui.UI, entries ...*widget.Button) {
	c := widget.NewContainer(
		// Set the background to a translucent black.
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.RGBA{R: 0, G: 0, B: 0, A: 125})),

		// Menu entries should be arranged vertically.
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(4),
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 1, Bottom: 1}),
			),
		),

		// Set the minimum size for the menu.
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(64, 0)),
	)

	for _, entry := range entries {
		c.AddChild(entry)
	}

	w, h := c.PreferredSize()

	window := widget.NewWindow(
		// Set the menu to be a modal. This makes it block UI interactions to anything ese.
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),

		// Close the menu if the user clicks outside of it.
		widget.WindowOpts.CloseMode(widget.CLICK),

		// Position the menu below the menu button that it belongs to.
		widget.WindowOpts.Location(
			goimage.Rect(
				opener.Rect.Min.X,
				opener.Rect.Min.Y+opener.Rect.Max.Y,
				opener.Rect.Min.X+w,
				opener.Rect.Min.Y+opener.Rect.Max.Y+opener.Rect.Min.Y+h,
			),
		),
	)

	// Immediately add the menu to the UI.
	ui.AddWindow(window)
}

func toggleFullscreenBorderless(cfg *config.CritSprinklerConfiguration) {
	cfg.IsFullscreenBorderless = !cfg.IsFullscreenBorderless
	go func() {
		if cfg.IsFullscreenBorderless {
			ebiten.SetWindowDecorated(false)
			ebiten.SetWindowPosition(0, 0)
			ebiten.SetWindowSize(ebiten.Monitor().Size())
		} else {
			ebiten.SetWindowDecorated(true)
		}
	}()
}
