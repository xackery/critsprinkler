// toolbar.go
//
// Toolbar struct and related functions.
//

package main

import (
	goimage "image"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
)

// NOTE: It's not strictly necessary to store references to all the buttons in the toolbar struct, but this example does
// so for completeness' sake. When you keep a reference to buttons in the struct, you can later configure them to respond
// to certain events in your application, and keep your program's logic outside the toolbar.
type toolbar struct {
	container    *widget.Container
	mnuFile      *widget.Button
	btnLoadEQLog *widget.Button
	btnSave      *widget.Button
	btnQuit      *widget.Button

	mnuGlobal        *widget.Button
	btnGlobalHitOut  *widget.Button
	btnGlobalHitIn   *widget.Button
	btnGlobalCritOut *widget.Button
	btnGlobalCritIn  *widget.Button
	btnGlobalMissOut *widget.Button
	btnGlobalMissIn  *widget.Button

	mnuMelee        *widget.Button
	btnMeleeHitOut  *widget.Button
	btnMeleeHitIn   *widget.Button
	btnMeleeCritOut *widget.Button
	btnMeleeCritIn  *widget.Button
	btnMeleeMissOut *widget.Button
	btnMeleeMissIn  *widget.Button

	mnuSpell        *widget.Button
	btnSpellHitOut  *widget.Button
	btnSpellHitIn   *widget.Button
	btnSpellCritOut *widget.Button
	btnSpellCritIn  *widget.Button
	btnSpellMissOut *widget.Button
	btnSpellMissIn  *widget.Button

	mnuHeal        *widget.Button
	btnHealHitOut  *widget.Button
	btnHealHitIn   *widget.Button
	btnHealCritOut *widget.Button
	btnHealCritIn  *widget.Button

	mnuRune       *widget.Button
	btnRuneHitOut *widget.Button
	btnRuneHitIn  *widget.Button

	btnGlobalTest *widget.Button
}

