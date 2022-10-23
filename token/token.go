package token

import (
	"astro/habit"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

var (
	configDir string
	tokenPath string
)

func Init() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home dir: %w", err)
	}
	configDir = path.Join(home, ".config", "astro")
	tokenPath = path.Join(configDir, "token")

	if err := ensureTokenExists(); err != nil {
		return "", fmt.Errorf("could no create token: %w", err)
	}

	f, err := os.Open(tokenPath)
	if err != nil {
		return "", fmt.Errorf("could not open token file: %w", err)
	}

	token, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("could not read token file(%q): %w", tokenPath, err)
	}

	return string(token), nil
}

func ensureTokenExists() error {
	_, err := os.Stat(tokenPath)
	if err == nil {
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("could not stat token file: %w", err)
	}

	err = os.Mkdir(configDir, 0755)
	if err != nil {
		return fmt.Errorf("could not create dir %q: %w", configDir, err)
	}

	res, err := habit.NewAPI().CreateToken()
	if err != nil {
		return fmt.Errorf("could not create token: %w", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}

	err = os.WriteFile(tokenPath, b, 0644)
	if err != nil {
		return fmt.Errorf("could not write token file: %w", err)
	}

	return nil
}
