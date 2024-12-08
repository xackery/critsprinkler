package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CritSprinklerConfiguration struct {
	IsNew     bool
	LogPath   string
	SettingsX int
	SettingsY int
	SettingsW int
	SettingsH int
}

// LoadCritSprinklerConfig loads an CritSprinkler config file
func LoadCritSprinklerConfig(path string) (*CritSprinklerConfiguration, error) {
	_, err := os.Stat(path)
	if err != nil {
		return &CritSprinklerConfiguration{IsNew: true,
			SettingsW: 365,
			SettingsH: 371,
		}, nil
	}

	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %s", strings.TrimPrefix(err.Error(), "open critsprinkler.ini: "))
	}
	defer r.Close()

	var config CritSprinklerConfiguration

	reader := bufio.NewScanner(r)
	for reader.Scan() {
		line := reader.Text()
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
			switch key {
			case "log_path":
				config.LogPath = value
			case "settings_x":
				config.SettingsX, err = strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("parse settings_x: %w", err)
				}
			case "settings_y":
				config.SettingsY, err = strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("parse settings_y: %w", err)
				}
			case "settings_w":
				config.SettingsW, err = strconv.Atoi(value)
				if err != nil {

					return nil, fmt.Errorf("parse settings_w: %w", err)
				}
			case "settings_h":
				config.SettingsH, err = strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("parse settings_h: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown key in critsprinkler.ini: %s", key)
			}
		}
	}

	return &config, nil
}

// Save saves the config
func (c *CritSprinklerConfiguration) Save() error {
	fi, err := os.Stat("critsprinkler.ini")
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat critsprinkler.ini: %w", err)
		}
		w, err := os.Create("critsprinkler.ini")
		if err != nil {
			return fmt.Errorf("create critsprinkler.ini: %w", err)
		}
		w.Close()
	}
	if fi != nil && fi.IsDir() {
		return fmt.Errorf("critsprinkler.ini is a directory")
	}

	r, err := os.Open("critsprinkler.ini")
	if err != nil {
		return fmt.Errorf("open: %s", strings.TrimPrefix(err.Error(), "open critsprinkler.ini: "))
	}
	defer r.Close()

	tmpConfig := CritSprinklerConfiguration{}

	out := ""
	reader := bufio.NewScanner(r)
	for reader.Scan() {
		line := reader.Text()
		if strings.HasPrefix(line, "#") {
			out += line + "\n"
			continue
		}
		if !strings.Contains(line, "=") {
			out += line + "\n"
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "log_path":
			if tmpConfig.LogPath == "1" {
				continue
			}
			out += fmt.Sprintf("%s = %s\n", key, c.LogPath)
			tmpConfig.LogPath = "1"
			continue
		case "settings_x":
			if tmpConfig.SettingsX == 1 {
				continue
			}
			out += fmt.Sprintf("%s = %d\n", key, c.SettingsX)
			tmpConfig.SettingsX = 1
			continue
		case "settings_y":
			if tmpConfig.SettingsY == 1 {
				continue
			}

			out += fmt.Sprintf("%s = %d\n", key, c.SettingsY)
			tmpConfig.SettingsY = 1
			continue
		case "settings_w":
			if tmpConfig.SettingsW == 1 {
				continue
			}

			out += fmt.Sprintf("%s = %d\n", key, c.SettingsW)
			tmpConfig.SettingsW = 1
			continue
		case "settings_h":
			if tmpConfig.SettingsH == 1 {
				continue
			}

			out += fmt.Sprintf("%s = %d\n", key, c.SettingsH)
			tmpConfig.SettingsH = 1
			continue
		}

		line = fmt.Sprintf("%s = %s", key, value)
		out += line + "\n"
	}

	if tmpConfig.LogPath != "1" {
		out += fmt.Sprintf("log_path = %s\n", c.LogPath)
	}

	if tmpConfig.SettingsX != 1 {
		out += fmt.Sprintf("settings_x = %d\n", c.SettingsX)
	}

	if tmpConfig.SettingsY != 1 {
		out += fmt.Sprintf("settings_y = %d\n", c.SettingsY)
	}

	if tmpConfig.SettingsW != 1 {
		out += fmt.Sprintf("settings_w = %d\n", c.SettingsW)
	}

	if tmpConfig.SettingsH != 1 {

		out += fmt.Sprintf("settings_h = %d\n", c.SettingsH)
	}

	err = os.WriteFile("critsprinkler.ini", []byte(out), 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
