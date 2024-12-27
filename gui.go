package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/config"
	"github.com/xackery/critsprinkler/popup"
	"github.com/xackery/critsprinkler/tracker"
)

var (
	lastGlobalTestCategory common.PopupCategory
)

// NewGUI creates a new CritSprinkler GUI
func NewGUI(cfg *config.CritSprinklerConfiguration, game *Game) (*ebitenui.UI, error) {

	res, err := newUIResources()
	if err != nil {
		return nil, fmt.Errorf("new ui resources: %w", err)
	}

	game.res = res

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	toolbar := newToolbar(eui, res)
	rootContainer.AddChild(toolbar.container)

	toolbar.btnGlobalTest.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			lastGlobalTestCategory++
			if lastGlobalTestCategory >= common.PopupCategoryGlobalCritOut {
				lastGlobalTestCategory = 0
			}

			setting := popup.SettingByCategory(lastGlobalTestCategory)
			if setting == nil {
				game.statusText = fmt.Sprintf("Failed: No settings found for %d", lastGlobalTestCategory)
				return
			}

			onDamageEvent(&common.DamageEvent{
				Category:  common.PopupCategory(lastGlobalTestCategory),
				Source:    tracker.PlayerName(),
				Target:    "Test",
				SpellName: "Ice Comet",
				Damage:    setting.Title,
			})
		}),

		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Test crit effects on screen" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnLoadEQLog.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			err := FileDialogBox(cfg.LogPath)
			if err != nil {
				MsgBox("Error", fmt.Sprintf("Error loading log: %v", err))
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
			if cfg.LogPath == "" {
				game.statusText = "Load EQ Log file (Currently not set)"
			} else {
				logPath := filepath.Base(cfg.LogPath)
				game.statusText = "Load EQ Log file (Currently " + logPath + ")"
			}
		}),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnSave.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			err := updateSave()
			if err != nil {
				MsgBox("Error", fmt.Sprintf("Error saving config: %v", err))
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Save changes to " + config.FileName() }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnQuit.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			os.Exit(0)
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Quit CritSprinkler" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnGlobalCritOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryGlobalCritOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Global Crits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Global Crits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnGlobalCritIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryGlobalCritIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Global Crits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Global Crits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnGlobalHitOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryGlobalHitOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Global Hits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Global Hits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnGlobalHitIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryGlobalHitIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Global Hits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Global Hits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnGlobalMissOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryGlobalMissOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Global Misses"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Global Misses" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnGlobalMissIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryGlobalMissIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Global Misses"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Global Misses" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnMeleeCritOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryMeleeCritOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Melee Crits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Melee Crits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnMeleeCritIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryMeleeCritIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Melee Crits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Melee Crits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnMeleeHitOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryMeleeHitOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Melee Hits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Melee Hits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnMeleeHitIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryMeleeHitIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Melee Hits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Melee Hits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnMeleeMissOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryMeleeMissOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Melee Misses"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Melee Misses" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnMeleeMissIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryMeleeMissIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Melee Misses"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Melee Misses" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnSpellCritOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategorySpellCritOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Spell Crits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Spell Crits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnSpellCritIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategorySpellCritIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Spell Crits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Spell Crits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnSpellHitOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategorySpellHitOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Spell Hits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Spell Hits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnSpellHitIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategorySpellHitIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Spell Hits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Spell Hits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnSpellMissOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategorySpellMissOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Spell Misses"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Spell Misses" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnSpellMissIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategorySpellMissIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Spell Misses"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Spell Misses" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnHealCritOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryHealCritOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Heal Crits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Heal Crits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnHealCritIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryHealCritIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Heal Crits"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Heal Crits" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnHealHitOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryHealHitOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Heals"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Heals" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnHealHitIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryHealHitIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Heals"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Heals" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnRuneHitOut.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryRuneHitOut)
			if setting == nil {
				game.statusText = "Failed: No settings found for Outgoing Runes"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Outgoing Runes" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	toolbar.btnRuneHitIn.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			setting := popup.SettingByCategory(common.PopupCategoryRuneHitIn)
			if setting == nil {
				game.statusText = "Failed: No settings found for Incoming Runes"
				return
			}
			err = openPopupSettingsWindow(res, eui, setting)
			if err != nil {
				game.statusText = fmt.Sprintf("Error opening popup settings: %v", err)
			}
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "Configure Incoming Runes" }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { game.statusText = "" }),
	)

	return eui, nil
}

func (g *Game) statusBarDraw(screen *ebiten.Image) {
	//x, y := ebiten.WindowPosition()
	screenWidth, screenHeight := ebiten.WindowSize()

	vector.DrawFilledRect(screen, 0, float32(screenHeight-40), float32(screenWidth), float32(screenHeight), color.RGBA{0, 0, 0, 190}, true)
	op := &text.DrawOptions{}
	op.GeoM.Translate(10, float64(screenHeight-35))
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, g.statusText, g.res.button.face, op)
}
