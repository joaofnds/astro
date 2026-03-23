package list_test

import (
	"astro/api"
	"astro/domain"
	"astro/models/list"
	"astro/models/textinput"
	"astro/msgs"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestList(t *testing.T, habits []*domain.Habit, groups []*domain.Group) list.List {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	client := api.NewClient(srv.URL, "test-token", api.WithHTTPClient(srv.Client()))
	return list.NewList(client, habits, groups, 80, 24)
}

func keyPress(code rune, text string) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code, Text: text})
}

func enterKey() tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter})
}

func TestListView_ContainsHabitNames(t *testing.T) {
	habits := []*domain.Habit{
		{ID: "h1", Name: "Read"},
		{ID: "h2", Name: "Meditate"},
	}
	m := newTestList(t, habits, nil)
	v := m.View()

	for _, h := range habits {
		if !strings.Contains(v.Content, h.Name) {
			t.Errorf("View should contain habit name %q, got:\n%s", h.Name, v.Content)
		}
	}
}

func TestListView_ContainsTitle(t *testing.T) {
	m := newTestList(t, []*domain.Habit{{ID: "h1", Name: "Run"}}, nil)
	v := m.View()

	if !strings.Contains(v.Content, "Habits") {
		t.Errorf("View should contain title 'Habits', got:\n%s", v.Content)
	}
}

func TestListView_ContainsHelpKeys(t *testing.T) {
	habits := []*domain.Habit{{ID: "h1", Name: "Run"}}
	m := newTestList(t, habits, nil)
	v := m.View()

	for _, label := range []string{"add", "view", "delete"} {
		if !strings.Contains(v.Content, label) {
			t.Errorf("View should contain help label %q, got:\n%s", label, v.Content)
		}
	}
}

func TestHabitCreatedMsg(t *testing.T) {
	m := newTestList(t, nil, nil)

	created := msgs.HabitCreatedMsg{Habit: &domain.Habit{ID: "h1", Name: "New Habit"}}
	updated, cmd := m.Update(created)
	m = updated.(list.List)

	if cmd == nil {
		t.Error("expected non-nil cmd (status message)")
	}

	v := m.View()
	if !strings.Contains(v.Content, "New Habit") {
		t.Errorf("View should contain newly created habit, got:\n%s", v.Content)
	}
}

func TestHabitDeletedMsg(t *testing.T) {
	habits := []*domain.Habit{
		{ID: "h1", Name: "Keep"},
		{ID: "h2", Name: "Remove"},
	}
	m := newTestList(t, habits, nil)

	updated, _ := m.Update(msgs.HabitDeletedMsg{ID: "h2"})
	m = updated.(list.List)

	v := m.View()
	if strings.Contains(v.Content, "Remove") {
		t.Error("View should not contain deleted habit 'Remove'")
	}
	if !strings.Contains(v.Content, "Keep") {
		t.Error("View should still contain 'Keep'")
	}
}

func TestHabitUpdatedMsg(t *testing.T) {
	habits := []*domain.Habit{{ID: "h1", Name: "Old Name"}}
	m := newTestList(t, habits, nil)

	updated, _ := m.Update(msgs.HabitUpdatedMsg{Habit: &domain.Habit{ID: "h1", Name: "New Name"}})
	m = updated.(list.List)

	v := m.View()
	if !strings.Contains(v.Content, "New Name") {
		t.Errorf("View should contain updated name 'New Name', got:\n%s", v.Content)
	}
}

func TestCheckInResultMsg(t *testing.T) {
	habits := []*domain.Habit{{ID: "h1", Name: "Run"}}
	m := newTestList(t, habits, nil)

	checkedIn := &domain.Habit{ID: "h1", Name: "Run", Activities: []domain.Activity{{ID: "a1"}}}
	updated, cmd := m.Update(msgs.CheckInResultMsg{Habit: checkedIn})
	_ = updated.(list.List)

	if cmd == nil {
		t.Error("expected non-nil cmd (status message)")
	}
}

