package config

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/critsprinkler/common"
)

const (
	fileName = "critsprinkler.ini"
)

var (
	mu sync.RWMutex
)

type CritSprinklerConfiguration struct {
	IsNew      bool
	LogPath    string          `config:"log_path"`
	MainWindow image.Rectangle `config:"main_window"`

	MeleeHitOut  image.Rectangle `config:"melee_hit_out"`
	MeleeHitIn   image.Rectangle `config:"melee_hit_in"`
	MeleeCritOut image.Rectangle `config:"melee_crit_out"`
	MeleeCritIn  image.Rectangle `config:"melee_crit_in"`
	MeleeMissOut image.Rectangle `config:"melee_miss_out"`
	MeleeMissIn  image.Rectangle `config:"melee_miss_in"`

	GlobalHitOut  image.Rectangle `config:"global_hit_out"`
	GlobalHitIn   image.Rectangle `config:"global_hit_in"`
	GlobalCritOut image.Rectangle `config:"global_crit_out"`
	GlobalCritIn  image.Rectangle `config:"global_crit_in"`
	GlobalMissOut image.Rectangle `config:"global_miss_out"`
	GlobalMissIn  image.Rectangle `config:"global_miss_in"`

	SpellHitOut  image.Rectangle `config:"spell_hit_out"`
	SpellHitIn   image.Rectangle `config:"spell_hit_in"`
	SpellCritOut image.Rectangle `config:"spell_crit_out"`
	SpellCritIn  image.Rectangle `config:"spell_crit_in"`
	SpellMissOut image.Rectangle `config:"spell_miss_out"`
	SpellMissIn  image.Rectangle `config:"spell_miss_in"`

	HealHitOut  image.Rectangle `config:"heal_hit_out"`
	HealHitIn   image.Rectangle `config:"heal_hit_in"`
	HealCritOut image.Rectangle `config:"heal_crit_out"`
	HealCritIn  image.Rectangle `config:"heal_crit_in"`

	RuneHitOut image.Rectangle `config:"rune_hit_out"`
	RuneHitIn  image.Rectangle `config:"rune_hit_in"`

	MeleeHitOutColor   color.RGBA `config:"melee_hit_out_color"`
	MeleeHitInColor    color.RGBA `config:"melee_hit_in_color"`
	MeleeCritOutColor  color.RGBA `config:"melee_crit_out_color"`
	MeleeCritInColor   color.RGBA `config:"melee_crit_in_color"`
	MeleeMissOutColor  color.RGBA `config:"melee_miss_out_color"`
	MeleeMissInColor   color.RGBA `config:"melee_miss_in_color"`
	GlobalHitOutColor  color.RGBA `config:"global_hit_out_color"`
	GlobalHitInColor   color.RGBA `config:"global_hit_in_color"`
	GlobalCritOutColor color.RGBA `config:"global_crit_out_color"`
	GlobalCritInColor  color.RGBA `config:"global_crit_in_color"`
	GlobalMissOutColor color.RGBA `config:"global_miss_out_color"`
	GlobalMissInColor  color.RGBA `config:"global_miss_in_color"`
	SpellHitOutColor   color.RGBA `config:"spell_hit_out_color"`
	SpellHitInColor    color.RGBA `config:"spell_hit_in_color"`
	SpellCritOutColor  color.RGBA `config:"spell_crit_out_color"`
	SpellCritInColor   color.RGBA `config:"spell_crit_in_color"`
	SpellMissOutColor  color.RGBA `config:"spell_miss_out_color"`
	SpellMissInColor   color.RGBA `config:"spell_miss_in_color"`
	HealHitOutColor    color.RGBA `config:"heal_hit_out_color"`
	HealHitInColor     color.RGBA `config:"heal_hit_in_color"`
	HealCritOutColor   color.RGBA `config:"heal_crit_out_color"`
	HealCritInColor    color.RGBA `config:"heal_crit_in_color"`
	RuneHitOutColor    color.RGBA `config:"rune_hit_out_color"`
	RuneHitInColor     color.RGBA `config:"rune_hit_in_color"`

	MeleeHitOutDirection   common.Direction `config:"melee_hit_out_direction"`
	MeleeHitInDirection    common.Direction `config:"melee_hit_in_direction"`
	MeleeCritOutDirection  common.Direction `config:"melee_crit_out_direction"`
	MeleeCritInDirection   common.Direction `config:"melee_crit_in_direction"`
	MeleeMissOutDirection  common.Direction `config:"melee_miss_out_direction"`
	MeleeMissInDirection   common.Direction `config:"melee_miss_in_direction"`
	GlobalHitOutDirection  common.Direction `config:"global_hit_out_direction"`
	GlobalHitInDirection   common.Direction `config:"global_hit_in_direction"`
	GlobalCritOutDirection common.Direction `config:"global_crit_out_direction"`
	GlobalCritInDirection  common.Direction `config:"global_crit_in_direction"`
	GlobalMissOutDirection common.Direction `config:"global_miss_out_direction"`
	GlobalMissInDirection  common.Direction `config:"global_miss_in_direction"`
	SpellHitOutDirection   common.Direction `config:"spell_hit_out_direction"`
	SpellHitInDirection    common.Direction `config:"spell_hit_in_direction"`
	SpellCritOutDirection  common.Direction `config:"spell_crit_out_direction"`
	SpellCritInDirection   common.Direction `config:"spell_crit_in_direction"`
	SpellMissOutDirection  common.Direction `config:"spell_miss_out_direction"`
	SpellMissInDirection   common.Direction `config:"spell_miss_in_direction"`
	HealHitOutDirection    common.Direction `config:"heal_hit_out_direction"`
	HealHitInDirection     common.Direction `config:"heal_hit_in_direction"`
	HealCritOutDirection   common.Direction `config:"heal_crit_out_direction"`
	HealCritInDirection    common.Direction `config:"heal_crit_in_direction"`
	RuneHitOutDirection    common.Direction `config:"rune_hit_out_direction"`
	RuneHitInDirection     common.Direction `config:"rune_hit_in_direction"`

	MeleeHitOutIsEnabled  bool `config:"melee_hit_out_is_enabled"`
	MeleeHitInIsEnabled   bool `config:"melee_hit_in_is_enabled"`
	MeleeCritOutIsEnabled bool `config:"melee_crit_out_is_enabled"`
	MeleeCritInIsEnabled  bool `config:"melee_crit_in_is_enabled"`
	MeleeMissOutIsEnabled bool `config:"melee_miss_out_is_enabled"`
	MeleeMissInIsEnabled  bool `config:"melee_miss_in_is_enabled"`
	SpellHitOutIsEnabled  bool `config:"spell_hit_out_is_enabled"`
	SpellHitInIsEnabled   bool `config:"spell_hit_in_is_enabled"`
	SpellCritOutIsEnabled bool `config:"spell_crit_out_is_enabled"`
	SpellCritInIsEnabled  bool `config:"spell_crit_in_is_enabled"`
	SpellMissOutIsEnabled bool `config:"spell_miss_out_is_enabled"`
	SpellMissInIsEnabled  bool `config:"spell_miss_in_is_enabled"`
	HealHitOutIsEnabled   bool `config:"heal_hit_out_is_enabled"`
	HealHitInIsEnabled    bool `config:"heal_hit_in_is_enabled"`
	HealCritOutIsEnabled  bool `config:"heal_crit_out_is_enabled"`
	HealCritInIsEnabled   bool `config:"heal_crit_in_is_enabled"`
	RuneHitOutIsEnabled   bool `config:"rune_hit_out_is_enabled"`
	RuneHitInIsEnabled    bool `config:"rune_hit_in_is_enabled"`
}

