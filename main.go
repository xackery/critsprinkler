package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"os"

	"github.com/ebitenui/ebitenui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/xackery/critsprinkler/aa"
	"github.com/xackery/critsprinkler/bubble"
	"github.com/xackery/critsprinkler/config"
	"github.com/xackery/critsprinkler/dialog"
	"github.com/xackery/critsprinkler/dps"
	"github.com/xackery/critsprinkler/library"
	"github.com/xackery/critsprinkler/menu"
	"github.com/xackery/critsprinkler/money"
	"github.com/xackery/critsprinkler/placement"
	"github.com/xackery/critsprinkler/popup"
	"github.com/xackery/critsprinkler/sound"
	"github.com/xackery/critsprinkler/status"
	"github.com/xackery/critsprinkler/tracker"
)

var (
	cfg                      *config.CritSprinklerConfiguration
	game                     *Game
	Version                  string
	fontDefault              text.Face
	lastLayoutX, lastLayoutY int
)

type Game struct {
	ui         *ebitenui.UI
	statusText string
	//hwnd       windows.HWND

	isEditMode bool
}

func main() {
	err := run()
	if err != nil {
		dialog.MsgBox("Error", err.Error())
		os.Exit(1)
	}
}

func run() error {
	var err error
	cfg, err = config.LoadCritSprinklerConfig()
	if err != nil {
		return fmt.Errorf("%s: %w", config.FileName(), err)
	}

	err = library.New()
	if err != nil {
		return fmt.Errorf("library new: %w", err)
	}

	game = &Game{
		//aaPerHour: "AA per Hour: 0",
	}
	game.ui, err = menu.New(cfg, onSave, onEQPathLoad)
	if err != nil {
		return fmt.Errorf("gui: %w", err)
	}
	err = placement.New(game.ui, cfg)
	if err != nil {
		return fmt.Errorf("placement new: %w", err)
	}
	fontDefault, err = library.FontByKey(library.FontDefault)
	if err != nil {
		return fmt.Errorf("font %s: %w", library.FontDefault.String(), err)
	}

	t, err := tracker.New(cfg.LogPath)
	if err != nil {
		return fmt.Errorf("tracker: %w", err)
	}

	_, err = aa.New()
	if err != nil {
		return fmt.Errorf("aa: %w", err)
	}

	err = bubble.New()
	if err != nil {
		return fmt.Errorf("bubble: %w", err)
	}
	err = dps.New()
	if err != nil {
		return fmt.Errorf("dps: %w", err)
	}
	err = popup.New(cfg)
	if err != nil {
		return fmt.Errorf("popup: %w", err)
	}
	err = money.New(game.ui, cfg)
	if err != nil {
		return fmt.Errorf("money: %w", err)
	}
	err = sound.New(cfg)
	if err != nil {
		return fmt.Errorf("sound: %w", err)
	}

	err = t.Start(true)
	if err != nil {
		return fmt.Errorf("tracker start: %w", err)
	}
	//ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("CritSprinkler " + Version)
	icons, err := library.AppIcons()
	if err != nil {
		return fmt.Errorf("appicons: %w", err)
	}

	ebiten.SetWindowIcon(icons)

	//ebiten.SetWindowDecorated(false)
	//ebiten.SetWindowMousePassthrough(true)
	//ebiten.SetWindowFloating(true)
	//ebiten.SetWindowPosition(0, 0)
	fmt.Println("Showing game window")
	ebiten.SetWindowPosition(cfg.MainWindow.Min.X, cfg.MainWindow.Min.Y)
	ebiten.SetWindowSize(cfg.MainWindow.Dx(), cfg.MainWindow.Dy())
	ebiten.SetWindowSizeLimits(700, 700, -1, -1)

	err = ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{
		//SkipTaskbar:       true,
		ScreenTransparent: true,
		//InitUnfocused:     true,
	})
	if err != nil {
		return fmt.Errorf("rungame: %v", err)
	}

	return nil
}

