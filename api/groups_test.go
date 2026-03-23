package api

import (
	"astro/domain"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newGroupsTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return NewClient(srv.URL, "test-token", WithHTTPClient(srv.Client()))
}

func TestGroupsAndHabits(t *testing.T) {
	t.Run("success returns sorted groups and habits", func(t *testing.T) {
		c := newGroupsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/groups" {
				t.Errorf("Path = %s, want /groups", r.URL.Path)
			}
			if got := r.Header.Get("Authorization"); got != "test-token" {
				t.Errorf("Authorization = %q, want %q", got, "test-token")
			}
			payload := GroupsAndHabitsPayload{
				Groups: []*domain.Group{
					{ID: "g2", Name: "Zzz"},
					{ID: "g1", Name: "Aaa"},
				},
				Habits: []*domain.Habit{
					{ID: "h2", Name: "Zzz"},
					{ID: "h1", Name: "Aaa"},
				},
			}
			json.NewEncoder(w).Encode(payload)
		}))

		groups, habits, err := c.GroupsAndHabits(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(groups) != 2 {
			t.Fatalf("len(groups) = %d, want 2", len(groups))
		}
		if groups[0].Name != "Aaa" {
			t.Errorf("groups[0].Name = %q, want %q", groups[0].Name, "Aaa")
		}
		if groups[1].Name != "Zzz" {
			t.Errorf("groups[1].Name = %q, want %q", groups[1].Name, "Zzz")
		}
		if len(habits) != 2 {
			t.Fatalf("len(habits) = %d, want 2", len(habits))
		}
		if habits[0].Name != "Aaa" {
			t.Errorf("habits[0].Name = %q, want %q", habits[0].Name, "Aaa")
		}
	})

	t.Run("server error returns api.Error", func(t *testing.T) {
		c := newGroupsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		_, _, err := c.GroupsAndHabits(context.Background())
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

func TestCreateGroup(t *testing.T) {
	t.Run("success posts JSON body and returns group", func(t *testing.T) {
		c := newGroupsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/groups" {
				t.Errorf("Path = %s, want /groups", r.URL.Path)
			}
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}
			group := domain.Group{ID: "g1", Name: "Morning"}
			json.NewEncoder(w).Encode(group)
		}))

		g, err := c.CreateGroup(context.Background(), "Morning")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if g.ID != "g1" {
			t.Errorf("ID = %q, want %q", g.ID, "g1")
		}
		if g.Name != "Morning" {
			t.Errorf("Name = %q, want %q", g.Name, "Morning")
		}
	})

	t.Run("server error returns api.Error", func(t *testing.T) {
		c := newGroupsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		_, err := c.CreateGroup(context.Background(), "x")
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

func TestAddToGroup(t *testing.T) {
	t.Run("success sends POST to correct path", func(t *testing.T) {
		c := newGroupsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Method = %s, want POST", r.Method)
			}
			// AddToGroup(habitID="h1", groupID="g1") => POST /groups/g1/h1
			if r.URL.Path != "/groups/g1/h1" {
				t.Errorf("Path = %s, want /groups/g1/h1", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}))

		if err := c.AddToGroup(context.Background(), "h1", "g1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found returns api.Error", func(t *testing.T) {
		c := newGroupsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		err := c.AddToGroup(context.Background(), "h1", "missing")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *api.Error, got %T: %v", err, err)
		}
		if apiErr.StatusCode != 404 {
			t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
		}
	})
}

func TestRemoveFromGroup(t *testing.T) {
	t.Run("success sends DELETE to correct path", func(t *testing.T) {
		c := newGroupsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" {
				t.Errorf("Method = %s, want DELETE", r.Method)
			}
			// RemoveFromGroup(habitID="h1", groupID="g1") => DELETE /groups/g1/h1
			if r.URL.Path != "/groups/g1/h1" {
				t.Errorf("Path = %s, want /groups/g1/h1", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}))

		if err := c.RemoveFromGroup(context.Background(), "h1", "g1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found returns api.Error", func(t *testing.T) {
		c := newGroupsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		err := c.RemoveFromGroup(context.Background(), "h1", "missing")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *api.Error, got %T: %v", err, err)
		}
		if apiErr.StatusCode != 404 {
			t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
		}
	})
}

func TestDeleteGroup(t *testing.T) {
	t.Run("success sends DELETE", func(t *testing.T) {
		c := newGroupsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" {
				t.Errorf("Method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/groups/g1" {
				t.Errorf("Path = %s, want /groups/g1", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}))

		if err := c.DeleteGroup(context.Background(), "g1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found returns api.Error", func(t *testing.T) {
		c := newGroupsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		err := c.DeleteGroup(context.Background(), "missing")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *api.Error, got %T: %v", err, err)
		}
		if apiErr.StatusCode != 404 {
			t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
		}
	})
}
