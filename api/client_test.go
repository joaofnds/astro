package api

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClientDefaultTimeout(t *testing.T) {
	c := NewClient("http://example.com", "tok")
	if c.httpClient.Timeout != 10*time.Second {
		t.Errorf("Timeout = %v, want 10s", c.httpClient.Timeout)
	}
}

func TestNewClientNoDefaultClient(t *testing.T) {
	c := NewClient("http://example.com", "tok")
	if c.httpClient == http.DefaultClient {
		t.Error("httpClient must not be http.DefaultClient")
	}
}

func TestWithTimeout(t *testing.T) {
	c := NewClient("http://example.com", "tok", WithTimeout(5*time.Second))
	if c.httpClient.Timeout != 5*time.Second {
		t.Errorf("Timeout = %v, want 5s", c.httpClient.Timeout)
	}
}
