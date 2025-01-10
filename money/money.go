package money

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/config"
	"github.com/xackery/critsprinkler/library"
	"github.com/xackery/critsprinkler/sound"
	"github.com/xackery/critsprinkler/tracker"
	"github.com/xackery/critsprinkler/util"
	"golang.org/x/exp/rand"
)

const (
	gravity        = 0.3
	bounceFactor   = 0.8
	horizontalDrag = 0.98
)

var (
	ui             *ebitenui.UI
	cfg            *config.CritSprinklerConfiguration
	placement      *common.Placement
	panelContainer *widget.Container
	face           text.Face

	platinumImage     *ebiten.Image
	platinumCollected int
	platinumBuffer    int

	goldImage     *ebiten.Image
	goldCollected int
	goldBuffer    int

	silverImage     *ebiten.Image
	silverCollected int
	silverBuffer    int

	copperImage     *ebiten.Image
	copperCollected int
	copperBuffer    int
	upgradeTimer    time.Time

	favorImage   *ebiten.Image
	sprinkleChan = make(chan showerEvent, 1000)
	sprinkles    []*sprinkle
)

type showerEvent struct {
	Currency library.Misc
	Amount   int
}

type sprinkle struct {
	Currency       library.Misc
	Amount         int
	IsTallyEnabled bool
	image          *ebiten.Image
	maxLife        int
	life           int
	x, y           float64
	vx, vy         float64
	fade           float32
}

func New(eui *ebitenui.UI, ecfg *config.CritSprinklerConfiguration) error {
	cfg = ecfg
	placement = &cfg.Money
	ui = eui
	err := tracker.Subscribe(onLine)
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	return Open()
}

func SetEditMode(ui *ebitenui.UI, editMode bool) {
	var err error
	state := widget.Visibility_Show
	if !editMode {
		state = widget.Visibility_Hide
		panelContainer.BackgroundImage = nil
	} else {
		panelContainer.BackgroundImage, err = library.NinesliceByKey(library.NineSlicePanelIdle)
		if err != nil {
			fmt.Println("ninesliceByKey", err)
		}
	}
	placement.TitleBar.GetWidget().Visibility = state
}

func Open() error {
	var err error

	face, err = library.FontByKey(library.FontSmall)
	if err != nil {
		return fmt.Errorf("fontByKey: %w", err)
	}

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
		widget.ContainerOpts.Layout(widget.NewGridLayout(widget.GridLayoutOpts.Columns(2), widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true})))) // widget.GridLayoutOpts.Padding(widget.Insets{

	// 	Left:   10,
	// 	Right:  5,
	// 	Top:    6,
	// 	Bottom: 5,
	// })	),
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

			sound.Play(sound.SoundEffectBuyItem)
			sprinkleChan <- showerEvent{library.MiscSilver, 10}
		}),
		widget.ButtonOpts.TabOrder(99),
	))
	placement.TitleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(buttonInvisibleImage),
		//widget.ButtonOpts.TextPadding(widget.Insets{Left: 30, Right: 30}),
		widget.ButtonOpts.Text("Money", face, &widget.ButtonTextColor{
			Idle:     util.HexToColor("dFF4FFFF"),
			Disabled: util.HexToColor("5A7A91FF"),
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {

		}),
		widget.ButtonOpts.TabOrder(99),
	))
	placement.TitleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(buttonCloseImage),
		widget.ButtonOpts.TextPadding(widget.Insets{Left: 30, Right: 30}),
		//widget.ButtonOpts.TextPadding(widget.Insets{Left: 30, Right: 30}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			Close()
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

	panelContainer = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(panelNineSlice),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.Insets{Left: 10, Right: 10, Top: 10, Bottom: 10}),
			),
		),
	)

	placement.Window = widget.NewWindow(
		//widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(panelContainer),
		widget.WindowOpts.TitleBar(placement.TitleBar, 30),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(100, 80),
		//widget.WindowOpts.MaxSize(300, 1000),
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			*placement.WindowRect = args.Rect
			OnResize()
		}),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			*placement.WindowRect = args.Rect
			OnResize()
		}),
	)
	//windowSize := input.GetWindowSize()

	//r := image.Rect(500, 50, 100, 100)
	//r = r.Add(image.Point{windowSize.X / 4 / 2, windowSize.Y * 2 / 3 / 5})

	//window.SetLocation(r)
	placement.Window.SetLocation(*placement.WindowRect)

	placement.IsVisible = 1
	_ = ui.AddWindow(placement.Window)
	panelContainer.RequestRelayout()
	return nil
}

