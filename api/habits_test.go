package api

import (
	"astro/domain"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newHabitsTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return NewClient(srv.URL, "test-token", WithHTTPClient(srv.Client()))
}

func TestListHabits(t *testing.T) {
	t.Run("success returns sorted habits with auth header", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/habits" {
				t.Errorf("Path = %s, want /habits", r.URL.Path)
			}
			if got := r.Header.Get("Authorization"); got != "test-token" {
				t.Errorf("Authorization = %q, want %q", got, "test-token")
			}
			habits := []*domain.Habit{
				{ID: "2", Name: "Zzz"},
				{ID: "1", Name: "Aaa"},
			}
			json.NewEncoder(w).Encode(habits)
		}))
		t.Cleanup(srv.Close)
		c := NewClient(srv.URL, "test-token", WithHTTPClient(srv.Client()))

		habits, err := c.ListHabits(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(habits) != 2 {
			t.Fatalf("len = %d, want 2", len(habits))
		}
		if habits[0].Name != "Aaa" {
			t.Errorf("habits[0].Name = %q, want %q", habits[0].Name, "Aaa")
		}
		if habits[1].Name != "Zzz" {
			t.Errorf("habits[1].Name = %q, want %q", habits[1].Name, "Zzz")
		}
	})

	t.Run("server error returns api.Error", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		_, err := c.ListHabits(context.Background())
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

	t.Run("invalid JSON returns decode error", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("not json"))
		}))

		_, err := c.ListHabits(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var apiErr *Error
		if errors.As(err, &apiErr) {
			t.Fatalf("expected non-api.Error decode error, got *api.Error: %v", apiErr)
		}
	})
}

