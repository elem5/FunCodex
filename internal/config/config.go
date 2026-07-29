package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	AccentColor        string   `json:"accent_color"` // hex (e.g., "#FF5733")
	ShowAscii          bool     `json:"show_ascii"`
	SeparatorCharacter string   `json:"separator_character"`
	Keybinds           Keybinds `json:"keybinds"`
}

type Keybinds struct {
	EncryptMode string `json:"encrypt_mode"`
	DecryptMode string `json:"decrypt_mode"`
	Length      string `json:"length_mode"`
}

// default returns reasonable fallback settings if no config file exists yet.
func Default() Config {
	return Config{
		AccentColor:        "#55559b",
		ShowAscii:          true,
		SeparatorCharacter: "|",
		Keybinds: Keybinds{
			EncryptMode: "ctrl+e",
			DecryptMode: "ctrl+d",
			Length:      "ctrl+l",
		},
	}
}

func Path() (string, error) {
	var dir string

	if runtime.GOOS == "windows" {
		dir = os.Getenv("LOCALAPPDATA")
		if dir == "" {
			var err error
			dir, err = os.UserConfigDir()
			if err != nil {
				return "", err
			}
		}
	} else {
		var err error
		dir, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(
		dir,
		"FunCodex",
		"config.json",
	), nil
}

func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}

	// create directory if it doesn't exist
	err = os.MkdirAll(
		filepath.Dir(path),
		0755,
	)
	if err != nil {
		return nil, err
	}

	// create default config file if missing
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := Default()

		data, err := json.MarshalIndent(
			cfg,
			"",
			"    ",
		)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(
			path,
			data,
			0644,
		)
		if err != nil {
			return nil, err
		}

		return &cfg, nil
	}

	// load existing config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(
		data,
		&cfg,
	)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
