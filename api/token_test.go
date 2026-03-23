package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTokenFunc(t *testing.T) {
	t.Run("success returns token string", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/token" {
				t.Errorf("Path = %s, want /token", r.URL.Path)
			}
			if ct := r.Header.Get("Content-Type"); ct != "application/text" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/text")
			}
			w.Write([]byte("new-token-abc"))
		}))
		t.Cleanup(srv.Close)

		token, err := CreateToken(srv.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "new-token-abc" {
			t.Errorf("token = %q, want %q", token, "new-token-abc")
		}
	})

	t.Run("server error returns api.Error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		t.Cleanup(srv.Close)

		_, err := CreateToken(srv.URL)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *api.Error, got %T: %v", err, err)
		}
		if apiErr.StatusCode != 500 {
			t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
		}
	})
}

func TestTestTokenFunc(t *testing.T) {
	t.Run("success returns nil for valid token", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/tokentest" {
				t.Errorf("Path = %s, want /tokentest", r.URL.Path)
			}
			if got := r.Header.Get("Authorization"); got != "test-token" {
				t.Errorf("Authorization = %q, want %q", got, "test-token")
			}
			w.WriteHeader(http.StatusOK)
		}))
		t.Cleanup(srv.Close)

		if err := TestToken(srv.URL, "test-token"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unauthorized returns api.Error with 401", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		t.Cleanup(srv.Close)

		err := TestToken(srv.URL, "bad-token")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *api.Error, got %T: %v", err, err)
		}
		if apiErr.StatusCode != 401 {
			t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
		}
	})
}