func TestCreateHabit(t *testing.T) {
	t.Run("success posts name and returns habit with sorted activities", func(t *testing.T) {
		now := time.Now().UTC().Truncate(time.Second)
		earlier := now.Add(-time.Hour)

		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/habits" {
				t.Errorf("Path = %s, want /habits", r.URL.Path)
			}
			if got := r.URL.Query().Get("name"); got != "Read books" {
				t.Errorf("name = %q, want %q", got, "Read books")
			}
			h := domain.Habit{
				ID:   "h1",
				Name: "Read books",
				Activities: []domain.Activity{
					{ID: "a2", Desc: "later", CreatedAt: now},
					{ID: "a1", Desc: "earlier", CreatedAt: earlier},
				},
			}
			json.NewEncoder(w).Encode(h)
		}))

		h, err := c.CreateHabit(context.Background(), "Read books")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if h.Name != "Read books" {
			t.Errorf("Name = %q, want %q", h.Name, "Read books")
		}
		if len(h.Activities) != 2 {
			t.Fatalf("len(Activities) = %d, want 2", len(h.Activities))
		}
		if h.Activities[0].Desc != "earlier" {
			t.Errorf("Activities[0].Desc = %q, want %q (should be sorted)", h.Activities[0].Desc, "earlier")
		}
	})

	t.Run("server error returns api.Error", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		_, err := c.CreateHabit(context.Background(), "x")
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

func TestGetHabit(t *testing.T) {
	t.Run("success returns habit with sorted activities", func(t *testing.T) {
		now := time.Now().UTC().Truncate(time.Second)
		earlier := now.Add(-time.Hour)

		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/habits/h1" {
				t.Errorf("Path = %s, want /habits/h1", r.URL.Path)
			}
			h := domain.Habit{
				ID:   "h1",
				Name: "Meditate",
				Activities: []domain.Activity{
					{ID: "a2", Desc: "later", CreatedAt: now},
					{ID: "a1", Desc: "earlier", CreatedAt: earlier},
				},
			}
			json.NewEncoder(w).Encode(h)
		}))

		h, err := c.GetHabit(context.Background(), "h1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if h.ID != "h1" {
			t.Errorf("ID = %q, want %q", h.ID, "h1")
		}
		if h.Activities[0].Desc != "earlier" {
			t.Errorf("Activities[0].Desc = %q, want %q (should be sorted)", h.Activities[0].Desc, "earlier")
		}
	})

	t.Run("not found returns api.Error with 404", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		_, err := c.GetHabit(context.Background(), "missing")
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

func TestUpdateHabit(t *testing.T) {
	t.Run("success sends PATCH with JSON body", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PATCH" {
				t.Errorf("Method = %s, want PATCH", r.Method)
			}
			if r.URL.Path != "/habits/h1" {
				t.Errorf("Path = %s, want /habits/h1", r.URL.Path)
			}
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}
			w.WriteHeader(http.StatusOK)
		}))

		if err := c.UpdateHabit(context.Background(), "h1", "new name"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found returns api.Error", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		err := c.UpdateHabit(context.Background(), "missing", "x")
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

func TestDeleteHabit(t *testing.T) {
	t.Run("success sends DELETE", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" {
				t.Errorf("Method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/habits/h1" {
				t.Errorf("Path = %s, want /habits/h1", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}))

		if err := c.DeleteHabit(context.Background(), "h1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found returns api.Error", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		err := c.DeleteHabit(context.Background(), "missing")
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

func TestAddActivity(t *testing.T) {
	t.Run("success sends POST with JSON body", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/habits/h1" {
				t.Errorf("Path = %s, want /habits/h1", r.URL.Path)
			}
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}
			w.WriteHeader(http.StatusOK)
		}))

		err := c.AddActivity(context.Background(), "h1", "did it", time.Now())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("server error returns api.Error", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		err := c.AddActivity(context.Background(), "h1", "x", time.Now())
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

func TestUpdateActivity(t *testing.T) {
	t.Run("success sends PATCH with JSON body", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "PATCH" {
				t.Errorf("Method = %s, want PATCH", r.Method)
			}
			if r.URL.Path != "/habits/h1/a1" {
				t.Errorf("Path = %s, want /habits/h1/a1", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}))

		if err := c.UpdateActivity(context.Background(), "h1", "a1", "updated"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found returns api.Error", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		err := c.UpdateActivity(context.Background(), "h1", "missing", "x")
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

func TestDeleteActivity(t *testing.T) {
	t.Run("success sends DELETE", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" {
				t.Errorf("Method = %s, want DELETE", r.Method)
			}
			if r.URL.Path != "/habits/h1/a1" {
				t.Errorf("Path = %s, want /habits/h1/a1", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		}))

		if err := c.DeleteActivity(context.Background(), "h1", "a1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found returns api.Error", func(t *testing.T) {
		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		err := c.DeleteActivity(context.Background(), "h1", "missing")
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

func TestCheckIn(t *testing.T) {
	t.Run("success calls AddActivity then GetHabit", func(t *testing.T) {
		now := time.Now().UTC().Truncate(time.Second)
		var addCalled, getCalled bool

		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == "POST" && r.URL.Path == "/habits/h1":
				addCalled = true
				w.WriteHeader(http.StatusOK)
			case r.Method == "GET" && r.URL.Path == "/habits/h1":
				getCalled = true
				h := domain.Habit{
					ID:   "h1",
					Name: "Exercise",
					Activities: []domain.Activity{
						{ID: "a1", Desc: "checked in", CreatedAt: now},
					},
				}
				json.NewEncoder(w).Encode(h)
			default:
				t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
				w.WriteHeader(http.StatusNotFound)
			}
		}))

		h, err := c.CheckIn(context.Background(), CheckInDTO{
			ID:   "h1",
			Desc: "checked in",
			Date: now,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !addCalled {
			t.Error("AddActivity was not called")
		}
		if !getCalled {
			t.Error("GetHabit was not called")
		}
		if h.ID != "h1" {
			t.Errorf("ID = %q, want %q", h.ID, "h1")
		}
	})

	t.Run("AddActivity failure returns error without calling GetHabit", func(t *testing.T) {
		var getCalled bool

		c := newHabitsTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == "POST":
				w.WriteHeader(http.StatusInternalServerError)
			case r.Method == "GET":
				getCalled = true
				w.WriteHeader(http.StatusOK)
			}
		}))

		_, err := c.CheckIn(context.Background(), CheckInDTO{
			ID:   "h1",
			Desc: "x",
			Date: time.Now(),
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if getCalled {
			t.Error("GetHabit should not be called when AddActivity fails")
		}
	})
}

func TestConnectionRefused(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	c := NewClient(srv.URL, "test-token", WithHTTPClient(srv.Client()))
	srv.Close()

	_, err := c.ListHabits(context.Background())
	if err == nil {
		t.Fatal("expected error for closed server, got nil")
	}
	var apiErr *Error
	if errors.As(err, &apiErr) {
		t.Fatalf("expected non-api.Error transport error, got *api.Error: %v", apiErr)
	}
}