func Close() error {
	if placement.Window == nil {
		return nil
	}
	placement.Window.Close()
	placement.Window = nil
	placement.IsVisible = 0
	return nil
}

// Toggle opens or closes the window
func Toggle() error {
	if placement.IsVisible == 1 {
		return Close()
	}
	return Open()
}

func Update() {
	for {
		isFound := false
		select {
		case sprinkle := <-sprinkleChan:
			fmt.Println("sprinkle", sprinkle)
			for i := 0; i < sprinkle.Amount; i++ {
				sprinkleOut(sprinkle.Currency, 1)
			}
			isFound = true
		default:
		}
		if !isFound {
			break
		}
	}

	for i := len(sprinkles) - 1; i >= 0; i-- {
		sprinkle := sprinkles[i]

		sprinkle.life--
		sprinkle.vy += gravity
		sprinkle.y += sprinkle.vy
		sprinkle.x += sprinkle.vx
		sprinkle.vx *= horizontalDrag

		ground := placement.WindowRect.Max.Y - placement.WindowRect.Min.Y - 20

		if sprinkle.y > float64(ground) {
			sprinkle.y = float64(ground)

			sprinkle.vy = -sprinkle.vy * bounceFactor

			if math.Abs(sprinkle.vy) < 1 {
				sprinkle.vy = 0
			} else {

				//sound.PlayBounceRandom()
			}
		}
		if sprinkle.vy == 0 {
			sprinkle.fade -= 0.01
		}

		// if we are at the target, remove the sprinkle
		// if int(sprinkle.x) == *sprinkle.baseX && int(sprinkle.y) == *sprinkle.baseY {
		// 	switch sprinkle.Currency {
		// 	case library.MiscPlatinum:
		// 		platinumBuffer += sprinkle.Amount
		// 	case library.MiscGold:
		// 		goldBuffer += sprinkle.Amount
		// 	case library.MiscSilver:
		// 		silverBuffer += sprinkle.Amount
		// 	case library.MiscCopper:
		// 		copperBuffer += sprinkle.Amount
		// 	}

		// 	sprinkles = append(sprinkles[:i], sprinkles[i+1:]...)
		//} else if sprinkle.life <= 0 {
		if sprinkle.fade <= 0 {
			switch sprinkle.Currency {
			case library.MiscPlatinum:
				platinumBuffer += sprinkle.Amount
			case library.MiscGold:
				goldBuffer += sprinkle.Amount
			case library.MiscSilver:
				silverBuffer += sprinkle.Amount
			case library.MiscCopper:
				copperBuffer += sprinkle.Amount
			}
			upgradeTimer = time.Now().Add(time.Second * 3)
			sprinkles = append(sprinkles[:i], sprinkles[i+1:]...)

		}
		//}

	}

	platinumCollected = bufferApply(platinumCollected, platinumBuffer)
	goldCollected = bufferApply(goldCollected, goldBuffer)
	silverCollected = bufferApply(silverCollected, silverBuffer)
	copperCollected = bufferApply(copperCollected, copperBuffer)
	if upgradeTimer.Before(time.Now()) {
		for copperBuffer >= 10 {
			copperBuffer -= 10
			silverBuffer++
		}
		for silverBuffer >= 10 {
			silverBuffer -= 10
			goldBuffer++
		}
		for goldBuffer >= 10 {
			goldBuffer -= 10
			platinumBuffer++
		}

		upgradeTimer = time.Now().Add(time.Second * 3)
	}

}

func OnResize() {
	w, h := ebiten.WindowSize()

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

	// Apply the changes only if the rectangle has changed
	if newRect != originalRect {
		fmt.Println("going from", originalRect, "to", newRect)
		placement.Window.SetLocation(newRect)
		*placement.WindowRect = newRect
	}
}

