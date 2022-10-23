package habit

import (
	"errors"
	"io"
	"log"
	"os"
	"path"
)

var (
	configDir string
	tokenPath string
)

func Init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get user home dir: %v\n", err)
	}
	configDir = path.Join(home, ".config", "astro")
	tokenPath = path.Join(configDir, "token")

	ensureTokenExists()

	f, err := os.Open(tokenPath)
	if err != nil {
		log.Fatalf("could not open token file: %v\n", err)
	}

	token, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("could not read token file(%q): %s", tokenPath, err)
	}

	InitClient(string(token))
}

func ensureTokenExists() {
	_, err := os.Stat(tokenPath)
	if err == nil {
		return
	}

	if !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("could not stat token file: %v\n", err)
	}

	err = os.Mkdir(configDir, 0755)
	if err != nil {
		log.Fatalf("could not create dir %q: %v\n", configDir, err)
	}

	res, err := NewAPI().CreateToken()
	if err != nil {
		log.Fatalf("could not create token: %v\n", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("could not read response body: %v\n", err)
	}

	err = os.WriteFile(tokenPath, b, 0644)
	if err != nil {
		log.Fatalf("could not write token file: %v\n", err)
	}
}
