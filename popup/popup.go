package popup

import (
	"fmt"
	"image/color"
	"strconv"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/xackery/critsprinkler/bubble"
	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/config"
	"github.com/xackery/critsprinkler/placement"
	"github.com/xackery/critsprinkler/tracker"
	"golang.org/x/exp/rand"
)

var (
	mu             sync.RWMutex
	popups         []*Popup
	tallyDuration  *time.Duration
	isCommaEnabled *bool
)

type Popup struct {
	category       common.PopupCategory
	text           string
	isTallyEnabled bool
	currentDamage  int
	targetDamage   int
	baseX, baseY   *int
	face           *text.Face
	startX, startY float64
	x, y           float64
	vx, vy         float64
	life           float64
	maxLife        float64
	color          color.RGBA
	isWave         bool
	waveMax        float64
	waveMin        float64
	isSmall        bool
	tallyEndTime   time.Time
}

func New(cfg *config.CritSprinklerConfiguration) error {
	mu.Lock()
	defer mu.Unlock()
	isCommaEnabled = &cfg.IsCommaEnabled
	tallyDuration = &cfg.PopupTallyDuration

	return nil
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

// Update is called by ebiten to update the popup animations
func Update() {

	spawns := bubble.Spawns()
	for _, event := range spawns {
		err := spawn(event)
		if err != nil {
			fmt.Println("spawn error:", err)
		}
	}

	for i := len(popups) - 1; i >= 0; i-- {
		popup := popups[i]

		popup.y += popup.vy * 0.9
		popup.x += popup.vx
		popup.life -= 1

		if popup.currentDamage < popup.targetDamage {
			delta := int((float64(popup.targetDamage-popup.currentDamage) * 0.5))
			if delta < 1 {
				delta = 1
			}
			popup.currentDamage += delta

			popup.text = strconv.Itoa(popup.currentDamage)
			if *isCommaEnabled && popup.currentDamage > 0 {
				popup.text = commaFormat(popup.currentDamage)
			}
		}

		// Reverse velocity after reaching the hover point
		/* if !popup.isTallyEnabled && popup.life < popup.maxLife*0.5 && popup.vy < 0 {
			popup.vy = -popup.vy * 0.2 // Slows down
		} */

		// Remove popups that have expired
		if popup.life <= 0 {
			popups = append(popups[:i], popups[i+1:]...)
		}
	}
}

// Draw is called by ebiten to draw the popups
func Draw(screen *ebiten.Image) {
	for _, popup := range popups {
		//alpha := uint8(255 * (popup.life + 100/popup.maxLife))
		shadow := color.RGBA{0, 0, 0, popup.color.A}
		col := color.RGBA{popup.color.R, popup.color.G, popup.color.B, popup.color.A}

		offset := 2

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(*popup.baseX)+popup.x+float64(offset), float64(*popup.baseY)+popup.y+float64(offset))
		op.ColorScale.ScaleWithColor(shadow)
		text.Draw(screen, popup.text, *popup.face, op)
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(*popup.baseX)+popup.x, float64(*popup.baseY)+popup.y)
		op.ColorScale.ScaleWithColor(col)

		text.Draw(screen, popup.text, *popup.face, op)
	}
}

