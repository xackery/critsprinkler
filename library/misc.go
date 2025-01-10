package library

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"time"

	"github.com/dblezek/tga"
	"github.com/hajimehoshi/ebiten/v2"
)

type Misc int

const (
	MiscPlatinum Misc = iota
	MiscGold
	MiscSilver
	MiscCopper
	MiscFavor
)

var (
	miscs = make(map[Misc]*ebiten.Image)
)

// MiscLoad loads all miscs
func MiscLoad(eqPath string) error {
	var err error
	start := time.Now()

	err = miscCoins(eqPath)
	if err != nil {
		return fmt.Errorf("coins: %w", err)
	}

	fmt.Printf("loaded misc assets in %0.2fs\n", time.Since(start).Seconds())
	return nil
}

// MiscByID returns an misc by misc id
func MiscByID(id Misc) *ebiten.Image {
	return miscs[id]
}

func miscCoins(eqPath string) error {
	fileName := "window_pieces01.tga"
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

	elements := []struct {
		id   Misc
		pos  image.Point
		size image.Point
	}{
		{MiscPlatinum, image.Point{X: 90, Y: 140}, image.Point{X: 18, Y: 18}},
		{MiscGold, image.Point{X: 108, Y: 140}, image.Point{X: 18, Y: 18}},
		{MiscSilver, image.Point{X: 126, Y: 140}, image.Point{X: 18, Y: 18}},
		{MiscCopper, image.Point{X: 144, Y: 140}, image.Point{X: 18, Y: 18}},
		{MiscFavor, image.Point{X: 88, Y: 159}, image.Point{X: 20, Y: 20}},
	}

	for _, element := range elements {
		dst := image.NewRGBA(image.Rect(0, 0, element.size.X, element.size.Y))
		draw.Draw(dst, dst.Bounds(), nrgba.SubImage(image.Rect(element.pos.X, element.pos.Y, element.pos.X+element.size.X, element.pos.Y+element.size.Y)), image.Point{X: element.pos.X, Y: element.pos.Y}, draw.Src)

		miscs[element.id] = ebiten.NewImageFromImage(dst)
	}

	return nil
}
