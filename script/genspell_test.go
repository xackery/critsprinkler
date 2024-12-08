package script

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGenSpell(t *testing.T) {

	err := Spell("c:/games/eq/thj")
	if err != nil {
		t.Fatalf("Spell: %v", err)
	}
	outputPath := "../spell_colors.txt"
	err = os.Rename(outputPath, strings.Replace(outputPath, ".txt", ".go", 1))
	if err != nil {
		t.Fatalf("Rename: %v", err)
	}
}

func Spell(path string) error {
	start := time.Now()
	defer func() {
		fmt.Printf("Finished in %0.2f seconds\n", time.Since(start).Seconds())
	}()

	data, err := os.ReadFile(path + "/spells_us.txt")
	if err != nil {
		return fmt.Errorf("read %s: %w", path+"/spells_us.txt", err)
	}

	lines := strings.Split(string(data), "\n")
	outputPath := "../spell_colors.txt"
	w, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", outputPath, err)
	}
	defer w.Close()

	w.WriteString(`package main

import "image/color"

var (
	spellColors = map[string]color.RGBA{`)

	resistColors := map[string]color.RGBA{
		"Magic":      {186, 143, 206, 255},
		"Fire":       {230, 126, 34, 255},
		"Cold":       {0, 244, 255, 255},
		"Poison":     {152, 250, 60, 255},
		"Disease":    {195, 254, 10, 255},
		"Chromatic":  {128, 0, 128, 255},
		"Prismatic":  {255, 255, 255, 255},
		"Physical":   {128, 128, 128, 255},
		"Corruption": {63, 63, 63, 255},
	}

	names := map[string]bool{}

	lineNumber := 0
	columnCount := 0
	for _, line := range lines {
		lineNumber++
		records := strings.Split(line, "^")
		if lineNumber == 1 {
			columnCount = len(records)
			fmt.Printf("Number of columns: %d\n", columnCount)
		}
		if len(records) != columnCount {
			continue
			//return fmt.Errorf("line %d: expected %d columns, got %d", lineNumber, columnCount, len(records))
		}

		spellName := records[1]

		resistType := "None"
		switch records[85] {
		case "0":
			continue
		case "1":
			resistType = "Magic"
		case "2":
			resistType = "Fire"
		case "3":
			resistType = "Cold"
		case "4":
			resistType = "Poison"
		case "5":
			resistType = "Disease"
		case "6":
			resistType = "Chromatic"
		case "7":
			resistType = "Prismatic"
		case "8":
			resistType = "Physical"
		case "9":
			resistType = "Corruption"
		default:
			continue
		}

		resistCode, ok := resistColors[resistType]
		if !ok {
			return fmt.Errorf("line %d: %s unknown resist type %s", lineNumber, spellName, resistType)
		}

		_, ok = names[spellName]
		if ok {
			continue
		}

		w.WriteString(fmt.Sprintf(`
					"%s": {%d, %d, %d, %d},`, spellName, resistCode.R, resistCode.G, resistCode.B, resistCode.A))

		names[spellName] = true
	}

	w.WriteString(`
		}
)`)

	fmt.Printf("Data written to %s\n", outputPath)

	return nil
}
