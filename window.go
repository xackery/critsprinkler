package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"strings"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/popup"
	"github.com/xackery/critsprinkler/tracker"
	"github.com/xackery/encdec"
)

/*
var (
	custColorPalette = [16]win.COLORREF{
		0x00FFFFFF, 0x00FFFFFF, 0x00FFFFFF, 0x00FFFFFF,
		0x00FFFFFF, 0x00FFFFFF, 0x00FFFFFF, 0x00FFFFFF,
		0x00FFFFFF, 0x00FFFFFF, 0x00FFFFFF, 0x00FFFFFF,
		0x00FFFFFF, 0x00FFFFFF, 0x00FFFFFF, 0x00FFFFFF,
	}
) */

func openPopupSettingsWindow(res *uiResources, ui *ebitenui.UI, setting *popup.SettingProperty) error {
	var rw widget.RemoveWindowFunc
	var window *widget.Window

	titleBar := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.titleBar),
		widget.ContainerOpts.Layout(widget.NewGridLayout(widget.GridLayoutOpts.Columns(3), widget.GridLayoutOpts.Stretch([]bool{true, true, false}, []bool{true}), widget.GridLayoutOpts.Padding(widget.Insets{
			Left:   10,
			Right:  5,
			Top:    6,
			Bottom: 5,
		}))))

	titleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Test", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if popup.IsGlobalCategory(setting.Category) {
				categories := popup.GlobalCategoryToCategory(setting.Category)
				for _, category := range categories {
					onDamageEvent(&common.DamageEvent{
						Category:  category,
						Source:    tracker.PlayerName(),
						Target:    "Test",
						SpellName: "Ice Comet",
						Damage:    category.String(),
					})
				}

				return
			}

			onDamageEvent(&common.DamageEvent{
				Category:  setting.Category,
				Source:    tracker.PlayerName(),
				Target:    "Test",
				SpellName: "Ice Comet",
				Damage:    fmt.Sprintf("%d", rand.Int()%1000),
			})
		}),
		widget.ButtonOpts.TabOrder(99),
	))

	sc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}))),
	)

	titleBar.AddChild(sc)

	titleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("X", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
		widget.ButtonOpts.TabOrder(99),
	))

	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, false}),
				widget.GridLayoutOpts.Padding(res.panel.padding),
				widget.GridLayoutOpts.Spacing(0, 15),
			),
		),
	)

	c.AddChild(widget.NewText(
		widget.TextOpts.Text(setting.Title, res.text.face, res.text.idleColor),
	))

	/*
		tOpts := []widget.TextInputOpt{
			widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
			widget.TextInputOpts.Image(res.textInput.image),
			widget.TextInputOpts.Color(res.textInput.color),
			widget.TextInputOpts.Padding(widget.Insets{
				Left:   13,
				Right:  13,
				Top:    7,
				Bottom: 7,
			}),
			widget.TextInputOpts.Face(res.textInput.face),
			widget.TextInputOpts.CaretOpts(
				widget.CaretOpts.Size(res.textInput.face, 2),
			),
		}
			t := widget.NewTextInput(append(
				tOpts,
				widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					StretchHorizontal: true,
				})),
				widget.TextInputOpts.Placeholder("Enter text here"))...,
			)
			textContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
			textContainer.AddChild(t)
			c.AddChild(textContainer) */

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(15),
			),
			// widget.NewGridLayout(
			// 	widget.GridLayoutOpts.Columns(2),
			// 	widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{true}),
			// 	widget.GridLayoutOpts.Padding(res.panel.padding),
			// 	widget.GridLayoutOpts.Spacing(0, 15),
			// ),
		),
	)
	c.AddChild(bc)

	if !popup.IsGlobalCategory(setting.Category) {
		chkEnable := newCheckbox("Enabled", func(args *widget.CheckboxChangedEventArgs) {
			*setting.IsEnabled = args.State == widget.WidgetChecked
		}, res)
		state := 0
		if *setting.IsEnabled {
			state = 1
		}
		chkEnable.SetState(widget.WidgetState(state))

		c.AddChild(chkEnable)
	}
	/* cb := widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Color", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			 fmt.Println("hwnd:")
			lpcc := &win.CHOOSECOLOR{
				LpCustColors: &custColorPalette,
				HwndOwner:    game.hwnd,
				Flags:        win.CC_RGBINIT | win.CC_FULLOPEN, // | win.CC_ANYCOLOR,
				RgbResult:    0,
			}
			lpcc.LStructSize = uint32(unsafe.Sizeof(*lpcc))
			if win.ChooseColor(lpcc) {
				rgb := lpcc.RgbResult
				r := uint8(rgb & 0xFF)         // Extract bits 0-7 for Red
				g := uint8((rgb >> 8) & 0xFF)  // Extract bits 8-15 for Green
				b := uint8((rgb >> 16) & 0xFF) // Extract bits 16-23 for Blue
				fmt.Printf("Color: %d %d %d\n", r, g, b)
			}
		}),
	)
	bc.AddChild(cb) */

	directionLabel := widget.NewText(
		widget.TextOpts.Text("Direction", res.text.face, res.text.idleColor),
	)

	cmmDirection := newListComboButton(
		[]interface{}{0, 1, 2, 3, 4, 5, 6, 7},
		func(e interface{}) string {
			dir := common.Direction(e.(int))
			return dir.String()
		},
		func(e interface{}) string {
			dir := common.Direction(e.(int))
			return dir.String()
		},
		func(args *widget.ListComboButtonEntrySelectedEventArgs) {
			c.RequestRelayout()
			popup.SetSettingDirection(setting.Category, common.Direction(args.Entry.(int)))
		},
		res)
	cmmDirection.SetSelectedEntry(int(*setting.Direction))
	newLabelProperty(c, directionLabel, cmmDirection)

	colorResultLabel := widget.NewText(
		widget.TextOpts.Text("Result", res.text.face, res.text.idleColor),
	)

	colorResult := widget.NewText(
		widget.TextOpts.Text("1234", res.fonts.popupBold36, res.text.idleColor),
	)

	colorLabel := widget.NewText(
		widget.TextOpts.Text("Hex", res.text.face, res.text.idleColor),
	)

	colorInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:   false,
			MaxWidth:  100,
			MaxHeight: 30,
		})),
		widget.TextInputOpts.Image(res.textInput.image),
		widget.TextInputOpts.Color(res.textInput.color),
		widget.TextInputOpts.Padding(widget.Insets{
			Left:   13,
			Right:  13,
			Top:    7,
			Bottom: 7,
		}),
		widget.TextInputOpts.Face(res.textInput.face),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(res.textInput.face, 2),
		),
		widget.TextInputOpts.Placeholder("FF00FFFF"),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			hexValue := args.InputText
			if len(hexValue) != 8 {
				return
			}
			dec := encdec.NewDecoder(bytes.NewReader([]byte(hexValue)), binary.LittleEndian)
			rStr := strings.ToUpper(dec.StringFixed(2))
			gStr := strings.ToUpper(dec.StringFixed(2))
			bStr := strings.ToUpper(dec.StringFixed(2))
			aStr := strings.ToUpper(dec.StringFixed(2))

			// convert FF to 255

			rgba := color.RGBA{}

			fmt.Sscanf(rStr, "%X", &rgba.R)
			fmt.Sscanf(gStr, "%X", &rgba.G)
			fmt.Sscanf(bStr, "%X", &rgba.B)
			fmt.Sscanf(aStr, "%X", &rgba.A)

			popup.SetSettingColorByCategory(setting.Category, rgba)
			colorResult.Color = rgba
		}),
	)

	colorInput.SetText(fmt.Sprintf("%02X%02X%02X%02X", setting.Color.R, setting.Color.G, setting.Color.B, setting.Color.A))

	if !popup.IsGlobalCategory(setting.Category) {
		newLabelProperty(c, colorLabel, colorInput)
		newLabelProperty(c, colorResultLabel, colorResult)
	}

	o2bLabel := widget.NewText(
		widget.TextOpts.Text("Location", res.text.face, res.text.idleColor),
	)

	o2b := widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Set", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			err := openAdjustmentWindow(res, ui, setting)
			if err != nil {
				fmt.Println("Error opening adjustment window: ", err)
				game.statusText = fmt.Sprintf("Error opening adjustment window: %v", err)
			}
		}),
	)
	newLabelProperty(c, o2bLabel, o2b)
	//bc.AddChild(o2b)
	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.TitleBar(titleBar, 30),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(250, 80),
		widget.WindowOpts.MaxSize(700, 600),
	)
	windowSize := input.GetWindowSize()
	r := image.Rect(0, 0, 350, 550)
	r = r.Add(image.Point{windowSize.X / 4 / 2, windowSize.Y * 2 / 3 / 5})

	window.SetLocation(r)

	rw = ui.AddWindow(window)
	return nil
}

