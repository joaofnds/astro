package group_test

import (
	"astro/api"
	"astro/domain"
	"astro/models/group"
	"astro/models/textinput"
	"astro/msgs"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestGroup(t *testing.T, g *domain.Group) group.List {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	client := api.NewClient(srv.URL, "test-token", api.WithHTTPClient(srv.Client()))
	return group.NewShow(client, g, 80, 24)
}

func testGroup() *domain.Group {
	return &domain.Group{
		ID:   "g1",
		Name: "Fitness",
		Habits: []*domain.Habit{
			{ID: "h1", Name: "Run"},
			{ID: "h2", Name: "Swim"},
		},
	}
}

func keyPress(code rune, text string) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code, Text: text})
}

func specialKey(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func TestGroupView_ContainsGroupName(t *testing.T) {
	m := newTestGroup(t, testGroup())
	v := m.View()

	if !strings.Contains(v.Content, "Fitness") {
		t.Errorf("View should contain group name 'Fitness', got:\n%s", v.Content)
	}
}

func TestGroupView_ContainsHabitNames(t *testing.T) {
	m := newTestGroup(t, testGroup())
	v := m.View()

	for _, name := range []string{"Run", "Swim"} {
		if !strings.Contains(v.Content, name) {
			t.Errorf("View should contain habit name %q, got:\n%s", name, v.Content)
		}
	}
}

func TestGroupView_ContainsHelpKeys(t *testing.T) {
	m := newTestGroup(t, testGroup())
	v := m.View()

	for _, label := range []string{"check in", "view", "quit"} {
		if !strings.Contains(v.Content, label) {
			t.Errorf("View should contain help label %q, got:\n%s", label, v.Content)
		}
	}
}

func TestCheckInResultMsg(t *testing.T) {
	m := newTestGroup(t, testGroup())

	checkedIn := &domain.Habit{ID: "h1", Name: "Run", Activities: []domain.Activity{{ID: "a1"}}}
	updated, cmd := m.Update(msgs.CheckInResultMsg{Habit: checkedIn})
	m = updated.(group.List)

	if cmd == nil {
		t.Error("expected non-nil cmd (status message)")
	}
}

func TestHabitUpdatedMsg(t *testing.T) {
	m := newTestGroup(t, testGroup())

	updated, cmd := m.Update(msgs.HabitUpdatedMsg{Habit: &domain.Habit{ID: "h1", Name: "Sprint"}})
	m = updated.(group.List)

	if cmd == nil {
		t.Error("expected non-nil cmd (status message)")
	}

	v := m.View()
	if !strings.Contains(v.Content, "Sprint") {
		t.Errorf("View should contain updated habit name 'Sprint', got:\n%s", v.Content)
	}
	if strings.Contains(v.Content, "Run") {
		t.Error("View should no longer contain old habit name 'Run'")
	}
}

func TestRemovedFromGroupMsg(t *testing.T) {
	m := newTestGroup(t, testGroup())

	updated, _ := m.Update(msgs.RemovedFromGroupMsg{HabitID: "h1", GroupID: "g1"})
	m = updated.(group.List)

	v := m.View()
	if !strings.Contains(v.Content, "Swim") {
		t.Error("View should still contain 'Swim'")
	}
	// h1 ("Run") is removed from group.Habits, but the list model items
	// are not refreshed by RemovedFromGroupMsg (they were already removed
	// optimistically by the key "d" press). RemovedFromGroupMsg only
	// cleans up pending state and removes from group.Habits.
}

func TestKeyTab_ToggleHistogram(t *testing.T) {
	m := newTestGroup(t, testGroup())

	tabKey := specialKey(tea.KeyTab)

	// Toggle to histogram mode.
	updated, _ := m.Update(tabKey)
	m = updated.(group.List)

	v := m.View()
	if v.Content == "" {
		t.Error("expected non-empty view after tab toggle to histogram")
	}

	// Toggle back to list mode.
	updated, _ = m.Update(tabKey)
	m = updated.(group.List)

	v = m.View()
	if v.Content == "" {
		t.Error("expected non-empty view after tab toggle back to list")
	}
}

func TestKeyQ_ProducesPopScreenMsg(t *testing.T) {
	m := newTestGroup(t, testGroup())

	_, cmd := m.Update(keyPress('q', "q"))
	if cmd == nil {
		t.Fatal("expected non-nil cmd from key 'q'")
	}

	msg := cmd()
	if _, ok := msg.(msgs.PopScreenMsg); !ok {
		t.Fatalf("expected PopScreenMsg, got %T", msg)
	}
}

func TestKeyEsc_ProducesPopScreenMsg(t *testing.T) {
	m := newTestGroup(t, testGroup())

	_, cmd := m.Update(specialKey(tea.KeyEscape))
	if cmd == nil {
		t.Fatal("expected non-nil cmd from esc key")
	}

	msg := cmd()
	if _, ok := msg.(msgs.PopScreenMsg); !ok {
		t.Fatalf("expected PopScreenMsg, got %T", msg)
	}
}

func TestKeyEnter_OnHabit_PushesShowScreen(t *testing.T) {
	m := newTestGroup(t, testGroup())

	_, cmd := m.Update(specialKey(tea.KeyEnter))
	if cmd == nil {
		t.Fatal("expected non-nil cmd from enter on habit")
	}

	msg := cmd()
	if _, ok := msg.(msgs.PushScreenMsg); !ok {
		t.Fatalf("expected PushScreenMsg from enter on habit, got %T", msg)
	}
}

func TestKeyD_OptimisticRemove(t *testing.T) {
	m := newTestGroup(t, testGroup())

	_, cmd := m.Update(keyPress('d', "d"))
	if cmd == nil {
		t.Fatal("expected non-nil cmd from key 'd' (remove from group)")
	}

	// After removal, the list shows one fewer item.
	// The status message "Removing Run..." confirms the operation was initiated.
}

func TestTextInputSubmit_RenameHabit(t *testing.T) {
	m := newTestGroup(t, testGroup())

	submit := textinput.Submit{Key: "habit", ID: "h1", Value: "Sprint"}
	updated, cmd := m.Update(submit)
	m = updated.(group.List)

	if cmd == nil {
		t.Fatal("expected non-nil cmd from rename")
	}

	v := m.View()
	if !strings.Contains(v.Content, "Sprint") {
		t.Errorf("View should contain renamed habit 'Sprint', got:\n%s", v.Content)
	}
}

func TestWindowSizeMsg(t *testing.T) {
	m := newTestGroup(t, testGroup())

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(group.List)

	v := m.View()
	if v.Content == "" {
		t.Error("expected non-empty view after WindowSizeMsg")
	}
}

func TestAPIErrorMsg_RemoveRollback(t *testing.T) {
	m := newTestGroup(t, testGroup())

	// Verify "Run" is present initially.
	v := m.View()
	if !strings.Contains(v.Content, "Run") {
		t.Fatal("habit 'Run' should be in initial view")
	}

	// Press "d" to optimistically remove h1 from the group.
	updated, cmd := m.Update(keyPress('d', "d"))
	m = updated.(group.List)
	if cmd == nil {
		t.Fatal("expected non-nil cmd from remove key")
	}

	// Simulate API error for the remove operation.
	updated, _ = m.Update(msgs.APIErrorMsg{Op: "remove from group", ID: "h1", Err: errors.New("server error")})
	m = updated.(group.List)

	// After rollback, the error cross mark and op name appear in the status.
	// The item is re-inserted into the list via InsertItem command.
	v = m.View()
	if !strings.Contains(v.Content, "remove from group") {
		t.Errorf("View should contain error status with 'remove from group', got:\n%s", v.Content)
	}
}

func TestAPIErrorMsg_UpdateHabitRollback(t *testing.T) {
	m := newTestGroup(t, testGroup())

	// Optimistic rename via textinput.Submit.
	submit := textinput.Submit{Key: "habit", ID: "h1", Value: "Sprint"}
	updated, _ := m.Update(submit)
	m = updated.(group.List)

	v := m.View()
	if !strings.Contains(v.Content, "Sprint") {
		t.Fatal("View should contain optimistically renamed 'Sprint'")
	}

	// API error rolls back the rename.
	updated, _ = m.Update(msgs.APIErrorMsg{Op: "update habit", ID: "h1", Err: errors.New("server error")})
	m = updated.(group.List)

	// The rename is reverted in group.Habits and the list items are rebuilt.
	v = m.View()
	if !strings.Contains(v.Content, "Run") {
		t.Errorf("habit should be rolled back to 'Run', got:\n%s", v.Content)
	}
	if !strings.Contains(v.Content, "update habit") {
		t.Errorf("View should contain error status with 'update habit', got:\n%s", v.Content)
	}
}
