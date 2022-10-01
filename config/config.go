package config

import (
	"errors"
	"io"
	"log"
	"os"
	"path"
)

const (
	TimeFormat      = "Jan 02, 2006"
	TimeFrameInDays = 52 * 7

	Graphic         = "⬛"
	SelectedGraphic = "⚫"
)

var Token string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get user home dir")
	}

	tokenPath := path.Join(home, ".config", "astro", "token")
	f, err := os.Open(tokenPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("could not find token at %q", tokenPath)
		}
		panic(err)
	}

	token, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("could not read token file(%q) %s", tokenPath, err)
	}
	Token = string(token)
}