func openAdjustmentWindow(res *uiResources, ui *ebitenui.UI, setting *popup.SettingProperty) error {
	var rw widget.RemoveWindowFunc
	c := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(res.panel.padding),
			widget.RowLayoutOpts.Spacing(15),
		)),
	)
	titleBar := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.titleBar),
		widget.ContainerOpts.Layout(widget.NewGridLayout(widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{true, true, false}, []bool{true}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Left:   10,
				Right:  5,
				Top:    6,
				Bottom: 5,
			}))))
	titleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("Test", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {

			if popup.IsGlobalCategory(setting.Category) {
				categories := popup.GlobalCategoryToCategory(setting.Category)
				for _, category := range categories {
					onDamageEvent(&common.DamageEvent{
						Category:  category,
						Source:    tracker.PlayerName(),
						Target:    "Test",
						SpellName: "Ice Comet",
						Damage:    fmt.Sprintf("%d", rand.Int()%1000),
					})
				}

				return
			}

			onDamageEvent(&common.DamageEvent{
				Category:  setting.Category,
				Source:    tracker.PlayerName(),
				Target:    "Test",
				SpellName: "Ice Comet",
				Damage:    fmt.Sprintf("%d", rand.Int()%1000),
			})
		}),
		widget.ButtonOpts.TabOrder(99),
	))

	sc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}))),
	)

	titleBar.AddChild(sc)

	titleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(res.button.image),
		widget.ButtonOpts.TextPadding(res.button.padding),
		widget.ButtonOpts.Text("X", res.button.face, res.button.text),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
		widget.ButtonOpts.TabOrder(99),
	))
	w := widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.TitleBar(titleBar, 30),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(70, 70),
		//widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			popup.SetSettingPositionByCategory(setting.Category, args.Rect)
		}),
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			popup.SetSettingPositionByCategory(setting.Category, args.Rect)
		}),
	)
	rect, err := popup.SettingPositionByCategory(setting.Category)
	if err != nil {
		return fmt.Errorf("setting position by category: %w", err)
	}
	w.SetLocation(rect)

	rw = ui.AddWindow(w)
	return nil
}

