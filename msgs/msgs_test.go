package msgs_test

import (
	"astro/api"
	"astro/domain"
	"astro/msgs"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

// testModel is a minimal tea.Model for testing PushScreen.
type testModel struct{}

func (testModel) Init() tea.Cmd                           { return nil }
func (testModel) Update(tea.Msg) (tea.Model, tea.Cmd)     { return testModel{}, nil }
func (testModel) View() tea.View                          { return tea.NewView("") }

func TestPushScreen(t *testing.T) {
	m := testModel{}
	cmd := msgs.PushScreen(m)
	if cmd == nil {
		t.Fatal("PushScreen returned nil cmd")
	}
	msg := cmd()
	push, ok := msg.(msgs.PushScreenMsg)
	if !ok {
		t.Fatalf("expected PushScreenMsg, got %T", msg)
	}
	if push.Screen == nil {
		t.Fatal("PushScreenMsg.Screen is nil")
	}
}

func TestPopScreen(t *testing.T) {
	cmd := msgs.PopScreen()
	if cmd == nil {
		t.Fatal("PopScreen returned nil cmd")
	}
	msg := cmd()
	pop, ok := msg.(msgs.PopScreenMsg)
	if !ok {
		t.Fatalf("expected PopScreenMsg, got %T", msg)
	}
	if pop.Cmd != nil {
		t.Fatal("PopScreenMsg.Cmd should be nil for PopScreen()")
	}
}

func TestPopScreenWith(t *testing.T) {
	inner := func() tea.Msg { return nil }
	cmd := msgs.PopScreenWith(inner)
	if cmd == nil {
		t.Fatal("PopScreenWith returned nil cmd")
	}
	msg := cmd()
	pop, ok := msg.(msgs.PopScreenMsg)
	if !ok {
		t.Fatalf("expected PopScreenMsg, got %T", msg)
	}
	if pop.Cmd == nil {
		t.Fatal("PopScreenMsg.Cmd should not be nil for PopScreenWith()")
	}
}

// newTestClient creates an api.Client pointing at the given test server.
func newTestClient(srv *httptest.Server) *api.Client {
	return api.NewClient(srv.URL, "test-token", api.WithHTTPClient(srv.Client()))
}

func TestLoadAll_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/groups" {
			http.NotFound(w, r)
			return
		}
		resp := api.GroupsAndHabitsPayload{
			Groups: []*domain.Group{{ID: "g1", Name: "Group 1"}},
			Habits: []*domain.Habit{{ID: "h1", Name: "Habit 1"}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.LoadAll(client)
	msg := cmd()

	loaded, ok := msg.(msgs.DataLoadedMsg)
	if !ok {
		t.Fatalf("expected DataLoadedMsg, got %T", msg)
	}
	if len(loaded.Habits) != 1 {
		t.Fatalf("expected 1 habit, got %d", len(loaded.Habits))
	}
	if len(loaded.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(loaded.Groups))
	}
}

func TestLoadAll_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.LoadAll(client)
	msg := cmd()

	fatal, ok := msg.(msgs.FatalErrorMsg)
	if !ok {
		t.Fatalf("expected FatalErrorMsg, got %T", msg)
	}
	if fatal.Err == nil {
		t.Fatal("FatalErrorMsg.Err should not be nil")
	}
}

func TestCreateHabit_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/habits" {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(domain.Habit{ID: "h1", Name: "New Habit"})
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.CreateHabit(client, "New Habit")
	msg := cmd()

	created, ok := msg.(msgs.HabitCreatedMsg)
	if !ok {
		t.Fatalf("expected HabitCreatedMsg, got %T", msg)
	}
	if created.Habit.Name != "New Habit" {
		t.Fatalf("expected habit name 'New Habit', got %q", created.Habit.Name)
	}
}

func TestCreateHabit_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.CreateHabit(client, "fail")
	msg := cmd()

	apiErr, ok := msg.(msgs.APIErrorMsg)
	if !ok {
		t.Fatalf("expected APIErrorMsg, got %T", msg)
	}
	if apiErr.Op != "create habit" {
		t.Fatalf("expected Op 'create habit', got %q", apiErr.Op)
	}
}

func TestDeleteHabit_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/habits/h1" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.DeleteHabit(client, "h1")
	msg := cmd()

	deleted, ok := msg.(msgs.HabitDeletedMsg)
	if !ok {
		t.Fatalf("expected HabitDeletedMsg, got %T", msg)
	}
	if deleted.ID != "h1" {
		t.Fatalf("expected ID 'h1', got %q", deleted.ID)
	}
}