// FileName returns the config file name
func FileName() string {
	return fileName
}

// LoadCritSprinklerConfig loads an CritSprinkler config file
func LoadCritSprinklerConfig() (*CritSprinklerConfiguration, error) {
	mu.Lock()
	defer mu.Unlock()
	knownKeys := []string{}
	exePath := os.Args[0]
	path := filepath.Dir(exePath)
	path = filepath.Join(path, fileName)

	_, err := os.Stat(path)
	if err != nil {
		x, y := ebiten.Monitor().Size()
		width := 1000
		height := 700

		x = (x - width) / 2
		y = (y - height) / 2

		mainWindowRect := image.Rectangle{Min: image.Point{X: x, Y: y}, Max: image.Point{X: x + width, Y: y + height}}
		defaultRect := image.Rectangle{Min: image.Point{X: 220, Y: 307}, Max: image.Point{X: 420, Y: 407}}
		return &CritSprinklerConfiguration{
			IsNew:                 true,
			MainWindow:            mainWindowRect,
			MeleeHitOut:           defaultRect,
			MeleeHitOutColor:      color.RGBA{R: 255, G: 255, B: 255, A: 200},
			MeleeHitOutDirection:  common.DirectionDown,
			MeleeHitOutIsEnabled:  true,
			MeleeHitIn:            defaultRect,
			MeleeHitInColor:       color.RGBA{R: 255, G: 0, B: 0, A: 200},
			MeleeHitInDirection:   common.DirectionLeft,
			MeleeHitInIsEnabled:   true,
			MeleeCritOut:          defaultRect,
			MeleeCritOutColor:     color.RGBA{R: 255, G: 165, B: 0, A: 200}, // Orange color
			MeleeCritOutDirection: common.DirectionUp,
			MeleeCritOutIsEnabled: true,
			MeleeCritIn:           defaultRect,
			MeleeCritInColor:      color.RGBA{R: 255, G: 0, B: 0, A: 200},
			MeleeCritInDirection:  common.DirectionUpLeft,
			MeleeCritInIsEnabled:  true,
			MeleeMissOut:          defaultRect,
			MeleeMissOutColor:     color.RGBA{R: 255, G: 0, B: 0, A: 200},
			MeleeMissOutDirection: common.DirectionDownLeft,
			MeleeMissOutIsEnabled: false,
			MeleeMissIn:           defaultRect,
			MeleeMissInColor:      color.RGBA{R: 255, G: 255, B: 255, A: 200},
			MeleeMissInDirection:  common.DirectionDownRight,
			MeleeMissInIsEnabled:  false,

			GlobalHitOut:  defaultRect,
			GlobalHitIn:   defaultRect,
			GlobalCritOut: defaultRect,
			GlobalCritIn:  defaultRect,
			GlobalMissOut: defaultRect,
			GlobalMissIn:  defaultRect,

			SpellHitOut:           defaultRect,
			SpellHitOutColor:      color.RGBA{R: 0, G: 0, B: 255, A: 200}, // Blue color
			SpellHitOutDirection:  common.DirectionDown,
			SpellHitOutIsEnabled:  true,
			SpellHitIn:            defaultRect,
			SpellHitInColor:       color.RGBA{R: 255, G: 0, B: 255, A: 200}, // Purple
			SpellHitInDirection:   common.DirectionRight,
			SpellHitInIsEnabled:   true,
			SpellCritOut:          defaultRect,
			SpellCritOutColor:     color.RGBA{R: 255, G: 69, B: 0, A: 200}, // Red-Orange
			SpellCritOutDirection: common.DirectionUp,
			SpellCritOutIsEnabled: true,
			SpellCritIn:           defaultRect,
			SpellCritInColor:      color.RGBA{R: 255, G: 0, B: 0, A: 200}, // Red
			SpellCritInDirection:  common.DirectionDownRight,
			SpellCritInIsEnabled:  true,
			SpellMissOut:          defaultRect,
			SpellMissOutColor:     color.RGBA{R: 0, G: 0, B: 255, A: 200}, // blue
			SpellMissOutDirection: common.DirectionDownRight,
			SpellMissOutIsEnabled: true,
			SpellMissIn:           defaultRect,
			SpellMissInColor:      color.RGBA{R: 255, G: 0, B: 255, A: 200}, // pink
			SpellMissInDirection:  common.DirectionDownRight,
			SpellMissInIsEnabled:  true,

			HealHitOut:           defaultRect,
			HealHitOutColor:      color.RGBA{R: 0, G: 255, B: 0, A: 200}, // Green
			HealHitOutDirection:  common.DirectionDown,
			HealHitOutIsEnabled:  true,
			HealHitIn:            defaultRect,
			HealHitInColor:       color.RGBA{R: 0, G: 128, B: 128, A: 200}, //teal
			HealHitInDirection:   common.DirectionDown,
			HealHitInIsEnabled:   true,
			HealCritOut:          defaultRect,
			HealCritOutColor:     color.RGBA{R: 0, G: 100, B: 0, A: 200}, //dark green
			HealCritOutDirection: common.DirectionUp,
			HealCritOutIsEnabled: true,
			HealCritIn:           defaultRect,
			HealCritInColor:      color.RGBA{R: 0, G: 128, B: 128, A: 200}, // dark teal
			HealCritInDirection:  common.DirectionUp,
			HealCritInIsEnabled:  true,

			RuneHitOut:          defaultRect,
			RuneHitOutColor:     color.RGBA{R: 128, G: 128, B: 128, A: 200},
			RuneHitOutDirection: common.DirectionDown,
			RuneHitOutIsEnabled: true,
			RuneHitIn:           defaultRect,
			RuneHitInColor:      color.RGBA{R: 128, G: 128, B: 128, A: 200},
			RuneHitInDirection:  common.DirectionDown,
			RuneHitInIsEnabled:  true,

			LogPath: "",
		}, nil
	}

	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %s", strings.TrimPrefix(err.Error(), fmt.Sprintf("open %s: ", fileName)))
	}
	defer r.Close()

	var config CritSprinklerConfiguration

	foundKeys := map[string]bool{}

	for i := range reflect.TypeOf(config).NumField() {
		sKey, ok := reflect.StructTag(reflect.TypeOf(config).Field(i).Tag).Lookup("config")
		if !ok {
			continue
		}

		knownKeys = append(knownKeys, sKey)
	}

	lineNumber := 0
	reader := bufio.NewScanner(r)
	for reader.Scan() {
		line := reader.Text()
		lineNumber++
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				continue
			}
			key := strings.ToLower(strings.TrimSpace(parts[0]))
			value := strings.TrimSpace(parts[1])
			// reflect struct tags of cfg for config and see if they match key

			isKnown := false
			for i := range reflect.TypeOf(config).NumField() {
				sKey, ok := reflect.StructTag(reflect.TypeOf(config).Field(i).Tag).Lookup("config")
				if !ok {
					continue
				}

				if sKey != key {
					continue
				}

				field := reflect.ValueOf(&config).Elem().Field(i)
				switch field.Kind() {
				case reflect.Int:
					val, err := strconv.Atoi(value)
					if err != nil {
						return nil, fmt.Errorf("line %d parse %s=%s to int: %w", lineNumber, key, value, err)
					}

					field.SetInt(int64(val))
				case reflect.String:
					field.SetString(value)
				case reflect.Bool:
					val, err := strconv.ParseBool(value)
					if err != nil {
						return nil, fmt.Errorf("line %d parse %s=%s to bool: %w", lineNumber, key, value, err)
					}

					field.SetBool(val)
				case reflect.Struct:
					switch field.Interface().(type) {
					case image.Rectangle:
						parts := strings.Split(value, ",")
						if len(parts) != 4 {
							return nil, fmt.Errorf("line %d parse %s=%s to image.Rectangle: invalid number of parts", lineNumber, key, value)
						}

						var rect image.Rectangle
						for i := range parts {
							val, err := strconv.Atoi(parts[i])
							if err != nil {
								return nil, fmt.Errorf("line %d parse %s=%s to image.Rectangle: %w", lineNumber, key, value, err)
							}

							switch i {
							case 0:
								rect.Min.X = val
							case 1:
								rect.Min.Y = val
							case 2:
								rect.Max.X = val
							case 3:
								rect.Max.Y = val
							}
						}

						field.Set(reflect.ValueOf(rect))
					case color.RGBA:
						parts := strings.Split(value, ",")
						if len(parts) != 4 {
							return nil, fmt.Errorf("line %d parse %s=%s to color.RGBA: invalid number of parts", lineNumber, key, value)
						}

						var rgba color.RGBA
						for i := range parts {
							val, err := strconv.Atoi(parts[i])
							if err != nil {
								return nil, fmt.Errorf("line %d parse %s=%s to color.RGBA: %w", lineNumber, key, value, err)
							}

							switch i {
							case 0:
								rgba.R = uint8(val)
							case 1:
								rgba.G = uint8(val)
							case 2:
								rgba.B = uint8(val)
							case 3:
								rgba.A = uint8(val)
							}
						}

						field.Set(reflect.ValueOf(rgba))
					case common.Direction:
						val, err := strconv.Atoi(value)
						if err != nil {
							return nil, fmt.Errorf("line %d parse %s=%s to common.Direction: %w", lineNumber, key, value, err)
						}
						field.Set(reflect.ValueOf(common.Direction(val)))

					default:
						return nil, fmt.Errorf("line %d unknown struct type %s", lineNumber, field.Kind())
					}
				default:
					return nil, fmt.Errorf("line %d unknown type %s", lineNumber, field.Kind())
				}
				isKnown = true
				foundKeys[key] = true
			}

			if !isKnown {
				return nil, fmt.Errorf("line %d unknown key %s", lineNumber, key)
			}
		}
	}

	if len(foundKeys) != len(knownKeys) {
		for i := range knownKeys {
			if !foundKeys[knownKeys[i]] {
				return nil, fmt.Errorf("missing key %s", knownKeys[i])
			}
		}
	}

	return &config, nil
}

