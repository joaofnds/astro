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

var (
	Width  int
	Height int

	ConfigDirPath string
	LogFilePath   string
	TokenFilePath string
)

func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home dir: %w", err)
	}

	ConfigDirPath := path.Join(home, ".config", "astro")
	if err := os.MkdirAll(ConfigDirPath, 0755); err != nil {
		return fmt.Errorf("could not create dir %q: %w", ConfigDirPath, err)
	}

	LogFilePath = path.Join(ConfigDirPath, "log.log")
	TokenFilePath = path.Join(ConfigDirPath, "token")

	return nil
}
