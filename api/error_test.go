package api

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	e := &Error{StatusCode: 404, Method: "GET", Path: "/habits/abc"}
	want := "api: GET /habits/abc returned 404"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestErrorMessageWithMessage(t *testing.T) {
	e := &Error{StatusCode: 422, Method: "POST", Path: "/habits", Message: "name required"}
	want := "api: POST /habits returned 422: name required"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestErrorAs(t *testing.T) {
	var err error = &Error{StatusCode: 500, Method: "GET", Path: "/test"}
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatal("errors.As should find *Error")
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
	}
}

func TestErrorNotAs(t *testing.T) {
	err := fmt.Errorf("network failure")
	var apiErr *Error
	if errors.As(err, &apiErr) {
		t.Fatal("errors.As should not find *Error in a plain error")
	}
}