func newToolbar(ui *ebitenui.UI, res *uiResources) *toolbar {
	// Create a root container for the toolbar.
	root := widget.NewContainer(
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

	//
	// "File" menu
	//
	file := newToolbarButton(res, "File")
	var (
		load = newToolbarMenuEntry(res, "Load EQ Log")
		save = newToolbarMenuEntry(res, "Save")
		quit = newToolbarMenuEntry(res, "Quit")
	)

	// Make the toolbar entry open a menu with our "save" and "load" entries  when the user clicks it.
	file.Configure(
		// Make the toolbar entry open a menu with our "save" and "load" entries  when the user clicks it.
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, load, save, quit)
		}),
	)
	root.AddChild(file)

	btnGlobal := newToolbarButton(res, "Global")
	var (
		btnGlobalCritOut = newToolbarMenuEntry(res, "Outgoing Crits")
		btnGlobalHitOut  = newToolbarMenuEntry(res, "Outgoing Hits")
		btnGlobalMissOut = newToolbarMenuEntry(res, "Outgoing Misses")
		btnGlobalCritIn  = newToolbarMenuEntry(res, "Incoming Crits")
		btnGlobalHitIn   = newToolbarMenuEntry(res, "Incoming Hits")
		btnGlobalMissIn  = newToolbarMenuEntry(res, "Incoming Misses")
	)
	btnGlobal.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, btnGlobalCritOut, btnGlobalHitOut, btnGlobalMissOut, btnGlobalCritIn, btnGlobalHitIn, btnGlobalMissIn)
		}),
	)
	root.AddChild(btnGlobal)

	btnMelee := newToolbarButton(res, "Melee")
	var (
		btnMeleeCritOut = newToolbarMenuEntry(res, "Outgoing Crits")
		btnMeleeHitOut  = newToolbarMenuEntry(res, "Outgoing Hits")
		btnMeleeMissOut = newToolbarMenuEntry(res, "Outgoing Misses")
		btnMeleeCritIn  = newToolbarMenuEntry(res, "Incoming Crits")
		btnMeleeHitIn   = newToolbarMenuEntry(res, "Incoming Hits")
		btnMeleeMissIn  = newToolbarMenuEntry(res, "Incoming Misses")
	)
	btnMelee.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, btnMeleeCritOut, btnMeleeHitOut, btnMeleeMissOut, btnMeleeCritIn, btnMeleeHitIn, btnMeleeMissIn)
		}),
	)
	root.AddChild(btnMelee)

	btnSpell := newToolbarButton(res, "Spell")
	var (
		btnSpellCritOut = newToolbarMenuEntry(res, "Outgoing Crits")
		btnSpellHitOut  = newToolbarMenuEntry(res, "Outgoing Hits")
		btnSpellMissOut = newToolbarMenuEntry(res, "Outgoing Resists")
		btnSpellCritIn  = newToolbarMenuEntry(res, "Incoming Crits")
		btnSpellHitIn   = newToolbarMenuEntry(res, "Incoming Hits")
		btnSpellMissIn  = newToolbarMenuEntry(res, "Incoming Resists")
	)

	btnSpell.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, btnSpellCritOut, btnSpellHitOut, btnSpellMissOut, btnSpellCritIn, btnSpellHitIn, btnSpellMissIn)
		}),
	)
	root.AddChild(btnSpell)

	btnHeal := newToolbarButton(res, "Heal")
	var (
		btnHealCritOut = newToolbarMenuEntry(res, "Outgoing Crits")
		btnHealHitOut  = newToolbarMenuEntry(res, "Outgoing Hits")
		btnHealCritIn  = newToolbarMenuEntry(res, "Incoming Crits")
		btnHealHitIn   = newToolbarMenuEntry(res, "Incoming Hits")
	)

	btnHeal.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, btnHealCritOut, btnHealHitOut, btnHealCritIn, btnHealHitIn)
		}),
	)
	root.AddChild(btnHeal)

	btnRune := newToolbarButton(res, "Rune")
	var (
		btnRuneHitOut = newToolbarMenuEntry(res, "Outgoing Hits")
		btnRuneHitIn  = newToolbarMenuEntry(res, "Incoming Hits")
	)

	btnRune.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openToolbarMenu(args.Button.GetWidget(), ui, btnRuneHitOut, btnRuneHitIn)
		}),
	)
	root.AddChild(btnRune)

	//
	// "Help" button
	// Unlike the "File" and "Edit" menu, this is just a regular button on the toolbar - it does not open a menu.
	// You can configure it to do something else when it's pressed, like opening a "Help" window.
	//
	btnTest := newToolbarButton(res, "Test")
	root.AddChild(btnTest)

	return &toolbar{
		container:     root,
		mnuFile:       file,
		btnGlobalTest: btnTest,
		btnSave:       save,
		btnLoadEQLog:  load,
		btnQuit:       quit,

		mnuGlobal:        btnGlobal,
		btnGlobalCritOut: btnGlobalCritOut,
		btnGlobalCritIn:  btnGlobalCritIn,
		btnGlobalHitOut:  btnGlobalHitOut,
		btnGlobalHitIn:   btnGlobalHitIn,
		btnGlobalMissOut: btnGlobalMissOut,
		btnGlobalMissIn:  btnGlobalMissIn,

		mnuMelee:        btnMelee,
		btnMeleeHitOut:  btnMeleeHitOut,
		btnMeleeHitIn:   btnMeleeHitIn,
		btnMeleeCritOut: btnMeleeCritOut,
		btnMeleeCritIn:  btnMeleeCritIn,
		btnMeleeMissOut: btnMeleeMissOut,
		btnMeleeMissIn:  btnMeleeMissIn,

		mnuSpell:        btnSpell,
		btnSpellHitOut:  btnSpellHitOut,
		btnSpellHitIn:   btnSpellHitIn,
		btnSpellCritOut: btnSpellCritOut,
		btnSpellCritIn:  btnSpellCritIn,
		btnSpellMissOut: btnSpellMissOut,
		btnSpellMissIn:  btnSpellMissIn,

		mnuHeal:        btnHeal,
		btnHealHitOut:  btnHealHitOut,
		btnHealHitIn:   btnHealHitIn,
		btnHealCritOut: btnHealCritOut,
		btnHealCritIn:  btnHealCritIn,

		mnuRune:       btnRune,
		btnRuneHitOut: btnRuneHitOut,
		btnRuneHitIn:  btnRuneHitIn,
	}
}

func newToolbarButton(res *uiResources, label string) *widget.Button {
	// Create a button for the toolbar.
	return widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.Text(label, res.button.face, &widget.ButtonTextColor{
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

func newToolbarMenuEntry(res *uiResources, label string) *widget.Button {
	// Create a button for a menu entry.
	return widget.NewButton(
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.Transparent),
			Hover:   image.NewNineSliceColor(colornames.Darkgray),
			Pressed: image.NewNineSliceColor(colornames.White),
		}),
		widget.ButtonOpts.Text(label, res.button.face, &widget.ButtonTextColor{
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

func openToolbarMenu(opener *widget.Widget, ui *ebitenui.UI, entries ...*widget.Button) {
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