func TestGroupCreatedMsg(t *testing.T) {
	m := newTestList(t, nil, nil)

	created := msgs.GroupCreatedMsg{Group: &domain.Group{ID: "g1", Name: "Fitness"}}
	updated, cmd := m.Update(created)
	m = updated.(list.List)

	if cmd == nil {
		t.Error("expected non-nil cmd (status message)")
	}

	v := m.View()
	if !strings.Contains(v.Content, "Fitness") {
		t.Errorf("View should contain new group name 'Fitness', got:\n%s", v.Content)
	}
}

func TestGroupDeletedMsg(t *testing.T) {
	groups := []*domain.Group{
		{ID: "g1", Name: "Fitness"},
		{ID: "g2", Name: "Work"},
	}
	m := newTestList(t, nil, groups)

	updated, _ := m.Update(msgs.GroupDeletedMsg{ID: "g1"})
	m = updated.(list.List)

	v := m.View()
	if strings.Contains(v.Content, "Fitness") {
		t.Error("View should not contain deleted group 'Fitness'")
	}
}

func TestWindowSizeMsg(t *testing.T) {
	m := newTestList(t, []*domain.Habit{{ID: "h1", Name: "Run"}}, nil)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(list.List)

	// Must not panic; verify View still works.
	v := m.View()
	if v.Content == "" {
		t.Error("expected non-empty view after WindowSizeMsg")
	}
}

func TestClearStatusMsg(t *testing.T) {
	m := newTestList(t, []*domain.Habit{{ID: "h1", Name: "Run"}}, nil)

	// Push an error to populate the error queue.
	updated, _ := m.Update(msgs.APIErrorMsg{Op: "test", Err: errors.New("fail")})
	m = updated.(list.List)

	// ClearStatusMsg should drain the queue without panic.
	updated, _ = m.Update(msgs.ClearStatusMsg{})
	m = updated.(list.List)

	v := m.View()
	if v.Content == "" {
		t.Error("expected non-empty view after ClearStatusMsg")
	}
}

func TestKeyA_PushesAddScreen(t *testing.T) {
	m := newTestList(t, []*domain.Habit{{ID: "h1", Name: "Run"}}, nil)

	_, cmd := m.Update(keyPress('a', "a"))
	if cmd == nil {
		t.Fatal("expected non-nil cmd from key 'a'")
	}

	msg := cmd()
	if _, ok := msg.(msgs.PushScreenMsg); !ok {
		t.Fatalf("expected PushScreenMsg, got %T", msg)
	}
}

func TestKeyG_PushesAddGroupScreen(t *testing.T) {
	m := newTestList(t, []*domain.Habit{{ID: "h1", Name: "Run"}}, nil)

	_, cmd := m.Update(keyPress('g', "G"))
	if cmd == nil {
		t.Fatal("expected non-nil cmd from key 'G'")
	}

	msg := cmd()
	if _, ok := msg.(msgs.PushScreenMsg); !ok {
		t.Fatalf("expected PushScreenMsg, got %T", msg)
	}
}

func TestKeyEnter_OnHabit_PushesShowScreen(t *testing.T) {
	habits := []*domain.Habit{{ID: "h1", Name: "Run"}}
	m := newTestList(t, habits, nil)

	_, cmd := m.Update(enterKey())
	if cmd == nil {
		t.Fatal("expected non-nil cmd from enter on habit")
	}

	msg := cmd()
	if _, ok := msg.(msgs.PushScreenMsg); !ok {
		t.Fatalf("expected PushScreenMsg from enter on habit, got %T", msg)
	}
}

func TestKeyEnter_OnGroup_PushesGroupShowScreen(t *testing.T) {
	// Create a list with only groups so the first selected item is a group.
	groups := []*domain.Group{{ID: "g1", Name: "Fitness"}}
	m := newTestList(t, nil, groups)

	_, cmd := m.Update(enterKey())
	if cmd == nil {
		t.Fatal("expected non-nil cmd from enter on group")
	}

	msg := cmd()
	if _, ok := msg.(msgs.PushScreenMsg); !ok {
		t.Fatalf("expected PushScreenMsg from enter on group, got %T", msg)
	}
}

