package menu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/xackery/critsprinkler/config"
	"github.com/xackery/critsprinkler/library"
	"github.com/xackery/critsprinkler/spell"
	"github.com/xackery/critsprinkler/status"
)

var (
	onEQPathLoad func()
)

// New creates a new UI
func New(cfg *config.CritSprinklerConfiguration, onSave func(), onEQPathLoadSrc func()) (*ebitenui.UI, error) {

	onEQPathLoad = onEQPathLoadSrc
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	eui := &ebitenui.UI{
		Container: rootContainer,
	}

	toolbar, err := toolbarNew(cfg, eui)
	if err != nil {
		return nil, fmt.Errorf("toolbarNew: %w", err)
	}
	rootContainer.AddChild(toolbar.container)

	toolbar.btnFileSave.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			onSave()
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { status.Setf("Save changes to %s", config.FileName()) }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { status.Set("") }),
	)

	toolbar.btnFileQuit.Configure(
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			os.Exit(0)
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) { status.Set("Quit CritSprinkler") }),
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) { status.Set("") }),
	)

	err = eqPathLoad(cfg)
	if err != nil {
		return nil, fmt.Errorf("eqPathLoad: %w", err)
	}
	return eui, nil

}

// SetVisible sets the visibility of the UI
func SetEditMode(ui *ebitenui.UI, isEditMode bool) {
	state := widget.Visibility_Show
	if !isEditMode {
		state = widget.Visibility_Hide
	}
	ui.Container.GetWidget().Visibility = state
}

func eqPathLoad(cfg *config.CritSprinklerConfiguration) error {

	path := cfg.EQPath
	if path == "" {
		return nil
	}

	err := spell.Load(filepath.Join(path, "spells_us.txt"))
	if err != nil {
		return fmt.Errorf("spell load: %w", err)
	}

	err = library.SpellLoad(path)
	if err != nil {
		return fmt.Errorf("spell load: %w", err)
	}

	err = library.MiscLoad(path)
	if err != nil {
		return fmt.Errorf("misc load: %w", err)
	}
	onEQPathLoad()
	return nil
}
