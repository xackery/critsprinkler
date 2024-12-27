package popup

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/config"
	"golang.org/x/exp/rand"
)

var (
	mu            sync.RWMutex
	popupSettings map[common.PopupCategory]*SettingProperty
	popups        []*Popup
)

type SettingProperty struct {
	IsEnabled    *bool
	Category     common.PopupCategory
	categoryName string // Melee, Spell, etc
	Title        string
	rect         *image.Rectangle
	lastSpawnX   float64
	Color        *color.RGBA
	font         text.Face
	fontBorder   text.Face
	Direction    *common.Direction
}

type Popup struct {
	setting *SettingProperty
	text    string
	x, y    float64
	vx, vy  float64
	life    float64
	maxLife float64
	color   color.RGBA
	isWave  bool
	waveMax float64
	waveMin float64
	isSmall bool
}

func New(cfg *config.CritSprinklerConfiguration, font text.Face) error {
	mu.Lock()
	defer mu.Unlock()
	popupSettings = make(map[common.PopupCategory]*SettingProperty)

	popupSettings[common.PopupCategoryGlobalCritOut] = &SettingProperty{
		Category:   common.PopupCategoryGlobalCritOut,
		font:       font,
		fontBorder: font,
		rect:       &cfg.GlobalCritOut,
		Title:      "Global Crit Outgoing",
		Color:      &cfg.GlobalCritOutColor,
		Direction:  &cfg.GlobalCritOutDirection,
	}

	popupSettings[common.PopupCategoryGlobalHitOut] = &SettingProperty{
		Category:     common.PopupCategoryGlobalHitOut,
		Title:        "Global Hit Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.GlobalHitOut,
		categoryName: "Global",
		Color:        &cfg.GlobalHitOutColor,
		Direction:    &cfg.GlobalHitOutDirection,
	}

	popupSettings[common.PopupCategoryGlobalMissOut] = &SettingProperty{
		Category:     common.PopupCategoryGlobalMissOut,
		Title:        "Global Miss Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.GlobalMissOut,
		categoryName: "Global",
		Color:        &cfg.GlobalMissOutColor,
		Direction:    &cfg.GlobalMissOutDirection,
	}

	popupSettings[common.PopupCategoryGlobalCritIn] = &SettingProperty{
		Category:     common.PopupCategoryGlobalCritIn,
		Title:        "Global Crit Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.GlobalCritIn,
		categoryName: "Global",
		Color:        &cfg.GlobalCritInColor,
		Direction:    &cfg.GlobalCritInDirection,
	}

	popupSettings[common.PopupCategoryGlobalHitIn] = &SettingProperty{
		Category:     common.PopupCategoryGlobalHitIn,
		Title:        "Global Hit Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.GlobalHitIn,
		categoryName: "Global",
		Color:        &cfg.GlobalHitInColor,
		Direction:    &cfg.GlobalHitInDirection,
	}

	popupSettings[common.PopupCategoryGlobalMissIn] = &SettingProperty{
		Category:     common.PopupCategoryGlobalMissIn,
		Title:        "Global Miss Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.GlobalMissIn,
		categoryName: "Global",
		Color:        &cfg.GlobalMissInColor,
		Direction:    &cfg.GlobalMissInDirection,
	}

	popupSettings[common.PopupCategoryMeleeCritOut] = &SettingProperty{
		Category:     common.PopupCategoryMeleeCritOut,
		Title:        "Melee Crit Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.MeleeCritOut,
		categoryName: "Melee",
		Color:        &cfg.MeleeCritOutColor,
		Direction:    &cfg.MeleeCritOutDirection,
		IsEnabled:    &cfg.MeleeCritOutIsEnabled,
	}

	popupSettings[common.PopupCategoryMeleeHitOut] = &SettingProperty{
		Category:     common.PopupCategoryMeleeHitOut,
		Title:        "Melee Hit Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.MeleeHitOut,
		categoryName: "Melee",
		Color:        &cfg.MeleeHitOutColor,
		Direction:    &cfg.MeleeHitOutDirection,
		IsEnabled:    &cfg.MeleeHitOutIsEnabled,
	}

	popupSettings[common.PopupCategoryMeleeMissOut] = &SettingProperty{
		Category:     common.PopupCategoryMeleeMissOut,
		Title:        "Melee Miss Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.MeleeMissOut,
		categoryName: "Melee",
		Color:        &cfg.MeleeMissOutColor,
		Direction:    &cfg.MeleeMissOutDirection,
		IsEnabled:    &cfg.MeleeMissOutIsEnabled,
	}

	popupSettings[common.PopupCategoryMeleeCritIn] = &SettingProperty{
		Category:     common.PopupCategoryMeleeCritIn,
		Title:        "Melee Crit Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.MeleeCritIn,
		categoryName: "Melee",
		Color:        &cfg.MeleeCritInColor,
		Direction:    &cfg.MeleeCritInDirection,
		IsEnabled:    &cfg.MeleeCritInIsEnabled,
	}

	popupSettings[common.PopupCategoryMeleeHitIn] = &SettingProperty{
		Category:     common.PopupCategoryMeleeHitIn,
		Title:        "Melee Hit Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.MeleeHitIn,
		categoryName: "Melee",
		Color:        &cfg.MeleeHitInColor,
		Direction:    &cfg.MeleeHitInDirection,
		IsEnabled:    &cfg.MeleeHitInIsEnabled,
	}

	popupSettings[common.PopupCategoryMeleeMissIn] = &SettingProperty{
		Category:     common.PopupCategoryMeleeMissIn,
		Title:        "Melee Miss Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.MeleeMissIn,
		categoryName: "Melee",
		Color:        &cfg.MeleeMissInColor,
		Direction:    &cfg.MeleeMissInDirection,
		IsEnabled:    &cfg.MeleeMissInIsEnabled,
	}

	popupSettings[common.PopupCategorySpellCritOut] = &SettingProperty{
		Category:     common.PopupCategorySpellCritOut,
		Title:        "Spell Crit Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.SpellCritOut,
		categoryName: "Spell",
		Color:        &cfg.SpellCritOutColor,
		Direction:    &cfg.SpellCritOutDirection,
		IsEnabled:    &cfg.SpellCritOutIsEnabled,
	}

	popupSettings[common.PopupCategorySpellHitOut] = &SettingProperty{
		Category:     common.PopupCategorySpellHitOut,
		Title:        "Spell Hit Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.SpellHitOut,
		categoryName: "Spell",
		Color:        &cfg.SpellHitOutColor,
		Direction:    &cfg.SpellHitOutDirection,
		IsEnabled:    &cfg.SpellHitOutIsEnabled,
	}

	popupSettings[common.PopupCategorySpellMissOut] = &SettingProperty{
		Category:     common.PopupCategorySpellMissOut,
		Title:        "Spell Miss Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.SpellMissOut,
		categoryName: "Spell",
		Color:        &cfg.SpellMissOutColor,
		Direction:    &cfg.SpellMissOutDirection,
		IsEnabled:    &cfg.SpellMissOutIsEnabled,
	}

	popupSettings[common.PopupCategorySpellCritIn] = &SettingProperty{
		Category:     common.PopupCategorySpellCritIn,
		Title:        "Spell Crit Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.SpellCritIn,
		categoryName: "Spell",
		Color:        &cfg.SpellCritInColor,
		Direction:    &cfg.SpellCritInDirection,
		IsEnabled:    &cfg.SpellCritInIsEnabled,
	}

	popupSettings[common.PopupCategorySpellHitIn] = &SettingProperty{
		Category:     common.PopupCategorySpellHitIn,
		Title:        "Spell Hit Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.SpellHitIn,
		categoryName: "Spell",
		Color:        &cfg.SpellHitInColor,
		Direction:    &cfg.SpellHitInDirection,
		IsEnabled:    &cfg.SpellHitInIsEnabled,
	}

	popupSettings[common.PopupCategorySpellMissIn] = &SettingProperty{
		Category:     common.PopupCategorySpellMissIn,
		Title:        "Spell Miss Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.SpellMissIn,
		categoryName: "Spell",
		Color:        &cfg.SpellMissInColor,
		Direction:    &cfg.SpellMissInDirection,
		IsEnabled:    &cfg.SpellMissInIsEnabled,
	}

	popupSettings[common.PopupCategoryHealCritOut] = &SettingProperty{
		Category:     common.PopupCategoryHealCritOut,
		Title:        "Heal Crit Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.HealCritOut,
		categoryName: "Heal",
		Color:        &cfg.HealCritOutColor,
		Direction:    &cfg.HealCritOutDirection,
		IsEnabled:    &cfg.HealCritOutIsEnabled,
	}

	popupSettings[common.PopupCategoryHealHitOut] = &SettingProperty{
		Category:     common.PopupCategoryHealHitOut,
		Title:        "Heal Hit Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.HealHitOut,
		categoryName: "Heal",
		Color:        &cfg.HealHitOutColor,
		Direction:    &cfg.HealHitOutDirection,
		IsEnabled:    &cfg.HealHitOutIsEnabled,
	}

	popupSettings[common.PopupCategoryHealCritIn] = &SettingProperty{
		Category:     common.PopupCategoryHealCritIn,
		Title:        "Heal Crit Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.HealCritIn,
		categoryName: "Heal",
		Color:        &cfg.HealCritInColor,
		Direction:    &cfg.HealCritInDirection,
		IsEnabled:    &cfg.HealCritInIsEnabled,
	}

	popupSettings[common.PopupCategoryHealHitIn] = &SettingProperty{
		Category:     common.PopupCategoryHealHitIn,
		Title:        "Heal Hit Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.HealHitIn,
		categoryName: "Heal",
		Color:        &cfg.HealHitInColor,
		Direction:    &cfg.HealHitInDirection,
		IsEnabled:    &cfg.HealHitInIsEnabled,
	}

	popupSettings[common.PopupCategoryRuneHitOut] = &SettingProperty{
		Category:     common.PopupCategoryRuneHitOut,
		Title:        "Rune Hit Outgoing",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.RuneHitOut,
		categoryName: "Rune",
		Color:        &cfg.RuneHitOutColor,
		Direction:    &cfg.RuneHitOutDirection,
		IsEnabled:    &cfg.RuneHitOutIsEnabled,
	}

	popupSettings[common.PopupCategoryRuneHitIn] = &SettingProperty{
		Category:     common.PopupCategoryRuneHitIn,
		Title:        "Rune Hit Incoming",
		font:         font,
		fontBorder:   font,
		rect:         &cfg.RuneHitIn,
		categoryName: "Rune",
		Color:        &cfg.RuneHitInColor,
		Direction:    &cfg.RuneHitInDirection,
		IsEnabled:    &cfg.RuneHitInIsEnabled,
	}

	return nil
}