func (g *Game) Update() error {
	if ebiten.IsFocused() && !g.IsEditMode() {
		g.setEditMode(true)
	}

	if !ebiten.IsFocused() && g.IsEditMode() {
		g.setEditMode(false)
	}

	if g.IsEditMode() {
		g.ui.Update()
	}

	bubble.Update()
	popup.Update()
	money.Update()
	// spawn a random pop up every 3s
	/*if rand.Intn(3) == 0 {
		g.spawnPopup(&dps.DamageEvent{
			Source:     "Test",
			Target:     "Test",
			SpellName:  "Ice Comet",
			Damage:     100,
			IsCritical: true,
		})
	}*/

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	//screen.Fill(color.RGBA{0, 0, 0, 0})
	screen.Clear()

	//	if g.IsEditMode() {
	if g.IsEditMode() { // || win.GetActiveWindowTitle() == "EverQuest" {
		//		vector.DrawFilledRect(screen, 0, 0, 50, 50, color.RGBA{0, 0, 0, 255}, true)
		g.statusBarDraw(screen)
	}
	g.ui.Draw(screen)
	popup.Draw(screen)
	money.Draw(screen)
	//win.GetActiveWindowTitle() == "EverQuest"

	//	}

	// fushia pink
	//screen.Fill(color.RGBA{255, 0, 255, 255})
	//screen.Fill(color.Black)

	// Draw AA per hour
	//if g.aaPerHour != "" {
	//	text.Draw(screen, g.aaPerHour, g.font, 50, 50, color.White)
	//}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != lastLayoutX || outsideHeight != lastLayoutY {
		lastLayoutX = outsideWidth
		lastLayoutY = outsideHeight
		g.onResize()
	}
	return outsideWidth, outsideHeight
}

func (g *Game) onResize() {
	placement.OnResize()
	money.OnResize()
}

func updateSave() error {
	var err error
	if cfg == nil {
		return fmt.Errorf("config not loaded")
	}

	cfg.MainWindow.Min.X, cfg.MainWindow.Min.Y = ebiten.WindowPosition()
	cfg.MainWindow.Max.X, cfg.MainWindow.Max.Y = ebiten.WindowSize()
	cfg.MainWindow.Max.X += cfg.MainWindow.Min.X
	cfg.MainWindow.Max.Y += cfg.MainWindow.Min.Y

	err = cfg.Save()
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

// SetEditMode sets if the game is in edit mode
func (g *Game) setEditMode(isEditMode bool) {
	g.isEditMode = isEditMode
	menu.SetEditMode(g.ui, isEditMode)
	placement.SetEditMode(g.ui, isEditMode)
	money.SetEditMode(g.ui, isEditMode)
	fmt.Println("edit mode is now", isEditMode)
	go func() {
		ebiten.SetWindowMousePassthrough(!isEditMode)
		ebiten.SetWindowFloating(!isEditMode)
		if !cfg.IsFullscreenBorderless {
			ebiten.SetWindowDecorated(isEditMode)

			resizeMode := ebiten.WindowResizingModeDisabled
			if isEditMode {
				resizeMode = ebiten.WindowResizingModeEnabled
			}
			ebiten.SetWindowResizingMode(resizeMode)
		}
	}()
}

// IsEditMode returns if the game is in edit mode
func (g *Game) IsEditMode() bool {
	return g.isEditMode
}

func (g *Game) statusBarDraw(screen *ebiten.Image) {
	//x, y := ebiten.WindowPosition()
	screenWidth, screenHeight := ebiten.WindowSize()

	vector.DrawFilledRect(screen, 0, float32(screenHeight-40), float32(screenWidth), float32(screenHeight), color.RGBA{0, 0, 0, 190}, true)
	op := &text.DrawOptions{}
	op.GeoM.Translate(10, float64(screenHeight-35))
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, status.String(), fontDefault, op)
}

func onSave() {
	err := updateSave()
	if err != nil {
		dialog.MsgBox("Error", err.Error())
	}
}

func onEQPathLoad() {
	money.OnEQPathLoad()
	sound.OnEQPathLoad()
}
