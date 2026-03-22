package config

import (
	"fmt"
	"os"
	"path"
	"time"
)

const (
	DateFormat      = "Jan 02, 2006"
	TimeFormat      = time.Kitchen
	TimeFrameInDays = 52 * 7
	ShortHistSize   = 14

	Graphic         = "⬛"
	SelectedGraphic = "⚫"
)

// Config holds file paths resolved from the user's home directory.
type Config struct {
	ConfigDirPath string
	LogFilePath   string
	TokenFilePath string
}

// Load reads the user's home directory, ensures the config directory exists,
// and returns a Config with all paths set.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get user home dir: %w", err)
	}

	configDir := path.Join(home, ".config", "astro")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("could not create dir %q: %w", configDir, err)
	}

	return &Config{
		ConfigDirPath: configDir,
		LogFilePath:   path.Join(configDir, "log.log"),
		TokenFilePath: path.Join(configDir, "token"),
	}, nil
}
