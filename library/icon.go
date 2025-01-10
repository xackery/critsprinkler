package library

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/malashin/dds"
)

var (
	icons = make(map[int]*ebiten.Image)
)

// IconLoad loads all icons
func IconLoad(eqPath string) error {
	start := time.Now()
	spellIconSize := image.Point{X: 40, Y: 40}
	spellIconPadding := image.Point{X: 0, Y: 0}
	iconRect := image.Rect(0, 0, spellIconSize.X, spellIconSize.Y)
	totalSpellIconCount := 500
	for i := 1; i < 179; i++ {
		fileName := fmt.Sprintf("dragitem%d.dds", i)
		path := filepath.Join(eqPath, fmt.Sprintf("uifiles/default/%s", fileName))
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", fileName, err)
		}

		img, err := dds.Decode(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("decode %s: %w", fileName, err)
		}

		var nrgba *image.NRGBA
		switch val := img.(type) {
		case *image.NRGBA:
			nrgba = val
		default:
			return fmt.Errorf("unknown type: %T", val)
		}

		spellToDump := 590
		//iconIndex := 0
		for x := 0; x < nrgba.Bounds().Dx(); x += spellIconSize.X + spellIconPadding.X {
			if x+spellIconSize.X > nrgba.Bounds().Dx() {
				break
			}
			for y := 0; y < nrgba.Bounds().Dy(); y += spellIconSize.Y + spellIconPadding.Y {

				if y+spellIconSize.Y > nrgba.Bounds().Dy() {
					break
				}

				icon := image.NewRGBA(iconRect)
				draw.Draw(icon, iconRect, nrgba.SubImage(image.Rect(x, y, x+spellIconSize.X, y+spellIconSize.Y)), image.Point{x, y}, draw.Src)
				//fmt.Println("icon", totalSpellIconCount, "x", x, "y", y, "dx", x+spellIconSize.X, "dy", y+spellIconSize.Y)
				// iconIndex++

				if spellToDump == totalSpellIconCount {
					w, err := os.Create(fmt.Sprintf("icon_%d.png", totalSpellIconCount))
					if err != nil {
						return fmt.Errorf("create icon: %w", err)
					}
					defer w.Close()

					err = png.Encode(w, icon)
					if err != nil {
						return fmt.Errorf("encode icon: %w", err)
					}
				}

				icons[totalSpellIconCount] = ebiten.NewImageFromImage(icon)

				totalSpellIconCount++
			}
		}
	}
	fmt.Printf("loaded %d icons in %0.2fs\n", totalSpellIconCount, time.Since(start).Seconds())
	return nil
}

// IconByID returns an icon by icon id
func IconByID(id int) *ebiten.Image {
	return icons[id]
}