func TestUpdateHabit_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPatch && r.URL.Path == "/habits/h1":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/habits/h1":
			json.NewEncoder(w).Encode(domain.Habit{ID: "h1", Name: "Updated"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.UpdateHabit(client, "h1", "Updated")
	msg := cmd()

	updated, ok := msg.(msgs.HabitUpdatedMsg)
	if !ok {
		t.Fatalf("expected HabitUpdatedMsg, got %T", msg)
	}
	if updated.Habit.Name != "Updated" {
		t.Fatalf("expected habit name 'Updated', got %q", updated.Habit.Name)
	}
}

func TestCheckIn_Success(t *testing.T) {
	now := time.Now()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/habits/h1":
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodGet && r.URL.Path == "/habits/h1":
			json.NewEncoder(w).Encode(domain.Habit{ID: "h1", Name: "My Habit"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.CheckIn(client, "h1", "done", now)
	msg := cmd()

	result, ok := msg.(msgs.CheckInResultMsg)
	if !ok {
		t.Fatalf("expected CheckInResultMsg, got %T", msg)
	}
	if result.Habit.ID != "h1" {
		t.Fatalf("expected habit ID 'h1', got %q", result.Habit.ID)
	}
}

func TestCreateGroup_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/groups" {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(domain.Group{ID: "g1", Name: "New Group"})
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.CreateGroup(client, "New Group")
	msg := cmd()

	created, ok := msg.(msgs.GroupCreatedMsg)
	if !ok {
		t.Fatalf("expected GroupCreatedMsg, got %T", msg)
	}
	if created.Group.Name != "New Group" {
		t.Fatalf("expected group name 'New Group', got %q", created.Group.Name)
	}
}

func TestDeleteGroup_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/groups/g1" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.DeleteGroup(client, "g1")
	msg := cmd()

	deleted, ok := msg.(msgs.GroupDeletedMsg)
	if !ok {
		t.Fatalf("expected GroupDeletedMsg, got %T", msg)
	}
	if deleted.ID != "g1" {
		t.Fatalf("expected ID 'g1', got %q", deleted.ID)
	}
}

func TestAddToGroup_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/groups/g1/h1" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.AddToGroup(client, "h1", "g1")
	msg := cmd()

	added, ok := msg.(msgs.AddedToGroupMsg)
	if !ok {
		t.Fatalf("expected AddedToGroupMsg, got %T", msg)
	}
	if added.HabitID != "h1" || added.GroupID != "g1" {
		t.Fatalf("expected h1/g1, got %s/%s", added.HabitID, added.GroupID)
	}
}

func TestRemoveFromGroup_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/groups/g1/h1" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.RemoveFromGroup(client, "h1", "g1")
	msg := cmd()

	removed, ok := msg.(msgs.RemovedFromGroupMsg)
	if !ok {
		t.Fatalf("expected RemovedFromGroupMsg, got %T", msg)
	}
	if removed.HabitID != "h1" || removed.GroupID != "g1" {
		t.Fatalf("expected h1/g1, got %s/%s", removed.HabitID, removed.GroupID)
	}
}

func TestUpdateActivity_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/habits/h1/a1" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.UpdateActivity(client, "h1", "a1", "updated desc")
	msg := cmd()

	updated, ok := msg.(msgs.ActivityUpdatedMsg)
	if !ok {
		t.Fatalf("expected ActivityUpdatedMsg, got %T", msg)
	}
	if updated.HabitID != "h1" || updated.ActivityID != "a1" || updated.Desc != "updated desc" {
		t.Fatalf("unexpected fields: %+v", updated)
	}
}

func TestDeleteActivity_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/habits/h1/a1" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := newTestClient(srv)
	cmd := msgs.DeleteActivity(client, "h1", "a1")
	msg := cmd()

	deleted, ok := msg.(msgs.ActivityDeletedMsg)
	if !ok {
		t.Fatalf("expected ActivityDeletedMsg, got %T", msg)
	}
	if deleted.HabitID != "h1" || deleted.ActivityID != "a1" {
		t.Fatalf("unexpected fields: %+v", deleted)
	}
}

func TestClearStatusAfter(t *testing.T) {
	cmd := msgs.ClearStatusAfter(time.Millisecond)
	if cmd == nil {
		t.Fatal("ClearStatusAfter returned nil cmd")
	}
	// ClearStatusAfter wraps tea.Tick, which returns a tea.Cmd.
	// We verify the cmd is non-nil; the actual ClearStatusMsg delivery
	// depends on the Bubbletea runtime processing the tick.
}

func TestAPIErrorMsg_Error(t *testing.T) {
	msg := msgs.APIErrorMsg{Err: errors.New("connection refused"), Op: "create habit"}
	if msg.Err == nil {
		t.Fatal("expected non-nil error")
	}
	if msg.Op != "create habit" {
		t.Fatalf("expected Op 'create habit', got %q", msg.Op)
	}
}

func TestFatalErrorMsg_Error(t *testing.T) {
	msg := msgs.FatalErrorMsg{Err: errors.New("token expired")}
	if msg.Err == nil {
		t.Fatal("expected non-nil error")
	}
}
