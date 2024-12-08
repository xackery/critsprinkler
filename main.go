package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"

	"fyne.io/systray"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/xackery/critsprinkler/aa"
	"github.com/xackery/critsprinkler/config"
	"github.com/xackery/critsprinkler/dps"
	"github.com/xackery/critsprinkler/tracker"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
	"github.com/xackery/wlk/win"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/sys/windows"
)

//go:embed critsprinkler.ico
var iconData []byte

type SafeCounter struct {
	locked *atomic.Bool
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{
		locked: new(atomic.Bool), // Initialize the pointer
	}
}

func (c *SafeCounter) Lock() {
	c.locked.Store(true)
}

func (c *SafeCounter) Unlock() {
	c.locked.Store(false)
}

func (c *SafeCounter) IsLocked() bool {
	return c.locked.Load()
}

var (
	cfg             *config.CritSprinklerConfiguration
	damageEventChan = make(chan *dps.DamageEvent, 1000)
	isLastLeft      = false
	game            *Game
	mSettings       *systray.MenuItem
	mQuit           *systray.MenuItem
	settingsWnd     *walk.MainWindow
	prevFilePath    string
	updateTicker    *time.Ticker
	setPathButton   *walk.PushButton
)

type Popup struct {
	text    string
	isLeft  bool
	x, y    float64
	vy      float64
	life    float64
	maxLife float64
	color   color.RGBA
	isWave  bool
	waveMax float64
	waveMin float64
	isSmall bool
}

type Game struct {
	popups          []*Popup
	font            font.Face
	fontBorder      font.Face
	smallFont       font.Face
	smallFontBorder font.Face
	lastSpawnX      float64
	isWindowActive  bool

	isSettingsBeingChanged SafeCounter
	settingsWindowX        int
	settingsWindowY        int
	settingsWindowW        int
	settingsWindowH        int
	//aaPerHour              string
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("CritSprinkler")
	systray.SetTooltip("CritSprinkler")

	mSettings = systray.AddMenuItem("Settings", "Settings")
	mQuit = systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	//mQuit.SetIcon(icon.Data)
	//mSettings.SetIcon(icon.Data)

	updateTicker = time.NewTicker(100 * time.Millisecond)
	go sprinklerLoop()
}