func Draw(screen *ebiten.Image) {
	if placement.IsVisible == 0 {
		return
	}

	windowX, windowY := placement.WindowRect.Min.X, placement.WindowRect.Min.Y
	//windowW := placement.WindowRect.Dx()
	windowH := placement.WindowRect.Dy()

	x := float64(windowX + 10)
	y := float64(windowY + windowH)
	dOp := &ebiten.DrawImageOptions{}
	tOp := &text.DrawOptions{}
	txt := ""
	if platinumImage != nil {
		dOp.GeoM.Translate(x, y)
		screen.DrawImage(platinumImage, dOp)
		dOp.GeoM.Reset()
	}
	x += 20
	tOp.GeoM.Translate(x, y)
	txt = fmt.Sprintf("%d", platinumCollected)
	tW, _ := text.Measure(txt, face, 0)
	text.Draw(screen, txt, face, tOp)
	tOp.GeoM.Reset()
	x += tW
	if goldImage != nil {
		dOp.GeoM.Translate(x, y)
		screen.DrawImage(goldImage, dOp)
		dOp.GeoM.Reset()
	}
	x += 20
	tOp.GeoM.Translate(x, y)
	txt = fmt.Sprintf("%d", goldCollected)
	tW, _ = text.Measure(txt, face, 0)
	text.Draw(screen, txt, face, tOp)
	tOp.GeoM.Reset()
	x += tW
	if silverImage != nil {
		dOp.GeoM.Translate(x, y)
		screen.DrawImage(silverImage, dOp)
		dOp.GeoM.Reset()
	}
	x += 20
	tOp.GeoM.Translate(x, y)
	txt = fmt.Sprintf("%d", silverCollected)
	tW, _ = text.Measure(txt, face, 0)
	text.Draw(screen, txt, face, tOp)
	tOp.GeoM.Reset()
	x += tW
	if copperImage != nil {
		dOp.GeoM.Translate(x, y)
		screen.DrawImage(copperImage, dOp)
		dOp.GeoM.Reset()
	}
	x += 20
	tOp.GeoM.Translate(x, y)
	txt = fmt.Sprintf("%d", copperCollected)
	text.Draw(screen, txt, face, tOp)
	tOp.GeoM.Reset()

	for _, sprinkle := range sprinkles {

		if sprinkle.image == nil {
			op := &text.DrawOptions{}
			op.GeoM.Translate(sprinkle.x, sprinkle.y)
			op.GeoM.Translate(float64(placement.WindowRect.Min.X), float64(placement.WindowRect.Min.Y))

			rgba := color.RGBA{255, 255, 255, 255}
			msg := fmt.Sprintf("%d", sprinkle.Amount)
			switch sprinkle.Currency {
			case library.MiscPlatinum:
				rgba = color.RGBA{255, 255, 255, 255}
				msg = "p"
			case library.MiscGold:
				rgba = color.RGBA{255, 215, 0, 255}
				msg = "g"
			case library.MiscSilver:
				rgba = color.RGBA{192, 192, 192, 255}
				msg = "s"
			case library.MiscCopper:
				rgba = color.RGBA{205, 127, 50, 255}
				msg = "c"
			case library.MiscFavor:
				rgba = color.RGBA{255, 255, 255, 255}
				msg = "f"
			}
			op.ColorScale.ScaleWithColor(rgba)
			text.Draw(screen, msg, face, op)
		} else {
			op := &ebiten.DrawImageOptions{}
			subImage := sprinkle.image
			if sprinkle.fade < 1 {
				//op.ColorScale.ScaleAlpha(sprinkle.fade * 4)
				meltHeight := int(float32(sprinkle.image.Bounds().Dy()) * sprinkle.fade)

				subImage = sprinkle.image.SubImage(image.Rect(0, 0, sprinkle.image.Bounds().Dx(), meltHeight)).(*ebiten.Image)
				op.GeoM.Translate(0, -float64(meltHeight)+18)
			}

			op.GeoM.Translate(sprinkle.x, sprinkle.y)
			op.GeoM.Translate(float64(placement.WindowRect.Min.X), float64(placement.WindowRect.Min.Y))

			screen.DrawImage(subImage, op)
		}

	}
}

