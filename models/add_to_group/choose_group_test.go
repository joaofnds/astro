package add_to_group_test

import (
	"astro/api"
	"astro/domain"
	"astro/models/add_to_group"
	"astro/msgs"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestChooseGroup(t *testing.T) add_to_group.ChooseGroup {
	t.Helper()
	client := api.NewClient("http://unused", "tok")
	habit := &domain.Habit{ID: "h1", Name: "Run"}
	groups := []*domain.Group{{ID: "g1", Name: "Fitness"}}
	return add_to_group.NewChooseGroup(client, habit, groups)
}

func TestView_ContainsTitle(t *testing.T) {
	m := newTestChooseGroup(t)

	// The list needs adequate dimensions to render the title.
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updated, _ := m.Update(sizeMsg)
	m = updated.(add_to_group.ChooseGroup)

	v := m.View()
	if !strings.Contains(v.Content, "Choose a group") {
		t.Fatalf("expected View to contain 'Choose a group', got %q", v.Content)
	}
}

func TestKeyEsc_ReturnsPopScreen(t *testing.T) {
	m := newTestChooseGroup(t)

	// Send WindowSizeMsg first so the list delegate has visible items.
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updated, _ := m.Update(sizeMsg)
	m = updated.(add_to_group.ChooseGroup)

	keyEsc := tea.KeyPressMsg(tea.Key{Code: tea.KeyEscape})
	_, cmd := m.Update(keyEsc)
	if cmd == nil {
		t.Fatal("expected non-nil cmd from esc")
	}
	msg := cmd()
	if _, ok := msg.(msgs.PopScreenMsg); !ok {
		t.Fatalf("expected PopScreenMsg from esc, got %T", msg)
	}
}

func TestKeyEnter_ReturnsPopWithAddToGroup(t *testing.T) {
	m := newTestChooseGroup(t)

	// Send WindowSizeMsg so items are visible and selectable.
	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updated, _ := m.Update(sizeMsg)
	m = updated.(add_to_group.ChooseGroup)

	keyEnter := tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter})
	_, cmd := m.Update(keyEnter)
	if cmd == nil {
		t.Fatal("expected non-nil cmd from enter with group selected")
	}
	msg := cmd()
	popMsg, ok := msg.(msgs.PopScreenMsg)
	if !ok {
		t.Fatalf("expected PopScreenMsg from enter, got %T", msg)
	}
	if popMsg.Cmd == nil {
		t.Fatal("expected PopScreenMsg to have follow-up Cmd (AddToGroup)")
	}
}

func TestWindowSizeMsg_DoesNotCrash(t *testing.T) {
	m := newTestChooseGroup(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updated, _ := m.Update(sizeMsg)
	_ = updated.(add_to_group.ChooseGroup) // Must not panic.
}
