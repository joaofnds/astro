package token

import (
	"astro/config"
	"astro/habit"
	"errors"
	"fmt"
	"io"
	"os"
)

func Init() (string, error) {
	if err := ensureTokenExists(); err != nil {
		return "", fmt.Errorf("could no create token: %w", err)
	}

	f, err := os.Open(config.TokenFilePath)
	if err != nil {
		return "", fmt.Errorf("could not open token file: %w", err)
	}

	token, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("could not read token file(%q): %w", config.TokenFilePath, err)
	}

	return string(token), nil
}

func ensureTokenExists() error {
	_, err := os.Stat(config.TokenFilePath)
	if err == nil {
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("could not stat token file: %w", err)
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

	if err = os.WriteFile(config.TokenFilePath, b, 0644); err != nil {
		return fmt.Errorf("could not write token file: %w", err)
	}

	return nil
}
