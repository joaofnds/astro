package token

import (
	"astro/api"
	"errors"
	"fmt"
	"io"
	"os"
)

// Init reads the token from disk, creating one via the API if it does not exist.
func Init(tokenFilePath, apiBaseURL string) (string, error) {
	if err := ensureTokenExists(tokenFilePath, apiBaseURL); err != nil {
		return "", fmt.Errorf("could not create token: %w", err)
	}

	f, err := os.Open(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("could not open token file: %w", err)
	}
	defer f.Close()

	token, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("could not read token file(%q): %w", tokenFilePath, err)
	}

	return string(token), nil
}

func ensureTokenExists(tokenFilePath, apiBaseURL string) error {
	_, err := os.Stat(tokenFilePath)
	if err == nil {
		return nil
	}

	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("could not stat token file: %w", err)
	}

	tok, err := api.CreateToken(apiBaseURL)
	if err != nil {
		return fmt.Errorf("could not create token: %w", err)
	}

	if err = os.WriteFile(tokenFilePath, []byte(tok), 0644); err != nil {
		return fmt.Errorf("could not write token file: %w", err)
	}

	return nil
}
