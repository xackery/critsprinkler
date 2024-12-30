package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/png"
	"os"
	"syscall"

	"github.com/ebitenui/ebitenui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/xackery/critsprinkler/aa"
	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/config"
	"github.com/xackery/critsprinkler/dps"
	"github.com/xackery/critsprinkler/popup"
	"github.com/xackery/critsprinkler/tracker"
	"github.com/xackery/wlk/win"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/sys/windows"
)

var (
	cfg             *config.CritSprinklerConfiguration
	damageEventChan = make(chan *common.DamageEvent, 10000)
	game            *Game
	path            string
	Version         string
)

type Game struct {
	ui         *ebitenui.UI
	statusText string
	res        *uiResources
	//hwnd       windows.HWND

	isEditMode bool
}

func main() {
	err := run()
	if err != nil {
		MsgBox("Error", err.Error())
		os.Exit(1)
	}
}

func run() error {
	var err error
	cfg, err = config.LoadCritSprinklerConfig()
	if err != nil {
		return fmt.Errorf("%s: %w", config.FileName(), err)
	}

	game = &Game{
		//aaPerHour: "AA per Hour: 0",
	}

	game.ui, err = NewGUI(cfg, game)
	if err != nil {
		return fmt.Errorf("new gui: %w", err)
	}
	if cfg.LogPath != "" {
		path = cfg.LogPath
	}

	t, err := tracker.New(path)
	if err != nil {
		return fmt.Errorf("tracker: %w", err)
	}

	_, err = aa.New()
	if err != nil {
		return fmt.Errorf("aa: %w", err)
	}

	err = dps.New()
	if err != nil {
		return fmt.Errorf("dps: %w", err)
	}

	// get windows system path

	winPath := os.Getenv("WINDIR")
	if winPath == "" {
		winPath = "C:/Windows"
	}

	game.isEditMode = true

	err = popup.New(cfg, game.res.fonts.popupBold42)
	if err != nil {
		return fmt.Errorf("popup: %w", err)
	}

	err = dps.SubscribeToDamageEvent(onDamageEvent)
	if err != nil {
		return fmt.Errorf("dps subscribe to damage event: %w", err)
	}

	err = t.Start(true)
	if err != nil {
		return fmt.Errorf("tracker start: %w", err)
	}

	//ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Critsprinkler " + Version)
	icons := []image.Image{}
	iconPaths := []string{
		"assets/icon_16.png",
		"assets/icon_24.png",
		"assets/icon_32.png",
	}
	for _, iconPath := range iconPaths {
		r, err := embeddedAssets.Open(iconPath)
		if err != nil {
			return fmt.Errorf("open icon: %w", err)
		}
		defer r.Close()
		img, err := png.Decode(r)
		if err != nil {
			return fmt.Errorf("decode icon: %w", err)
		}
		icons = append(icons, img)
	}

	ebiten.SetWindowIcon(icons)
	game.setEditMode(true)

	//ebiten.SetWindowIcon(iconData)
	//ebiten.SetWindowDecorated(false)
	//ebiten.SetWindowMousePassthrough(true)
	//ebiten.SetWindowFloating(true)
	//ebiten.SetWindowPosition(0, 0)
	fmt.Println("Showing game window")
	// go func() {
	// 	attempts := 0
	// 	for {
	// 		time.Sleep(50 * time.Millisecond)

	// 		game.hwnd = win.FindWindow(nil, StringToUTF16Ptr("Critsprinkler"))
	// 		if game.hwnd == 0 {
	// 			attempts++
	// 			if attempts > 1000 {
	// 				fmt.Println("Failed to find window")
	// 				os.Exit(1)
	// 			}
	// 			continue
	// 		}

	// 		//HideFromTaskbar(hwnd)
	// 		time.Sleep(50 * time.Millisecond)
	// 		return
	// 	}
	// }()
	ebiten.SetWindowPosition(cfg.MainWindow.Min.X, cfg.MainWindow.Min.Y)
	ebiten.SetWindowSize(cfg.MainWindow.Dx(), cfg.MainWindow.Dy())
	ebiten.SetWindowSizeLimits(700, 700, 10000, 10000)
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

func DisableMinMaxTitle(hwnd windows.HWND) {
	// Get current style
	style := win.GetWindowLong(hwnd, win.GWL_STYLE)

	// Modify styles: Remove WS_MAXIMIZEBOX, WS_MINIMIZEBOX, WS_CAPTION
	style &^= win.WS_MAXIMIZEBOX // Remove WS_MAXIMIZEBOX
	style &^= win.WS_MINIMIZEBOX // Remove WS_MINIMIZEBOX
	//style &^= win.WS_CAPTION     // Remove WS_CAPTION

	// Apply the new style
	win.SetWindowLong(hwnd, win.GWL_STYLE, style)
	win.UpdateWindow(hwnd)
}

func HideFromTaskbar(hwnd windows.HWND) {
	// Get current extended style
	exStyle := win.GetWindowLong(hwnd, win.GWL_EXSTYLE)

	// Modify styles: Remove WS_EX_APPWINDOW, add WS_EX_TOOLWINDOW
	//exStyle &^= win.WS_EX_APPWINDOW // Remove WS_EX_APPWINDOW
	exStyle |= win.WS_EX_TOOLWINDOW // Add WS_EX_TOOLWINDOW
	exStyle |= win.WS_EX_LAYERED

	// Apply the new style
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, exStyle)
	win.SetWindowPos(hwnd, 0, 0, 0, 0, 0, win.SWP_NOSIZE|win.SWP_NOMOVE|win.SWP_NOZORDER|win.SWP_FRAMECHANGED)
	// Ensure changes take effect
	//win.ShowWindow(hwnd, win.SW_HIDE) // Temporarily hide the window
	//win.ShowWindow(hwnd, win.SW_SHOW) // Show it again with the new style
}