func spawn(event *common.DamageEvent) error {

	setting := placement.ByCategory(event.Category)
	if setting == nil {
		return fmt.Errorf("no setting found for %d", event.Category)
	}

	if setting.IsVisible == 0 {
		return spawnTotal(event)
	}

	spellColor := color.RGBA{255, 255, 255, 255}

	spellColor.R = setting.FontColor.R
	spellColor.G = setting.FontColor.G
	spellColor.B = setting.FontColor.B
	spellColor.A = setting.FontColor.A

	spellColor = color.RGBA{100, 100, 255, 255}

	/*
		spellColor, ok := spellColors[event.SpellName]
		if !ok {
			//  spellColor = color.RGBA{
			// 	R: uint8(rand.Intn(256)),
			// 	G: uint8(rand.Intn(256)),
			// 	B: uint8(rand.Intn(256)),
			// 	A: 255,
			// }
			//fmt.Println("Spell not found", event.SpellName)
			spellColor = color.RGBA{255, 255, 255, 255}
		} */

	/* if event.Origin == "melee" {
		spellColor = color.RGBA{180, 180, 180, 255}
	}
	if event.Origin == "dot" {
		spellColor = color.RGBA{0, 150, 0, 255}
	} */

	/* if event.Source == tracker.PlayerName() && event.Target == tracker.PlayerName() && event.Origin == "direct" {
		spellColor = color.RGBA{255, 0, 0, 255}
	}

	if event.Origin == "heal" {
		spellColor = color.RGBA{0, 147, 255, 255}
	}

	if event.Target != tracker.PlayerName() && event.Source != tracker.PlayerName() {
		return nil
	}

	if event.Target == tracker.PlayerName() && event.Origin != "heal" {
		spellColor = color.RGBA{128, 0, 0, 255}
	} */

	//  0
	// 100
	// -50, 50

	if event.Source != tracker.PlayerName() && event.Target != tracker.PlayerName() {
		return nil
	}

	if setting.IsTallyEnabled == 1 {
		for i := 0; i < len(popups); i++ {
			popup := popups[i]
			if time.Now().After(popup.tallyEndTime) {
				continue
			}
			if popup.category != event.Category {
				continue
			}

			val, err := strconv.Atoi(event.Damage)
			if err != nil {
				return nil
			}

			popup.targetDamage += val
			popup.x -= popup.vx
			popup.y -= popup.vy
			popup.life += 10
			return spawnTotal(event)
		}
	}

	vx := float64(0)
	vy := float64(-0.5 - rand.Float64() - (rand.Float64() / 2))
	switch setting.Direction {
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

	fmt.Printf("%s->%s->%s (%s) %s %s\n", event.Source, event.Type, event.Target, event.SpellName, event.Damage, event.Category.String())

	val, err := strconv.Atoi(event.Damage)
	if err != nil {
		val = 0
	}
	damage := event.Damage
	if *isCommaEnabled && val > 0 {
		damage = commaFormat(val)
	}

	popup := &Popup{
		category:      event.Category,
		text:          damage,
		currentDamage: val,
		targetDamage:  val,
		baseX:         &setting.WindowRect.Min.X,
		baseY:         &setting.WindowRect.Min.Y,
		face:          &setting.FontFace,
		x:             randomSpawnRange(int(setting.LastSpawnX), 0, setting.WindowRect.Dx()-50, setting.WindowRect.Dx()/4, 0),
		y:             randomSpawnRange(int(setting.LastSpawnX), 0, setting.WindowRect.Dy()-50, setting.WindowRect.Dy()/4, 0),
		vx:            vx,
		vy:            vy,
		life:          240,
		maxLife:       240,
		tallyEndTime:  time.Now().Add(*tallyDuration),
		color:         spellColor,
	}
	if setting.IsTallyEnabled == 1 {
		popup.maxLife += 1000
		popup.life += 1000
		switch setting.Direction {
		case common.DirectionUp:
			popup.vy = -0.2
		case common.DirectionDown:
			popup.vy = 0.2
		case common.DirectionLeft:
			popup.vx = -0.4
		case common.DirectionRight:
			popup.vx = 0.4
		case common.DirectionUpLeft:
			popup.vx = -0.2
			popup.vy = -0.2
		case common.DirectionUpRight:
			popup.vx = 0.2
			popup.vy = -0.2
		case common.DirectionDownLeft:
			popup.vx = -0.2
			popup.vy = 0.2
		case common.DirectionDownRight:
			popup.vx = 0.2
			popup.vy = 0.2
		}
		popup.x = float64(setting.WindowRect.Dx() / 2)
		popup.y = float64(setting.WindowRect.Dy() / 2)
	}
	popup.startX = popup.x
	popup.startY = popup.y
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

	popups = append(popups, popup)
	return spawnTotal(event)
}

// ConfigUpdate updates the popup configuration
func ConfigUpdate(cfg *config.CritSprinklerConfiguration) {
	mu.Lock()
	defer mu.Unlock()
}

func spawnTotal(event *common.DamageEvent) error {
	if common.IsTotalDamageOut(event.Category) {
		fmt.Println("total damage out")
		event.Category = common.PopupCategoryTotalDamageOut
		return spawn(event)
	}
	if common.IsTotalDamageIn(event.Category) {
		event.Category = common.PopupCategoryTotalDamageIn
		return spawn(event)
	}
	if common.IsTotalHealOut(event.Category) {
		event.Category = common.PopupCategoryTotalHealOut
		return spawn(event)
	}
	if common.IsTotalHealIn(event.Category) {
		event.Category = common.PopupCategoryTotalHealIn
		return spawn(event)
	}
	return nil
}

func commaFormat(num int) string {
	if num < 1000 {
		return strconv.Itoa(num)
	}
	in := strconv.Itoa(num)
	n := len(in) % 3
	out := in[:n]
	for i := n; i < len(in); i += 3 {
		if len(out) > 0 {
			out += ","
		}
		out += in[i : i+3]
	}
	return out
}

func (p *Popup) Clone() *Popup {
	return &Popup{
		text:          p.text,
		currentDamage: p.currentDamage,
		targetDamage:  p.targetDamage,
		x:             p.x,
		y:             p.y,
		vx:            p.vx,
		vy:            p.vy,
		life:          p.life,
		maxLife:       p.maxLife,
		color:         p.color,
		isWave:        p.isWave,
		waveMax:       p.waveMax,
		waveMin:       p.waveMin,
		isSmall:       p.isSmall,
		tallyEndTime:  p.tallyEndTime,
	}
}