// Save saves the config
func (c *CritSprinklerConfiguration) Save() error {
	exePath := os.Args[0]
	path := filepath.Dir(exePath)
	path = filepath.Join(path, fileName)

	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat %s: %w", fileName, err)
		}
		w, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("create %s: %w", fileName, err)
		}
		w.Close()
	}
	if fi != nil && fi.IsDir() {
		return fmt.Errorf("%s is a directory", fileName)
	}

	out := ""

	for i := range reflect.TypeOf(*c).NumField() {
		sKey, ok := reflect.StructTag(reflect.TypeOf(*c).Field(i).Tag).Lookup("config")
		if !ok {
			continue
		}

		field := reflect.ValueOf(c).Elem().Field(i)
		switch field.Kind() {
		case reflect.Int:
			out += fmt.Sprintf("%s = %d\n", sKey, field.Int())
		case reflect.String:
			out += fmt.Sprintf("%s = %s\n", sKey, field.String())
		case reflect.Bool:
			out += fmt.Sprintf("%s = %t\n", sKey, field.Bool())
		case reflect.Struct:
			switch field.Interface().(type) {
			case image.Rectangle:
				rect := field.Interface().(image.Rectangle)
				out += fmt.Sprintf("%s = %d,%d,%d,%d\n", sKey, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
			case color.RGBA:
				rgba := field.Interface().(color.RGBA)
				out += fmt.Sprintf("%s = %d,%d,%d,%d\n", sKey, rgba.R, rgba.G, rgba.B, rgba.A)
			case common.Direction:
				out += fmt.Sprintf("%s = %d\n", sKey, field.Interface().(common.Direction))

			default:
				return fmt.Errorf("unknown struct type %s", field.Kind())
			}
		default:
			return fmt.Errorf("unknown type %s", field.Kind())
		}
	}

	err = os.WriteFile(path, []byte(out), 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