func OnEQPathLoad() {
	platinumImage = library.MiscByID(library.MiscPlatinum)
	goldImage = library.MiscByID(library.MiscGold)
	silverImage = library.MiscByID(library.MiscSilver)
	copperImage = library.MiscByID(library.MiscCopper)
	favorImage = library.MiscByID(library.MiscFavor)
}

func onLine(event time.Time, line string) {
	if parseCorpse(line) {
		return
	}
	if parseCorpseSplit(line) {
		return
	}
	if parseMerchant(line) {
		return
	}
	if parseTribute(line) {
		return
	}
}

func parseCorpse(line string) bool {
	match, ok := util.Parse(line, `\] You receive (.*) from the corpse.`, 1)
	if !ok {
		return false
	}
	var val int
	records := strings.Split(match[0], " ")
	for _, record := range records {
		newVal, err := strconv.Atoi(record)
		if err == nil {
			val = newVal
			continue
		}
		record = strings.ReplaceAll(record, "and", "")
		record = strings.TrimSpace(strings.ReplaceAll(record, ",", ""))
		if len(record) < 1 {
			continue
		}
		fmt.Println(record)

		sound.Play(sound.SoundEffectBuyItem)
		switch record {
		case "platinum":
			sprinkleChan <- showerEvent{library.MiscPlatinum, val}
		case "gold":
			sprinkleChan <- showerEvent{library.MiscGold, val}
		case "silver":
			sprinkleChan <- showerEvent{library.MiscSilver, val}
		case "copper":
			sprinkleChan <- showerEvent{library.MiscCopper, val}
		default:
			fmt.Println("unknown money type", record)
		}
	}
	return true
}

func parseCorpseSplit(line string) bool {
	match, ok := util.Parse(line, `\] You receive (.*) as your split.`, 1)
	if !ok {
		return false
	}
	var val int
	records := strings.Split(match[0], " ")
	for _, record := range records {
		newVal, err := strconv.Atoi(record)
		if err == nil {
			val = newVal
			continue
		}
		record = strings.ReplaceAll(record, "and", "")
		record = strings.TrimSpace(strings.ReplaceAll(record, ",", ""))
		if len(record) < 1 {
			continue
		}
		fmt.Println(record)

		sound.Play(sound.SoundEffectBuyItem)
		switch record {
		case "platinum":
			sprinkleChan <- showerEvent{library.MiscPlatinum, val}
		case "gold":
			sprinkleChan <- showerEvent{library.MiscGold, val}
		case "silver":
			sprinkleChan <- showerEvent{library.MiscSilver, val}
		case "copper":
			sprinkleChan <- showerEvent{library.MiscCopper, val}
		case "favor":
			sprinkleChan <- showerEvent{library.MiscFavor, val}
		default:
			fmt.Println("unknown money type", record)
		}
	}
	return true
}

func parseMerchant(line string) bool {
	match, ok := util.Parse(line, `\] You receive (.*) from .* for .*.`, 1)
	if !ok {
		return false
	}
	var val int
	records := strings.Split(match[0], " ")
	for _, record := range records {
		newVal, err := strconv.Atoi(record)
		if err == nil {
			val = newVal
			continue
		}
		record = strings.ReplaceAll(record, "and", "")
		record = strings.TrimSpace(strings.ReplaceAll(record, ",", ""))
		if len(record) < 1 {
			continue
		}
		fmt.Println(record)

		sound.Play(sound.SoundEffectBuyItem)

		switch record {
		case "platinum":
			sprinkleChan <- showerEvent{library.MiscPlatinum, val}
		case "gold":
			sprinkleChan <- showerEvent{library.MiscGold, val}
		case "silver":
			sprinkleChan <- showerEvent{library.MiscSilver, val}
		case "copper":
			sprinkleChan <- showerEvent{library.MiscCopper, val}
		default:
			fmt.Println("unknown money type", record)
		}
	}
	return true
}