func SettingByCategory(category common.PopupCategory) *SettingProperty {
	popup, ok := popupSettings[category]
	if !ok {
		return nil
	}

	return popup
}

func SetSettingPositionByCategory(category common.PopupCategory, rect image.Rectangle) error {
	setting, ok := popupSettings[category]
	if !ok {
		return fmt.Errorf("setting not found for %d", category)
	}
	setting.rect.Min = rect.Min
	setting.rect.Max = rect.Max

	switch category {
	case common.PopupCategoryGlobalCritOut:
		subCats := []common.PopupCategory{common.PopupCategoryMeleeCritOut, common.PopupCategorySpellCritOut, common.PopupCategoryHealCritOut}
		for _, subCat := range subCats {
			subSetting, ok := popupSettings[subCat]
			if !ok {
				return fmt.Errorf("setting not found for %d", subCat)
			}
			subSetting.rect.Min = rect.Min
			subSetting.rect.Max = rect.Max
		}
	case common.PopupCategoryGlobalHitOut:
		subCats := []common.PopupCategory{common.PopupCategoryMeleeHitOut, common.PopupCategorySpellHitOut, common.PopupCategoryHealHitOut, common.PopupCategoryRuneHitOut}
		for _, subCat := range subCats {
			subSetting, ok := popupSettings[subCat]
			if !ok {
				return fmt.Errorf("setting not found for %d", subCat)
			}
			subSetting.rect.Min = rect.Min
			subSetting.rect.Max = rect.Max
		}
	case common.PopupCategoryGlobalMissOut:
		subCats := []common.PopupCategory{common.PopupCategoryMeleeMissOut, common.PopupCategorySpellMissOut}
		for _, subCat := range subCats {
			subSetting, ok := popupSettings[subCat]
			if !ok {
				return fmt.Errorf("setting not found for %d", subCat)
			}
			subSetting.rect.Min = rect.Min
			subSetting.rect.Max = rect.Max
		}
	case common.PopupCategoryGlobalCritIn:
		subCats := []common.PopupCategory{common.PopupCategoryMeleeCritIn, common.PopupCategorySpellCritIn, common.PopupCategoryHealCritIn}
		for _, subCat := range subCats {
			subSetting, ok := popupSettings[subCat]
			if !ok {
				return fmt.Errorf("setting not found for %d", subCat)
			}
			subSetting.rect.Min = rect.Min
			subSetting.rect.Max = rect.Max
		}
	case common.PopupCategoryGlobalHitIn:
		subCats := []common.PopupCategory{common.PopupCategoryMeleeHitIn, common.PopupCategorySpellHitIn, common.PopupCategoryHealHitIn, common.PopupCategoryRuneHitIn}
		for _, subCat := range subCats {
			subSetting, ok := popupSettings[subCat]
			if !ok {
				return fmt.Errorf("setting not found for %d", subCat)
			}
			subSetting.rect.Min = rect.Min
			subSetting.rect.Max = rect.Max
		}
	case common.PopupCategoryGlobalMissIn:
		subCats := []common.PopupCategory{common.PopupCategoryMeleeMissIn, common.PopupCategorySpellMissIn}
		for _, subCat := range subCats {
			subSetting, ok := popupSettings[subCat]
			if !ok {
				return fmt.Errorf("setting not found for %d", subCat)
			}
			subSetting.rect.Min = rect.Min
			subSetting.rect.Max = rect.Max
		}
	}

	return nil
}

