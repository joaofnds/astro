package api

import (
	"fmt"
	"io"
	"net/http"
)

// CreateToken creates a new API token. This is an unauthenticated endpoint,
// so it uses its own http.Client rather than requiring a Client instance.
func CreateToken(baseURL string) (string, error) {
	client := &http.Client{Timeout: defaultTimeout}

	req, err := http.NewRequest("POST", baseURL+"/token", nil)
	if err != nil {
		return "", fmt.Errorf("api: creating token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/text")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("api: reading token response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &Error{
			StatusCode: resp.StatusCode,
			Method:     "POST",
			Path:       "/token",
			Body:       string(body),
		}
	}

	return string(body), nil
}

// TestToken validates a token against the API. Returns nil if valid.
func TestToken(baseURL, token string) error {
	client := &http.Client{Timeout: defaultTimeout}

	req, err := http.NewRequest("GET", baseURL+"/tokentest", nil)
	if err != nil {
		return fmt.Errorf("api: creating token test request: %w", err)
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("api: reading token test response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &Error{
			StatusCode: resp.StatusCode,
			Method:     "GET",
			Path:       "/tokentest",
			Body:       string(body),
		}
	}

	return nil
}
