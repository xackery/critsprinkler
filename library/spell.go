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

	"github.com/dblezek/tga"

	"github.com/hajimehoshi/ebiten/v2"

	eimage "github.com/ebitenui/ebitenui/image"
)

var (
	spells         = make(map[int]*ebiten.Image)
	spellNineSlice = make(map[int]*eimage.NineSlice)
)

// SpellLoad loads all spells
func SpellLoad(eqPath string) error {
	start := time.Now()
	spellSpellSize := image.Point{X: 40, Y: 40}
	spellSpellPadding := image.Point{X: 0, Y: 0}
	spellRect := image.Rect(0, 0, spellSpellSize.X, spellSpellSize.Y)
	totalSpellSpellCount := 0
	for i := 1; i < 8; i++ {
		fileName := fmt.Sprintf("spells%02d.tga", i)
		path := filepath.Join(eqPath, fmt.Sprintf("uifiles/default/%s", fileName))
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", fileName, err)
		}

		img, err := tga.Decode(bytes.NewReader(data))
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

		spellToDump := -1
		//spellIndex := 0
		for y := 0; y < nrgba.Bounds().Dy(); y += spellSpellSize.Y + spellSpellPadding.Y {
			if y+spellSpellSize.Y > nrgba.Bounds().Dy() {
				break
			}
			for x := 0; x < nrgba.Bounds().Dx(); x += spellSpellSize.X + spellSpellPadding.X {
				if x+spellSpellSize.X > nrgba.Bounds().Dx() {
					break
				}

				spell := image.NewRGBA(spellRect)
				draw.Draw(spell, spellRect, nrgba.SubImage(image.Rect(x, y, x+spellSpellSize.X, y+spellSpellSize.Y)), image.Point{x, y}, draw.Src)
				//fmt.Println("spell", totalSpellSpellCount, "x", x, "y", y, "dx", x+spellSpellSize.X, "dy", y+spellSpellSize.Y)
				// spellIndex++

				if spellToDump == totalSpellSpellCount {
					w, err := os.Create(fmt.Sprintf("spell_%02d.png", totalSpellSpellCount))
					if err != nil {
						return fmt.Errorf("create spell: %w", err)
					}
					defer w.Close()

					err = png.Encode(w, spell)
					if err != nil {
						return fmt.Errorf("encode spell: %w", err)
					}
				}

				spells[totalSpellSpellCount] = ebiten.NewImageFromImage(spell)

				totalSpellSpellCount++
			}
		}
	}
	fmt.Printf("loaded %d spells in %0.2fs\n", totalSpellSpellCount, time.Since(start).Seconds())
	return nil
}

// SpellByID returns an spell by spell id
func SpellByID(id int) *ebiten.Image {
	return spells[id]
}

// SpellByIDNineSlice returns a nine slice of an spell by spell id
func SpellByIDNineSlice(id int) *eimage.NineSlice {
	e, ok := spellNineSlice[id]
	if !ok {
		// generate a nine slice cache
		e = eimage.NewNineSlice(SpellByID(id), [3]int{20, 0, 20}, [3]int{20, 0, 20})
		spellNineSlice[id] = e
	}
	return e
}