func parseTribute(line string) bool {
	match, ok := util.Parse(line, `\] You have received (.*) favor for your tribute!`, 1)
	if !ok {
		return false
	}
	var val int
	newVal, err := strconv.Atoi(match[0])
	if err == nil {
		val = newVal
	}
	sound.Play(sound.SoundEffectBuyItem)
	sprinkleChan <- showerEvent{library.MiscFavor, val}
	return true
}

func sprinkleOut(currency library.Misc, val int) {
	if len(sprinkles) >= 100000 {
		switch currency {
		case library.MiscPlatinum:
			platinumBuffer += val
		case library.MiscGold:
			goldBuffer += val
		case library.MiscSilver:
			silverBuffer += val
		case library.MiscCopper:
			copperBuffer += val
		}
		return
	}
	vx := float64(0)
	vy := float64(-0.5 - rand.Float64() - (rand.Float64() / 2))
	switch placement.Direction {
	case common.DirectionUp:
		vy = -0.5 - rand.Float64() - (rand.Float64() / 2)
		vx = 0
	case common.DirectionDown:
		vy = 0.5 + rand.Float64() + (rand.Float64() / 2)
		vx = 0
	case common.DirectionLeft:
		vy = 0
		vx = -0.5 - rand.Float64() - (rand.Float64() / 2)
	case common.DirectionRight:
		vy = 0
		vx = 0.5 + rand.Float64() + (rand.Float64() / 2)
	case common.DirectionUpLeft:
		vy = -0.5 - rand.Float64() - (rand.Float64() / 2)
		vx = -0.5 - rand.Float64() - (rand.Float64() / 2)
	case common.DirectionUpRight:
		vy = -0.5 - rand.Float64() - (rand.Float64() / 2)
		vx = 0.5 + rand.Float64() + (rand.Float64() / 2)
	case common.DirectionDownLeft:
		vy = 0.5 + rand.Float64() + (rand.Float64() / 2)
		vx = -0.5 - rand.Float64() - (rand.Float64() / 2)
	case common.DirectionDownRight:
		vy = 0.5 + rand.Float64() + (rand.Float64() / 2)
		vx = 0.5 + rand.Float64() + (rand.Float64() / 2)
	}

	x := randomSpawnRange(int(placement.LastSpawnX), 0, placement.WindowRect.Dx()-50, placement.WindowRect.Dx()/4, 0)
	y := randomSpawnRange(int(placement.LastSpawnY), 0, placement.WindowRect.Dy()-50, placement.WindowRect.Dy()/4, 0)

	sprinkles = append(sprinkles, &sprinkle{
		Currency:       currency,
		Amount:         val,
		IsTallyEnabled: placement.IsTallyEnabled == 1,
		image:          library.MiscByID(currency),
		x:              x,
		y:              y,
		vx:             vx,
		vy:             vy,
		maxLife:        600,
		life:           600,
		fade:           1,
	})

	placement.LastSpawnX = int(x)
	placement.LastSpawnY = int(y)

}

func randomSpawnRange(lastSpawnX, minPos, maxPos, tolerance, maxAttempts int) float64 {
	if minPos == 0 && maxPos == 0 {
		return 0
	}
	if maxPos < 10 {
		maxPos = 10
	}
	for attempts := 0; attempts < maxAttempts; attempts++ {
		x := rand.Intn(maxPos-minPos) + minPos
		if x < int(lastSpawnX)-tolerance || x > int(lastSpawnX)+tolerance {
			return float64(x)
		}
		// Reduce tolerance to increase chances of finding a valid spawn point.
		tolerance /= 2
		if tolerance < 1 {
			tolerance = 1
		}
	}
	// If no valid value is found within maxAttempts, return a default value.
	return float64(minPos + rand.Intn(maxPos-minPos))
}

func bufferApply(collected, buffer int) int {
	if collected == buffer {
		return collected
	}
	if collected < buffer {
		delta := int((float64(buffer-collected) * 0.5))
		if delta < 1 {
			delta = 1
		}
		collected += delta
	}
	if collected > buffer {
		collected = buffer
	}

	return collected
}