func StringToUTF16Ptr(s string) *uint16 {
	ptr, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return ptr
}

func loadFont(path string, size float64) (text.Face, error) {
	fontData, err := embeddedAssets.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone, // Disable hinting for a crisp look
	})
	if err != nil {
		return nil, err
	}
	return text.NewGoXFace(face), nil
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
	select {
	case event := <-damageEventChan:
		err := popup.Spawn(event)
		if err != nil {
			g.statusText = fmt.Sprintf("Error spawning popup: %v", err)
			fmt.Println(g.statusText)
		}

	default:
	}

	popup.Update()

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
		g.ui.Draw(screen)
		g.statusBarDraw(screen)
	}

	//win.GetActiveWindowTitle() == "EverQuest"

	//	}

	// fushia pink
	//screen.Fill(color.RGBA{255, 0, 255, 255})
	//screen.Fill(color.Black)
	popup.Draw(screen)

	// Draw AA per hour
	//if g.aaPerHour != "" {
	//	text.Draw(screen, g.aaPerHour, g.font, 50, 50, color.White)
	//}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func onDamageEvent(event *common.DamageEvent) {
	damageEventChan <- event
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

	popup.ConfigUpdate(cfg)

	err = cfg.Save()
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

// SetEditMode sets if the game is in edit mode
func (g *Game) setEditMode(isEditMode bool) {
	g.isEditMode = isEditMode
	go func() {
		ebiten.SetWindowMousePassthrough(!isEditMode)
		ebiten.SetWindowDecorated(isEditMode)
		ebiten.SetWindowFloating(!isEditMode)
		resizeMode := ebiten.WindowResizingModeDisabled
		if isEditMode {
			resizeMode = ebiten.WindowResizingModeEnabled
		}
		ebiten.SetWindowResizingMode(resizeMode)
	}()
}

// IsEditMode returns if the game is in edit mode
func (g *Game) IsEditMode() bool {
	return g.isEditMode
}