func SetSettingColorByCategory(category common.PopupCategory, color color.RGBA) error {
	setting, ok := popupSettings[category]
	if !ok {
		return fmt.Errorf("setting not found for %d", category)
	}
	setting.Color.R = color.R
	setting.Color.G = color.G
	setting.Color.B = color.B
	setting.Color.A = color.A
	fmt.Println("setting color", category, color)
	return nil
}

func SettingPositionByCategory(category common.PopupCategory) (image.Rectangle, error) {
	mu.RLock()
	defer mu.RUnlock()
	setting, ok := popupSettings[category]
	if !ok {
		return image.Rectangle{}, fmt.Errorf("setting not found for %d", category)
	}
	rect := image.Rectangle{
		Min: setting.rect.Min,
		Max: setting.rect.Max,
	}
	return rect, nil
}

func SetSettingDirection(category common.PopupCategory, direction common.Direction) error {
	mu.Lock()
	defer mu.Unlock()
	setting, ok := popupSettings[category]
	if !ok {
		return fmt.Errorf("setting not found for %d", category)
	}

	*setting.Direction = direction
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
	for i := len(popups) - 1; i >= 0; i-- {
		popup := popups[i]

		popup.y += popup.vy * 0.9
		popup.x += popup.vx
		popup.life -= 1

		// Reverse velocity after reaching the hover point
		if popup.life < popup.maxLife*0.5 && popup.vy < 0 {
			popup.vy = -popup.vy * 0.2 // Slows down
		}

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
		op.GeoM.Translate(float64(popup.setting.rect.Min.X)+popup.x+float64(offset), float64(popup.setting.rect.Min.Y)+popup.y+float64(offset))
		op.ColorScale.ScaleWithColor(shadow)
		text.Draw(screen, popup.text, popup.setting.font, op)
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(popup.setting.rect.Min.X)+popup.x, float64(popup.setting.rect.Min.Y)+popup.y)
		op.ColorScale.ScaleWithColor(col)

		text.Draw(screen, popup.text, popup.setting.fontBorder, op)
	}
}