func onExit() {
	fmt.Println("Exiting")
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var err error
	exePath := os.Args[0]

	wd := filepath.Dir(exePath)

	cfg, err = config.LoadCritSprinklerConfig(wd + "/critsprinkler.ini")
	if err != nil {
		fmt.Printf("load crit sprinkler config: %v\n", err)
	}

	if cfg.SettingsX == 0 && cfg.SettingsY == 0 {
		x, y := ebiten.Monitor().Size()
		cfg.SettingsX = x / 2
		cfg.SettingsY = y / 2
	}

	cmw := cpl.MainWindow{
		Title:    "Critsprinkler Settings",
		Name:     "sink",
		AssignTo: &settingsWnd,
		Size:     cpl.Size{Width: cfg.SettingsW, Height: cfg.SettingsH},
		Layout:   cpl.VBox{},

		Children: []cpl.Widget{
			cpl.Composite{
				Layout: cpl.HBox{},
				Children: []cpl.Widget{
					cpl.PushButton{
						Text:    "Test",
						MaxSize: cpl.Size{Width: 45},
						OnClicked: func() {
							onDamageEvent(&dps.DamageEvent{
								Orientation: dps.OrientationTopLeft,
								Source:      tracker.PlayerName(),
								Target:      "Test",
								SpellName:   "Ice Comet",
								Damage:      rand.Int() % 1000,
								IsCritical:  true,
							})
						},
					},
					cpl.PushButton{
						Text:    "Test",
						MaxSize: cpl.Size{Width: 45},
						OnClicked: func() {
							onDamageEvent(&dps.DamageEvent{
								Orientation: dps.OrientationTop,
								Source:      tracker.PlayerName(),
								Target:      "Test",
								SpellName:   "Ice Comet",
								Damage:      rand.Int() % 1000,
								IsCritical:  true,
							})
						},
					},
					cpl.PushButton{
						Text:    "Test",
						MaxSize: cpl.Size{Width: 45},
						OnClicked: func() {
							onDamageEvent(&dps.DamageEvent{
								Orientation: dps.OrientationTopRight,
								Source:      tracker.PlayerName(),
								Target:      "Test",
								SpellName:   "Ice Comet",
								Damage:      rand.Int() % 1000,
								IsCritical:  true,
							})
						},
					},
				},
			},
			cpl.Label{
				Text:      "Drag this window\nand test the sprinkles\non close it saves\nquit via systray",
				Alignment: cpl.AlignHCenterVCenter,
			},
			cpl.PushButton{
				Text:      "Set EQ Log",
				OnClicked: onSetPath,
				AssignTo:  &setPathButton,
			},
			cpl.PushButton{
				Text:    "Crit Randomly",
				MaxSize: cpl.Size{Width: 45},
				OnClicked: func() {
					onDamageEvent(&dps.DamageEvent{
						Source:     tracker.PlayerName(),
						Target:     "Test",
						SpellName:  "Ice Comet",
						Damage:     rand.Int() % 1000,
						IsCritical: true,
					})
				},
			},
			cpl.PushButton{
				Text:    "Non-Crit Randomly",
				MaxSize: cpl.Size{Width: 45},
				OnClicked: func() {
					onDamageEvent(&dps.DamageEvent{
						Source:     tracker.PlayerName(),
						Target:     "Test",
						SpellName:  "Ice Comet",
						Damage:     rand.Int() % 1000,
						IsCritical: false,
					})
				},
			},
			/* cpl.Composite{
				Layout: cpl.HBox{},
				Children: []cpl.Widget{
					cpl.CheckBox{
						Text:    "Show AA per Hour",
						Checked: false,
					},
					cpl.PushButton{
						Text: "Adjust",
						MaxSize: cpl.Size{Width: 45},
						OnClicked: func() {
							// show input box
						},
					},
				},
			}, */
			cpl.PushButton{
				Text:    "Save",
				MaxSize: cpl.Size{Width: 45},
				OnClicked: func() {
					err := updateSave()
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}
					settingsWnd.SetVisible(false)

				},
			},
			cpl.Composite{
				Layout: cpl.HBox{},
				Children: []cpl.Widget{
					cpl.PushButton{
						Text:    "Test",
						MaxSize: cpl.Size{Width: 45},
						OnClicked: func() {
							onDamageEvent(&dps.DamageEvent{
								Orientation: dps.OrientationBottomLeft,
								Source:      tracker.PlayerName(),
								Target:      "Test",
								SpellName:   "Ice Comet",
								Damage:      rand.Int() % 1000,
								IsCritical:  false,
							})
						},
					},
					cpl.PushButton{
						Text:    "Test",
						MaxSize: cpl.Size{Width: 45},
						OnClicked: func() {
							onDamageEvent(&dps.DamageEvent{
								Orientation: dps.OrientationBottom,
								Source:      tracker.PlayerName(),
								Target:      "Test",
								SpellName:   "Ice Comet",
								Damage:      rand.Int() % 1000,
								IsCritical:  false,
							})
						},
					},
					cpl.PushButton{
						Text:    "Test",
						MaxSize: cpl.Size{Width: 45},
						OnClicked: func() {
							onDamageEvent(&dps.DamageEvent{
								Orientation: dps.OrientationBottomRight,
								Source:      tracker.PlayerName(),
								Target:      "Test",
								SpellName:   "Ice Comet",
								Damage:      rand.Int() % 1000,
								IsCritical:  false,
							})
						},
					},
				},
			},
		},
		OnSizeChanged: func() {
			if game == nil {
				return
			}
			game.isSettingsBeingChanged.Lock()
			fmt.Println(settingsWnd.X(), settingsWnd.Y(), settingsWnd.Width(), settingsWnd.Height())
			game.settingsWindowX = settingsWnd.X()
			game.settingsWindowY = settingsWnd.Y()
			game.settingsWindowW = settingsWnd.Width()
			game.settingsWindowH = settingsWnd.Height()
			game.isSettingsBeingChanged.Unlock()
		},
		OnMouseMove: func(x, y int, button walk.MouseButton) {
			if game == nil {
				return
			}
			game.isSettingsBeingChanged.Lock()
			game.settingsWindowX = settingsWnd.X()
			game.settingsWindowY = settingsWnd.Y()
			game.settingsWindowW = settingsWnd.Width()
			game.settingsWindowH = settingsWnd.Height()
			game.isSettingsBeingChanged.Unlock()
		},

		OnBoundsChanged: func() {
			if game == nil {
				return
			}
			game.isSettingsBeingChanged.Lock()
			game.settingsWindowX = settingsWnd.X()
			game.settingsWindowY = settingsWnd.Y()
			game.settingsWindowW = settingsWnd.Width()
			game.settingsWindowH = settingsWnd.Height()
			game.isSettingsBeingChanged.Unlock()
		},

		Visible: cfg.LogPath == "",
	}

	err = cmw.Create()
	if err != nil {
		return fmt.Errorf("create main window: %w", err)
	}

	settingsWnd.Closing().Attach(func(isCancel *bool, reason byte) {
		*isCancel = true
		err := updateSave()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		settingsWnd.SetVisible(false)
	})

	settingsWnd.SetWidth(cfg.SettingsW)
	settingsWnd.SetHeight(cfg.SettingsH)
	settingsWnd.SetX(cfg.SettingsX)
	settingsWnd.SetY(cfg.SettingsY)
	prevFilePath = cfg.LogPath

	go func() {

		code := settingsWnd.Run()
		if code != 0 {
			fmt.Printf("mainWalk Error: %v\n", code)
		}
		DisableMinMaxTitle(settingsWnd.Handle())

		if cfg.LogPath == "" {

			setPathButton.Button.Clicked()
			win.SetForegroundWindow(settingsWnd.Handle())
			win.SetActiveWindow(settingsWnd.Handle())

		}

	}()
	game = &Game{
		isSettingsBeingChanged: *NewSafeCounter(),
		//aaPerHour: "AA per Hour: 0",
	}

	path := ""

	if cfg.LogPath != "" {
		path = cfg.LogPath
	}

	game.isSettingsBeingChanged.Lock()
	game.settingsWindowX = cfg.SettingsX
	game.settingsWindowY = cfg.SettingsY
	game.settingsWindowW = cfg.SettingsW
	game.settingsWindowH = cfg.SettingsH
	game.isSettingsBeingChanged.Unlock()

	t, err := tracker.New(path)
	if err != nil {
		return fmt.Errorf("tracker: %w", err)
	}

	_, err = aa.New()
	if err != nil {
		return fmt.Errorf("aa: %w", err)
	}

	_, err = dps.New()
	if err != nil {
		return fmt.Errorf("dps: %w", err)
	}

	// get windows system path

	winPath := os.Getenv("WINDIR")
	if winPath == "" {
		winPath = "C:/Windows"
	}

	fontPath := fmt.Sprintf("%s/Fonts/arial.ttf", winPath)

	game.font, err = loadFont(fontPath, 42)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	game.fontBorder, err = loadFont(fontPath, 42)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	game.smallFont, err = loadFont(fontPath, 36)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}
	game.smallFontBorder, err = loadFont(fontPath, 36)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	err = dps.SubscribeToDamageEvent(onDamageEvent)
	if err != nil {
		return fmt.Errorf("dps subscribe to damage event: %w", err)
	}

	err = t.Start(true)
	if err != nil {
		return fmt.Errorf("tracker start: %w", err)
	}
	go func() {
		fmt.Println("Showing SysTray")
		systray.Run(onReady, onExit)
	}()

	// get screen size
	screenWidth, screenHeight := ebiten.Monitor().Size()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Critsprinkler")
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowPosition(0, 0)
	fmt.Println("Showing game window")
	go func() {
		time.Sleep(1 * time.Second)

		hwnd := win.FindWindow(nil, StringToUTF16Ptr("Critsprinkler"))

		fmt.Println("hwnd", hwnd)
		HideFromTaskbar(hwnd)
	}()
	err = ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{
		SkipTaskbar:       true,
		ScreenTransparent: true,
		InitUnfocused:     true,
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
	exStyle &^= win.WS_EX_APPWINDOW // Remove WS_EX_APPWINDOW
	exStyle |= win.WS_EX_TOOLWINDOW // Add WS_EX_TOOLWINDOW

	// Apply the new style
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE, exStyle)

	// Ensure changes take effect
	win.ShowWindow(hwnd, win.SW_HIDE) // Temporarily hide the window
	win.ShowWindow(hwnd, win.SW_SHOW) // Show it again with the new style
}