func TestTextInputSubmit_RenameHabit(t *testing.T) {
	habits := []*domain.Habit{{ID: "h1", Name: "Run"}}
	m := newTestList(t, habits, nil)

	submit := textinput.Submit{Key: "habit", ID: "h1", Value: "Sprint"}
	updated, cmd := m.Update(submit)
	m = updated.(list.List)

	if cmd == nil {
		t.Fatal("expected non-nil cmd from rename")
	}

	v := m.View()
	if !strings.Contains(v.Content, "Sprint") {
		t.Errorf("View should contain renamed habit 'Sprint', got:\n%s", v.Content)
	}
}

func TestAPIErrorMsg_DeleteHabitRollback(t *testing.T) {
	habits := []*domain.Habit{{ID: "h1", Name: "Run"}}
	m := newTestList(t, habits, nil)

	// Verify habit is present initially.
	v := m.View()
	if !strings.Contains(v.Content, "Run") {
		t.Fatal("habit 'Run' should be in initial view")
	}

	// Press "D" to trigger optimistic delete on the selected habit.
	updated, cmd := m.Update(keyPress('d', "D"))
	m = updated.(list.List)
	if cmd == nil {
		t.Fatal("expected non-nil cmd from delete key")
	}

	// Item is removed from the list. The status message still contains
	// the name ("Deleting Run..."), so verify via "No items" indicator.
	v = m.View()
	if !strings.Contains(v.Content, "No items") {
		t.Error("list should show 'No items' after optimistic habit removal")
	}

	// Simulate API error for the delete operation.
	updated, _ = m.Update(msgs.APIErrorMsg{Op: "delete habit", ID: "h1", Err: errors.New("server error")})
	m = updated.(list.List)

	// After processing the batch cmd from InsertItem, the item is restored.
	// The error message with "restored" confirms rollback occurred.
	v = m.View()
	if !strings.Contains(v.Content, "restored") {
		t.Errorf("View should contain 'restored' after rollback, got:\n%s", v.Content)
	}
}

func TestAPIErrorMsg_DeleteGroupRollback(t *testing.T) {
	groups := []*domain.Group{{ID: "g1", Name: "Fitness"}}
	m := newTestList(t, nil, groups)

	// Verify group is present initially.
	v := m.View()
	if !strings.Contains(v.Content, "Fitness") {
		t.Fatal("group 'Fitness' should be in initial view")
	}

	// For groups, delete is bound to "d" (lowercase) on groupBinds.
	updated, cmd := m.Update(keyPress('d', "d"))
	m = updated.(list.List)
	if cmd == nil {
		t.Fatal("expected non-nil cmd from group delete key")
	}

	// Item is removed from the list. Status message "Deleting Fitness..."
	// still contains the name, so verify removal via "No items" indicator.
	v = m.View()
	if !strings.Contains(v.Content, "No items") {
		t.Error("list should show 'No items' after optimistic group removal")
	}

	// Simulate API error.
	updated, _ = m.Update(msgs.APIErrorMsg{Op: "delete group", ID: "g1", Err: errors.New("server error")})
	m = updated.(list.List)

	// The error message with "restored" confirms rollback occurred.
	v = m.View()
	if !strings.Contains(v.Content, "restored") {
		t.Errorf("View should contain 'restored' after rollback, got:\n%s", v.Content)
	}
}

func TestAPIErrorMsg_UpdateHabitRollback(t *testing.T) {
	habits := []*domain.Habit{{ID: "h1", Name: "Run"}}
	m := newTestList(t, habits, nil)

	// Optimistic rename via textinput.Submit.
	submit := textinput.Submit{Key: "habit", ID: "h1", Value: "Sprint"}
	updated, _ := m.Update(submit)
	m = updated.(list.List)

	v := m.View()
	if !strings.Contains(v.Content, "Sprint") {
		t.Fatal("View should contain optimistically renamed 'Sprint'")
	}

	// API error rolls back the rename.
	updated, _ = m.Update(msgs.APIErrorMsg{Op: "update habit", ID: "h1", Err: errors.New("server error")})
	m = updated.(list.List)

	v = m.View()
	if !strings.Contains(v.Content, "Run") {
		t.Errorf("habit should be rolled back to 'Run', got:\n%s", v.Content)
	}
}