func Spawn(event *common.DamageEvent) error {

	setting := SettingByCategory(event.Category)
	if setting == nil {
		return fmt.Errorf("no setting found for %d", event.Category)
	}

	if !*setting.IsEnabled {
		return nil
	}

	spellColor := color.RGBA{255, 255, 255, 255}

	spellColor.R = setting.Color.R
	spellColor.G = setting.Color.G
	spellColor.B = setting.Color.B
	spellColor.A = setting.Color.A

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

	vx := float64(0)
	vy := float64(-0.5 - rand.Float64() - (rand.Float64() / 2))
	switch *setting.Direction {
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

	popup := &Popup{
		setting: setting,
		text:    event.Damage,
		x:       randomSpawnRange(int(setting.lastSpawnX), 0, setting.rect.Dx()-50, setting.rect.Dx()/4, 0),
		y:       randomSpawnRange(int(setting.lastSpawnX), 0, setting.rect.Dy()-50, setting.rect.Dy()/4, 0),
		vx:      vx,
		vy:      vy,
		life:    240,
		maxLife: 240,
		color:   spellColor,
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
	popups = append(popups, popup)
	return nil
}

// ConfigUpdate updates the popup configuration
func ConfigUpdate(cfg *config.CritSprinklerConfiguration) {
	mu.Lock()
	defer mu.Unlock()
}

// IsGlobalCategory returns true if the category is a global category
func IsGlobalCategory(category common.PopupCategory) bool {
	return category == common.PopupCategoryGlobalCritOut ||
		category == common.PopupCategoryGlobalHitOut ||
		category == common.PopupCategoryGlobalMissOut ||
		category == common.PopupCategoryGlobalCritIn ||
		category == common.PopupCategoryGlobalHitIn ||
		category == common.PopupCategoryGlobalMissIn
}

// GlobalCategoryToCategory returns a list of categories based on the global category
func GlobalCategoryToCategory(category common.PopupCategory) []common.PopupCategory {
	switch category {
	case common.PopupCategoryGlobalCritOut:
		return []common.PopupCategory{common.PopupCategoryMeleeCritOut, common.PopupCategorySpellCritOut, common.PopupCategoryHealCritOut}
	case common.PopupCategoryGlobalHitOut:
		return []common.PopupCategory{common.PopupCategoryMeleeHitOut, common.PopupCategorySpellHitOut, common.PopupCategoryHealHitOut, common.PopupCategoryRuneHitOut}
	case common.PopupCategoryGlobalMissOut:
		return []common.PopupCategory{common.PopupCategoryMeleeMissOut, common.PopupCategorySpellMissOut}
	case common.PopupCategoryGlobalCritIn:
		return []common.PopupCategory{common.PopupCategoryMeleeCritIn, common.PopupCategorySpellCritIn, common.PopupCategoryHealCritIn}
	case common.PopupCategoryGlobalHitIn:
		return []common.PopupCategory{common.PopupCategoryMeleeHitIn, common.PopupCategorySpellHitIn, common.PopupCategoryHealHitIn, common.PopupCategoryRuneHitIn}
	case common.PopupCategoryGlobalMissIn:
		return []common.PopupCategory{common.PopupCategoryMeleeMissIn, common.PopupCategorySpellMissIn}
	}
	return nil
}
