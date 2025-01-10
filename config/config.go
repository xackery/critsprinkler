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
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/dialog"
)

const (
	fileName = "critsprinkler.ini"
)

var (
	mu sync.RWMutex
)

type CritSprinklerConfiguration struct {
	IsNew      bool
	LogPath    string          `config:"log_path" config_default:""`
	EQPath     string          `config:"eq_path" config_default:""`
	MainWindow image.Rectangle `config:"main_window"`

	MeleeHitOut    common.Placement `config:"melee_hit_out" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	MeleeHitIn     common.Placement `config:"melee_hit_in" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	MeleeCritOut   common.Placement `config:"melee_crit_out" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	MeleeCritIn    common.Placement `config:"melee_crit_in" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	MeleeMissOut   common.Placement `config:"melee_miss_out" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	MeleeMissIn    common.Placement `config:"melee_miss_in" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	SpellHitOut    common.Placement `config:"spell_hit_out" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	SpellHitIn     common.Placement `config:"spell_hit_in" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	SpellCritOut   common.Placement `config:"spell_crit_out" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	SpellCritIn    common.Placement `config:"spell_crit_in" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	SpellMissOut   common.Placement `config:"spell_miss_out" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	SpellMissIn    common.Placement `config:"spell_miss_in" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	HealHitOut     common.Placement `config:"heal_hit_out" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	HealHitIn      common.Placement `config:"heal_hit_in" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	HealCritOut    common.Placement `config:"heal_crit_out" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	HealCritIn     common.Placement `config:"heal_crit_in" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	RuneHitOut     common.Placement `config:"rune_hit_out" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	RuneHitIn      common.Placement `config:"rune_hit_in" config_default:"0,1,220,307,420,407,255,0,255,255,0,2"`
	TotalDamageIn  common.Placement `config:"total_damage_in" config_default:"1,1,220,307,420,407,255,0,255,255,0,2"`
	TotalDamageOut common.Placement `config:"total_damage_out" config_default:"1,1,220,307,420,407,255,0,255,255,0,2"`
	TotalHealIn    common.Placement `config:"total_heal_in" config_default:"1,1,220,307,420,407,255,0,255,255,0,2"`
	TotalHealOut   common.Placement `config:"total_heal_out" config_default:"1,1,220,307,420,407,255,0,255,255,0,2"`
	Money          common.Placement `config:"money" config_default:"1,0,220,307,420,407,255,0,255,255,0,2"`

	IsFullscreenBorderless bool          `config:"is_fullscreen_borderless" config_default:"false"`
	IsCommaEnabled         bool          `config:"is_comma_enabled" config_default:"true"`
	PopupTallyDuration     time.Duration `config:"popup_tally_duration" config_default:"5000000000"`

	PopupIsCommaEnabled bool `config:"popup_is_comma_enabled" config_default:"true"`
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
		config := &CritSprinklerConfiguration{
			IsNew:      true,
			MainWindow: mainWindowRect,

			LogPath: "",
		}

		// fil config with default values
		for i := range reflect.TypeOf(*config).NumField() {
			sKey, ok := reflect.StructTag(reflect.TypeOf(*config).Field(i).Tag).Lookup("config")
			if !ok {
				continue
			}

			err = config.resetDefault(sKey)
			if err != nil {
				return nil, fmt.Errorf("set default value for %s: %w", sKey, err)
			}
		}

		return config, nil
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

				if strings.HasSuffix(key, "_color") {
					base := strings.TrimSuffix(key, "_color")
					for i := range reflect.TypeOf(config).NumField() {
						sKey, ok := reflect.StructTag(reflect.TypeOf(config).Field(i).Tag).Lookup("config")
						if !ok {
							continue
						}

						if sKey != base {
							continue
						}

						field := reflect.ValueOf(&config).Elem().Field(i)
						rgba := color.RGBA{}
						parts := strings.Split(value, ",")
						if len(parts) != 4 {
							return nil, fmt.Errorf("line %d parse %s=%s to color.RGBA: invalid number of parts", lineNumber, key, value)
						}
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
						isKnown = true
						foundKeys[key] = true
					}
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
				case reflect.Int64:
					val, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return nil, fmt.Errorf("line %d parse %s=%s to int64: %w", lineNumber, key, value, err)
					}

					field.SetInt(val)

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
					case common.Placement:
						parts := strings.Split(value, ",")
						if len(parts) == 4 {
							dialog.MsgBox("New Version", "CritSprinkler had a big change, so settings will be reset. Be sure to save your settings!")
							for i := range reflect.TypeOf(config).NumField() {
								sKey, ok := reflect.StructTag(reflect.TypeOf(config).Field(i).Tag).Lookup("config")
								if !ok {
									continue
								}

								err = config.resetDefault(sKey)
								if err != nil {
									return nil, fmt.Errorf("set default value for %s: %w", sKey, err)
								}
							}
							return &config, nil
						}
						windowRect := &image.Rectangle{}
						rgba := color.RGBA{}
						for i := 0; i < 11; i++ {
							val, err := strconv.Atoi(parts[i])
							if err != nil {
								return nil, fmt.Errorf("line %d parse %s=%s to image.Rectangle: %w", lineNumber, key, value, err)
							}

							switch i {
							case 0: // IsVisible
								field.Field(i).SetInt(int64(val))
							case 1: // IsTallyEnabled
								field.Field(i).SetInt(int64(val))
							case 2: // WindowRect
								windowRect.Min.X = val
							case 3:
								windowRect.Min.Y = val
							case 4:
								windowRect.Max.X = val
							case 5:
								windowRect.Max.Y = val
								field.Field(2).Set(reflect.ValueOf(windowRect))
							case 6: // FontColor
								rgba.R = uint8(val)
							case 7:
								rgba.G = uint8(val)
							case 8:
								rgba.B = uint8(val)
							case 9:
								rgba.A = uint8(val)
								field.Field(3).Set(reflect.ValueOf(rgba))
							case 10: // Direction
								field.Field(4).Set(reflect.ValueOf(common.Direction(val)))
							case 11: // Font
								field.Field(5).Set(reflect.ValueOf(common.Font(val)))
							}
						}

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
				fmt.Println("line", lineNumber, "unknown key", key)
				//return nil, fmt.Errorf("line %d unknown key %s", lineNumber, key)
				err = config.resetDefault(key)
				if err != nil {
					return nil, fmt.Errorf("line %d set default value for %s: %w", lineNumber, key, err)
				}

			}
		}
	}

	if len(foundKeys) != len(knownKeys) {
		for i := range knownKeys {
			if !foundKeys[knownKeys[i]] {
				fmt.Println("missing key", knownKeys[i])
				err = config.resetDefault(knownKeys[i])
				if err != nil {
					return nil, fmt.Errorf("set default value for %s: %w", knownKeys[i], err)
				}
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
		case reflect.Int64:
			out += fmt.Sprintf("%s = %d\n", sKey, field.Int())
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
			case common.Placement:
				placement := field.Interface().(common.Placement)
				out += fmt.Sprintf("%s = %d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n", sKey, placement.IsVisible, placement.IsTallyEnabled, placement.WindowRect.Min.X, placement.WindowRect.Min.Y, placement.WindowRect.Max.X, placement.WindowRect.Max.Y, placement.FontColor.R, placement.FontColor.G, placement.FontColor.B, placement.FontColor.A, placement.Direction, placement.Font)
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

// resetDefault sets a default value for a key based on config_default
func (c *CritSprinklerConfiguration) resetDefault(key string) error {
	for i := range reflect.TypeOf(*c).NumField() {
		sKey, ok := reflect.StructTag(reflect.TypeOf(*c).Field(i).Tag).Lookup("config")
		if !ok {
			continue
		}

		if sKey != key {
			continue
		}

		if reflect.TypeOf(*c).Field(i).Tag.Get("config_default") == "" {
			continue
		}

		field := reflect.ValueOf(c).Elem().Field(i)
		switch field.Kind() {
		case reflect.Int:
			val, err := strconv.Atoi(reflect.TypeOf(*c).Field(i).Tag.Get("config_default"))
			if err != nil {
				return fmt.Errorf("parse %s to int: %w", reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), err)
			}
			field.SetInt(int64(val))
		case reflect.String:
			field.SetString(reflect.TypeOf(*c).Field(i).Tag.Get("config_default"))
		case reflect.Bool:
			val, err := strconv.ParseBool(reflect.TypeOf(*c).Field(i).Tag.Get("config_default"))
			if err != nil {
				return fmt.Errorf("parse %s to bool: %w", reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), err)
			}
			field.SetBool(val)
		case reflect.Int64:
			val, err := strconv.ParseInt(reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), 10, 64)
			if err != nil {
				return fmt.Errorf("parse %s to int64: %w", reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), err)
			}
			field.SetInt(val)

		case reflect.Struct:
			switch field.Interface().(type) {
			case image.Rectangle:
				parts := strings.Split(reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), ",")
				if len(parts) != 4 {
					return fmt.Errorf("parse %s to image.Rectangle: invalid number of parts", reflect.TypeOf(*c).Field(i).Tag.Get("config_default"))
				}

				var rect image.Rectangle
				for i := range parts {
					val, err := strconv.Atoi(parts[i])
					if err != nil {
						return fmt.Errorf("parse %s to image.Rectangle: %w", reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), err)
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
				parts := strings.Split(reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), ",")
				if len(parts) != 4 {
					return fmt.Errorf("parse %s to color.RGBA: invalid number of parts", reflect.TypeOf(*c).Field(i).Tag.Get("config_default"))
				}

				var rgba color.RGBA
				for i := range parts {
					val, err := strconv.Atoi(parts[i])
					if err != nil {
						return fmt.Errorf("parse %s to color.RGBA: %w", reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), err)
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
				val, err := strconv.Atoi(reflect.TypeOf(*c).Field(i).Tag.Get("config_default"))
				if err != nil {
					return fmt.Errorf("parse %s to common.Direction: %w", reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), err)
				}
				field.Set(reflect.ValueOf(common.Direction(val)))
			case common.Placement:
				parts := strings.Split(reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), ",")
				windowRect := &image.Rectangle{}
				rgba := color.RGBA{}
				for i := 0; i < 11; i++ {
					val, err := strconv.Atoi(parts[i])
					if err != nil {
						return fmt.Errorf("parse %s to image.Rectangle: %w", reflect.TypeOf(*c).Field(i).Tag.Get("config_default"), err)
					}

					switch i {
					case 0: // IsVisible
						field.Field(i).SetInt(int64(val))
					case 1: // IsTallyEnabled
						field.Field(i).SetInt(int64(val))
					case 2: // WindowRect
						windowRect.Min.X = val
					case 3:
						windowRect.Min.Y = val
					case 4:
						windowRect.Max.X = val
					case 5:
						windowRect.Max.Y = val
						field.Field(2).Set(reflect.ValueOf(windowRect))
					case 6: // FontColor
						rgba.R = uint8(val)
					case 7:
						rgba.G = uint8(val)
					case 8:
						rgba.B = uint8(val)
					case 9:
						rgba.A = uint8(val)
						field.Field(3).Set(reflect.ValueOf(rgba))
					case 10: // Direction
						field.Field(4).Set(reflect.ValueOf(common.Direction(val)))
					case 11: // Font
						field.Field(5).Set(reflect.ValueOf(common.Font(val)))
					}
				}

			default:
				return fmt.Errorf("unknown struct type %s", field.Kind())
			}
		default:
			return fmt.Errorf("unknown type %s", field.Kind())
		}
	}

	return nil
}
