package main

import (
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/xackery/critsprinkler/aa"
	"github.com/xackery/critsprinkler/dps"
	"github.com/xackery/critsprinkler/tracker"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	damageEvent = make(chan *dps.DamageEvent)
	isLastLeft  = false
	spellColors = map[string]color.RGBA{
		"Ice Comet":           {120, 120, 255, 255},
		"Enticement of Flame": {255, 255, 0, 255},
		"Lifespike":           {255, 0, 0, 255},
		"Drain Soul":          {128, 0, 0, 255},
		"Thunderclap":         {255, 0, 255, 255},
		"Supernova":           {255, 165, 0, 255},
		"Word of Souls":       {230, 230, 250, 255},
	}
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
	spawnXLocs      []float64
	spawnXBuffer    []float64
	lastSpawnX      float64
	spawnYLocs      []float64
	spawnYBuffer    []float64
	lastSpawnY      float64
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {

	game := &Game{
		spawnXLocs: []float64{
			0,
			10, -10,
			20, -20,
			30, -30,
			40, -40,
			50, -50,
		},
		spawnYLocs: []float64{
			0,
			10, -10,
			20, -20,
			30, -30,
			40, -40,
			50, -50,
		},
	}
	new := flag.Bool("new", true, "Parse new log file")

	flag.Parse()
	if flag.NArg() < 1 {
		return fmt.Errorf("usage: %s <log file>, use -new to parse new data only, dps to enable dpsing", os.Args[0])
	}

	t, err := tracker.New(flag.Arg(0))
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
	game.font, err = loadFont("C:/Windows/Fonts/arial.ttf", 42)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	game.fontBorder, err = loadFont("C:/Windows/Fonts/arial.ttf", 42)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	game.smallFont, err = loadFont("C:/Windows/Fonts/arial.ttf", 24)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}
	game.smallFontBorder, err = loadFont("C:/Windows/Fonts/arial.ttf", 24)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	err = dps.SubscribeToDamageEvent(onDamageEvent)
	if err != nil {
		return fmt.Errorf("dps subscribe to damage event: %w", err)
	}

	if !*new {
		fmt.Println("Parsing entire log file")
	}

	err = t.Start(!*new)
	if err != nil {
		return fmt.Errorf("tracker start: %w", err)
	}
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Critsprinkler")
	err = ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{
		ScreenTransparent: true,
	})
	if err != nil {
		return fmt.Errorf("rungame: %v", err)
	}
	return nil
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

	select {
	case event := <-damageEvent:
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
	/* if rand.Intn(3) == 0 {
		g.spawnPopup(&dps.DamageEvent{
			Source:     "Test",
			Target:     "Test",
			SpellName:  "Ice Comet",
			Damage:     100,
			IsCritical: true,
		})
	} */

	return nil
}

func (g *Game) randomSpawnX(width int, attempts int) float64 {

	if len(g.spawnXBuffer) == 0 {
		g.spawnXBuffer = append(g.spawnXBuffer, g.spawnXLocs...)
	}

	if len(g.spawnXBuffer) == 0 {
		return float64(width / 2)
	}

	index := rand.Intn(len(g.spawnXBuffer))
	x := g.spawnXBuffer[index]
	if x < g.lastSpawnX+20 && x > g.lastSpawnX-20 && attempts < 10 {
		return g.randomSpawnX(width, attempts+1)
	}
	g.spawnXBuffer = append(g.spawnXBuffer[:index], g.spawnXBuffer[index+1:]...)
	return float64(width) + x
}

func (g *Game) randomSpawnY(height int, attempts int) float64 {

	if len(g.spawnYBuffer) == 0 {
		g.spawnYBuffer = append(g.spawnYBuffer, g.spawnYLocs...)
	}

	if len(g.spawnYBuffer) == 0 {
		return float64(height / 2)
	}

	index := rand.Intn(len(g.spawnYBuffer))
	y := g.spawnYBuffer[index]
	if y < g.lastSpawnY+20 && y > g.lastSpawnY-20 && attempts < 10 {
		return g.randomSpawnY(height, attempts+1)
	}
	g.spawnYBuffer = append(g.spawnYBuffer[:index], g.spawnYBuffer[index+1:]...)
	return float64(height) + y
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{0, 0, 0, 0})
	// fushia pink
	//screen.Fill(color.RGBA{255, 0, 255, 255})
	//screen.Fill(color.Black)
	// Draw all popups
	for _, popup := range g.popups {
		//alpha := uint8(255 * (popup.life + 100/popup.maxLife))
		col := color.RGBA{popup.color.R, popup.color.G, popup.color.B, 255}

		f := g.font
		fb := g.fontBorder
		if popup.isSmall {
			f = g.smallFont
			fb = g.smallFontBorder
		}

		text.Draw(screen, popup.text, f, int(popup.x)+1, int(popup.y)+1, color.Black)
		text.Draw(screen, popup.text, fb, int(popup.x), int(popup.y), col)

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) spawnPopup(event *dps.DamageEvent) {
	width, height := ebiten.WindowSize()
	x := g.randomSpawnX(width/2, 0)
	y := g.randomSpawnY(height/2, 0)

	isLastLeft = !isLastLeft

	spellColor, ok := spellColors[event.SpellName]
	if !ok {
		spellColor = color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		}
	}

	if event.Origin == "melee" {
		spellColor = color.RGBA{180, 180, 180, 255}
	}
	if event.Origin == "dot" {
		spellColor = color.RGBA{0, 150, 0, 255}
	}

	if event.Source != tracker.PlayerName() {
		return
	}
	popup := &Popup{
		text:    fmt.Sprintf("%d", event.Damage),
		x:       x,
		y:       y,
		isLeft:  (float64(width)/2 > float64(x)),
		vy:      -0.5 - rand.Float64() - (rand.Float64() / 2),
		life:    240,
		maxLife: 240,
		color:   spellColor,
		isSmall: !event.IsCritical,
	}
	if popup.isSmall || event.Origin == "melee" {
		popup.y += 200
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
	fmt.Println("Spawned popup at", x, y)
	g.popups = append(g.popups, popup)
}

func onDamageEvent(event *dps.DamageEvent) {
	damageEvent <- event
}