func StringToUTF16Ptr(s string) *uint16 {
	ptr, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return ptr
}

func sprinklerLoop() {
	for {
		select {
		case <-mSettings.ClickedCh:
			if !settingsWnd.Visible() {
				fmt.Println("Showing mainWalk")
				settingsWnd.SetVisible(true)
				DisableMinMaxTitle(settingsWnd.Handle())
			}
			win.SetForegroundWindow(settingsWnd.Handle())
			win.SetActiveWindow(settingsWnd.Handle())
		case <-mQuit.ClickedCh:
			systray.Quit()
			os.Exit(0)
		case <-updateTicker.C:

			if win.GetActiveWindowTitle() != "Critsprinkler Settings" &&
				win.GetActiveWindowTitle() != "EverQuest" {
				game.isWindowActive = false
				continue
			}

			game.isWindowActive = true
			/*
				if win.GetActiveWindowTitle() == "Critsprinkler Settings" {
					ebiten.SetWindowPosition(mainWalk.X(), mainWalk.Y())
					ebiten.SetWindowSize(mainWalk.Width(), mainWalk.Height())
				} */

		}
	}
}

func loadFont(path string, size float64) (font.Face, error) {
	fontData, err := ioutil.ReadFile(path)
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
	return face, nil
}

func (g *Game) Update() error {

	if game.isSettingsBeingChanged.IsLocked() {
		return nil
	}

	select {
	case event := <-damageEventChan:
		g.spawnPopup(event)
	default:
	}
	// Update existing popups
	for i := len(g.popups) - 1; i >= 0; i-- {
		popup := g.popups[i]
		popup.y += popup.vy * 0.9
		popup.life -= 1
		if popup.isSmall {
			popup.life -= 1
			popup.y += 0.5
		}

		if popup.isWave {
			if popup.y > popup.waveMax {
				popup.vy = -0.5
			}
			if popup.y < popup.waveMin {
				popup.vy = 0.5
			}
		}

		// Reverse velocity after reaching the hover point
		if popup.life < popup.maxLife*0.5 && popup.vy < 0 {
			popup.vy = -popup.vy * 0.2 // Slows down
			if popup.isSmall {
				popup.vy = -popup.vy * 0.1 // Slows down
			}
			popup.x += 0.4
			if popup.isLeft {
				popup.x -= 0.3
			}
		} else {
			popup.x += 0.2
			if popup.isLeft {
				popup.x -= 0.4
			}
		}

		// Remove popups that have expired
		if popup.life <= 0 {
			g.popups = append(g.popups[:i], g.popups[i+1:]...)
		}
	}

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

func (g *Game) randomSpawnX(minX, maxX, tolerance, attempts int) float64 {
	if minX == 0 && maxX == 0 {
		return 0
	}
	x := rand.Intn(maxX-minX) + minX
	if x < int(g.lastSpawnX)+tolerance && x > int(g.lastSpawnX)-tolerance && attempts < 10 {
		tolerance = tolerance / 2
		if tolerance < 1 {
			tolerance = 1
		}
		return g.randomSpawnX(minX, maxX, tolerance, attempts+1)
	}
	return float64(x)

}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{0, 0, 0, 0})

	if game.isSettingsBeingChanged.IsLocked() {
		return
	}

	if !game.isWindowActive {
		return
	}

	// fushia pink
	//screen.Fill(color.RGBA{255, 0, 255, 255})
	//screen.Fill(color.Black)
	// Draw all popups
	for _, popup := range g.popups {
		//alpha := uint8(255 * (popup.life + 100/popup.maxLife))
		col := color.RGBA{popup.color.R, popup.color.G, popup.color.B, 255}

		f := g.font
		fb := g.fontBorder
		offset := 2
		if popup.isSmall {
			f = g.smallFont
			fb = g.smallFontBorder
			offset = 2
		}

		text.Draw(screen, popup.text, f, int(popup.x)+offset, int(popup.y)+offset, color.Black)
		text.Draw(screen, popup.text, fb, int(popup.x), int(popup.y), col)

	}

	// Draw AA per hour
	//if g.aaPerHour != "" {
	//	text.Draw(screen, g.aaPerHour, g.font, 50, 50, color.White)
	//}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) spawnPopup(event *dps.DamageEvent) {

	isLastLeft = !isLastLeft

	spellColor, ok := spellColors[event.SpellName]
	if !ok {
		/* spellColor = color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		} */
		fmt.Println("Spell not found", event.SpellName)
		spellColor = color.RGBA{255, 255, 255, 255}
	}

	/* if event.Origin == "melee" {
		spellColor = color.RGBA{180, 180, 180, 255}
	}
	if event.Origin == "dot" {
		spellColor = color.RGBA{0, 150, 0, 255}
	} */

	if event.Source == tracker.PlayerName() && event.Target == tracker.PlayerName() && event.Origin == "direct" {
		spellColor = color.RGBA{255, 0, 0, 255}
	}

	if event.Origin == "heal" {
		spellColor = color.RGBA{0, 147, 255, 255}
	}

	if event.Target != tracker.PlayerName() && event.Source != tracker.PlayerName() {
		return
	}

	if event.Target == tracker.PlayerName() && event.Origin != "heal" {
		spellColor = color.RGBA{128, 0, 0, 255}
	}

	fmt.Println(event, spellColor)
	var x float64

	y := float64(g.settingsWindowY)

	centerX := float64(g.settingsWindowX + g.settingsWindowW/2)

	if !event.IsCritical {
		y = float64(g.settingsWindowY + g.settingsWindowH)
	}

	switch event.Orientation {
	case dps.OrientationTopLeft:
		x = float64(g.settingsWindowX)
		y = float64(g.settingsWindowY)
	case dps.OrientationTop:
		x = float64(g.settingsWindowX + g.settingsWindowW/2)
		y = float64(g.settingsWindowY)
	case dps.OrientationTopRight:
		x = float64(g.settingsWindowX + g.settingsWindowW)
		y = float64(g.settingsWindowY)
	case dps.OrientationBottomLeft:
		x = float64(g.settingsWindowX)
		y = float64(g.settingsWindowY + g.settingsWindowH)
	case dps.OrientationBottom:
		x = float64(g.settingsWindowX + g.settingsWindowW/2)
		y = float64(g.settingsWindowY + g.settingsWindowH)
	case dps.OrientationBottomRight:
		x = float64(g.settingsWindowX + g.settingsWindowW)
		y = float64(g.settingsWindowY + g.settingsWindowH)
	default: // random position
		x = g.randomSpawnX(g.settingsWindowX, g.settingsWindowX+g.settingsWindowW, g.settingsWindowW/2, 0)
	}

	popup := &Popup{
		text:   fmt.Sprintf("%d", event.Damage),
		x:      x,
		y:      y,
		isLeft: (centerX > float64(x)),
		//vy:     -0.5,
		vy:      -0.5 - rand.Float64() - (rand.Float64() / 2),
		life:    240,
		maxLife: 240,
		color:   spellColor,
		isSmall: !event.IsCritical,
	}
	if popup.isSmall || event.Origin == "melee" {
		//popup.y += 200
		popup.life = 250
		popup.maxLife = 500
	}
	if event.Origin == "dot" {
		popup.life = 240
		popup.maxLife = 240
		popup.isWave = true
		popup.waveMax = popup.y + 10
		popup.waveMin = popup.y - 10
	}
	g.popups = append(g.popups, popup)
}

func onDamageEvent(event *dps.DamageEvent) {
	damageEventChan <- event
}

func onSetPath() {
	var err error
	dia := new(walk.FileDialog)

	curDir := "."
	if prevFilePath != "" {
		curDir = filepath.Dir(prevFilePath)
	} else {
		curDir, err = os.Getwd()
		if err != nil {
			curDir = "."
		}
	}

	dia.FilePath = curDir
	dia.Filter = "Log Files (*.txt)|eqlog_*.txt"
	dia.Title = "Select Log File"

	ok, err := dia.ShowOpen(settingsWnd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if !ok {
		return
	}
	prevFilePath = dia.FilePath

	err = tracker.SetNewPath(prevFilePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Selected file", prevFilePath)
}

func updateSave() error {
	if cfg == nil {
		return fmt.Errorf("config not loaded")
	}
	cfg.LogPath = prevFilePath
	cfg.SettingsH = game.settingsWindowH
	cfg.SettingsW = game.settingsWindowW
	cfg.SettingsX = game.settingsWindowX
	cfg.SettingsY = game.settingsWindowY

	err := cfg.Save()
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}
