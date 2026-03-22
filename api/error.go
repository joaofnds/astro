package api

import "fmt"

// Error represents a non-successful HTTP response from the API.
type Error struct {
	StatusCode int
	Method     string
	Path       string
	Message    string
	Body       string
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("api: %s %s returned %d: %s", e.Method, e.Path, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("api: %s %s returned %d", e.Method, e.Path, e.StatusCode)
}