/*
func newSeparator(res *uiResources, ld interface{}) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(ld)))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.ImageNineSlice(eimage.NewNineSliceColor(res.separatorColor)),
	))

	return c
} */

func newLabelProperty(c *widget.Container, label widget.PreferredSizeLocateableWidget, property widget.PreferredSizeLocateableWidget) widget.RemoveChildFunc {
	pc := widget.NewContainer(

		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch([]bool{false, true}, nil),
			// widget.GridLayoutOpts.Padding(widget.Insets{
			// 	Left:   10,
			// 	Right:  10,
			// 	Top:    10,
			// 	Bottom: 10,
			// }),

			widget.GridLayoutOpts.Spacing(10, 0),
		)),
	)

	pc.AddChild(label)
	pc.AddChild(property)

	c.AddChild(pc)
	return func() {
		pc.RemoveChild(label)
		pc.RemoveChild(property)
		c.RemoveChild(pc)
	}
}

func newListComboButton(entries []interface{}, buttonLabel widget.SelectComboButtonEntryLabelFunc, entryLabel widget.ListEntryLabelFunc,
	entrySelectedHandler widget.ListComboButtonEntrySelectedHandlerFunc, res *uiResources) *widget.ListComboButton {

	return widget.NewListComboButton(
		widget.ListComboButtonOpts.SelectComboButtonOpts(
			widget.SelectComboButtonOpts.ComboButtonOpts(
				widget.ComboButtonOpts.ButtonOpts(
					widget.ButtonOpts.Image(res.comboButton.image),
					widget.ButtonOpts.TextPadding(res.comboButton.padding),
				),
			),
		),
		widget.ListComboButtonOpts.Text(res.comboButton.face, res.comboButton.graphic, res.comboButton.text),
		widget.ListComboButtonOpts.ListOpts(
			widget.ListOpts.Entries(entries),
			widget.ListOpts.ScrollContainerOpts(
				widget.ScrollContainerOpts.Image(res.list.image),
			),
			widget.ListOpts.SliderOpts(
				widget.SliderOpts.Images(res.list.track, res.list.handle),
				widget.SliderOpts.MinHandleSize(res.list.handleSize),
				widget.SliderOpts.TrackPadding(res.list.trackPadding)),
			widget.ListOpts.EntryFontFace(res.list.face),
			widget.ListOpts.EntryColor(res.list.entry),
			widget.ListOpts.EntryTextPadding(res.list.entryPadding),
		),
		widget.ListComboButtonOpts.EntryLabelFunc(buttonLabel, entryLabel),
		widget.ListComboButtonOpts.EntrySelectedHandler(entrySelectedHandler))
}

func newCheckbox(label string, changedHandler widget.CheckboxChangedHandlerFunc, res *uiResources) *widget.LabeledCheckbox {
	return widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(res.checkbox.spacing),
		widget.LabeledCheckboxOpts.LabelFirst(),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(res.checkbox.image)),
			widget.CheckboxOpts.Image(res.checkbox.graphic),
			widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
				if changedHandler != nil {
					changedHandler(args)
				}
			})),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text(label, res.label.face, res.label.text)))

}
