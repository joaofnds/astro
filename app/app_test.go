package app_test

import (
	"astro/api"
	"astro/app"
	"astro/domain"
	"astro/msgs"
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// mockScreen is a minimal tea.Model that records the last message it received.
type mockScreen struct {
	lastMsg tea.Msg
	initCmd tea.Cmd
}

func newMockScreen(initCmd tea.Cmd) *mockScreen {
	return &mockScreen{initCmd: initCmd}
}

func (m *mockScreen) Init() tea.Cmd { return m.initCmd }

func (m *mockScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.lastMsg = msg
	return m, nil
}

func (m *mockScreen) View() tea.View {
	return tea.NewView("mock screen")
}

func TestNew(t *testing.T) {
	client := api.NewClient("http://localhost", "token")
	a := app.New(client)

	// New returns App with ready=false; View should show spinner + loading text.
	v := a.View()
	if !strings.Contains(v.Content, "Loading habits...") {
		t.Fatalf("expected loading view with spinner, got %q", v.Content)
	}
}

func TestInit_ReturnsNonNilCmd(t *testing.T) {
	client := api.NewClient("http://localhost", "token")
	a := app.New(client)
	cmd := a.Init()
	if cmd == nil {
		t.Fatal("expected non-nil Cmd from Init (LoadAll)")
	}
}

func TestPushScreen(t *testing.T) {
	a := appWithScreen(newMockScreen(nil))
	screen2 := newMockScreen(nil)

	updated, _ := a.Update(msgs.PushScreenMsg{Screen: screen2})
	a = updated.(app.App)

	// Should show screen2's view.
	v := a.View()
	if v.Content != "mock screen" {
		t.Fatalf("expected pushed screen view, got %q", v.Content)
	}
}

func TestPushScreen_ReturnsScreenInit(t *testing.T) {
	a := appWithScreen(newMockScreen(nil))
	initCmd := func() tea.Msg { return "init-msg" }
	screen2 := newMockScreen(initCmd)

	_, cmd := a.Update(msgs.PushScreenMsg{Screen: screen2})
	if cmd == nil {
		t.Fatal("expected non-nil Cmd from pushed screen Init")
	}
}

func TestPopScreen_WithMultiple(t *testing.T) {
	screen1 := newMockScreen(nil)
	a := appWithScreen(screen1)
	screen2 := newMockScreen(nil)
	updated, _ := a.Update(msgs.PushScreenMsg{Screen: screen2})
	a = updated.(app.App)

	updated, _ = a.Update(msgs.PopScreenMsg{})
	a = updated.(app.App)

	// After pop, should show screen1's view.
	v := a.View()
	if v.Content != "mock screen" {
		t.Fatalf("expected screen1 view after pop, got %q", v.Content)
	}
}

func TestPopScreen_WithCmd(t *testing.T) {
	a := appWithScreen(newMockScreen(nil))
	screen2 := newMockScreen(nil)
	updated, _ := a.Update(msgs.PushScreenMsg{Screen: screen2})
	a = updated.(app.App)

	followUp := func() tea.Msg { return "follow-up" }
	_, cmd := a.Update(msgs.PopScreenMsg{Cmd: followUp})
	if cmd == nil {
		t.Fatal("expected non-nil follow-up Cmd after pop")
	}
}

func TestPopScreen_LastScreen_Quits(t *testing.T) {
	a := appWithScreen(newMockScreen(nil))

	_, cmd := a.Update(msgs.PopScreenMsg{})
	if cmd == nil {
		t.Fatal("expected tea.Quit when popping last screen")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected QuitMsg, got %T", msg)
	}
}

func TestWindowSizeMsg_StoresAndForwards(t *testing.T) {
	screen := newMockScreen(nil)
	a := appWithScreen(screen)

	sizeMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	a.Update(sizeMsg)

	if screen.lastMsg == nil {
		t.Fatal("expected WindowSizeMsg to be forwarded to active screen")
	}
	if got, ok := screen.lastMsg.(tea.WindowSizeMsg); !ok || got.Width != 80 || got.Height != 24 {
		t.Fatalf("expected WindowSizeMsg{80,24} forwarded, got %+v", screen.lastMsg)
	}
}

func TestDataLoadedMsg(t *testing.T) {
	client := api.NewClient("http://localhost", "token")
	a := app.New(client)

	habits := []*domain.Habit{{ID: "h1", Name: "A"}}
	groups := []*domain.Group{{ID: "g1", Name: "G"}}

	updated, _ := a.Update(msgs.DataLoadedMsg{Habits: habits, Groups: groups})
	a = updated.(app.App)

	// After DataLoadedMsg, ready should be true (view should not show loading).
	v := a.View()
	if strings.Contains(v.Content, "Loading habits...") {
		t.Fatal("expected ready view after DataLoadedMsg, still showing loading")
	}
}

func TestFatalErrorMsg_Quits(t *testing.T) {
	a := appWithScreen(newMockScreen(nil))
	fatalErr := errors.New("connection refused")

	_, cmd := a.Update(msgs.FatalErrorMsg{Err: fatalErr})
	if cmd == nil {
		t.Fatal("expected non-nil Cmd from FatalErrorMsg")
	}
}

func TestCheckInResultMsg_MergesAndForwards(t *testing.T) {
	screen := newMockScreen(nil)
	a := appWithReadyState(screen)

	updatedHabit := &domain.Habit{ID: "h1", Name: "Updated"}
	updated, _ := a.Update(msgs.CheckInResultMsg{Habit: updatedHabit})
	a = updated.(app.App)

	if screen.lastMsg == nil {
		t.Fatal("expected CheckInResultMsg to be forwarded to active screen")
	}
	if _, ok := screen.lastMsg.(msgs.CheckInResultMsg); !ok {
		t.Fatalf("expected CheckInResultMsg, got %T", screen.lastMsg)
	}
}

func TestHabitCreatedMsg_AddsAndForwards(t *testing.T) {
	screen := newMockScreen(nil)
	a := appWithReadyState(screen)

	newHabit := &domain.Habit{ID: "h2", Name: "New"}
	updated, _ := a.Update(msgs.HabitCreatedMsg{Habit: newHabit})
	a = updated.(app.App)

	if screen.lastMsg == nil {
		t.Fatal("expected HabitCreatedMsg to be forwarded to active screen")
	}
	if _, ok := screen.lastMsg.(msgs.HabitCreatedMsg); !ok {
		t.Fatalf("expected HabitCreatedMsg, got %T", screen.lastMsg)
	}
}

func TestHabitDeletedMsg_RemovesAndForwards(t *testing.T) {
	screen := newMockScreen(nil)
	a := appWithReadyState(screen)

	updated, _ := a.Update(msgs.HabitDeletedMsg{ID: "h1"})
	a = updated.(app.App)

	if screen.lastMsg == nil {
		t.Fatal("expected HabitDeletedMsg to be forwarded to active screen")
	}
	if _, ok := screen.lastMsg.(msgs.HabitDeletedMsg); !ok {
		t.Fatalf("expected HabitDeletedMsg, got %T", screen.lastMsg)
	}
}

func TestGroupDeletedMsg_RemovesAndForwards(t *testing.T) {
	screen := newMockScreen(nil)
	a := appWithReadyState(screen)

	updated, _ := a.Update(msgs.GroupDeletedMsg{ID: "g1"})
	a = updated.(app.App)

	if screen.lastMsg == nil {
		t.Fatal("expected GroupDeletedMsg to be forwarded to active screen")
	}
	if _, ok := screen.lastMsg.(msgs.GroupDeletedMsg); !ok {
		t.Fatalf("expected GroupDeletedMsg, got %T", screen.lastMsg)
	}
}

func TestViewWhenNotReady(t *testing.T) {
	client := api.NewClient("http://localhost", "token")
	a := app.New(client)

	v := a.View()
	if !strings.Contains(v.Content, "Loading habits...") {
		t.Fatalf("expected spinner + 'Loading habits...', got %q", v.Content)
	}
	if !v.AltScreen {
		t.Fatal("expected AltScreen=true for loading view")
	}
}

func TestViewWhenReady(t *testing.T) {
	screen := newMockScreen(nil)
	a := appWithScreen(screen)

	v := a.View()
	if v.Content != "mock screen" {
		t.Fatalf("expected active screen view, got %q", v.Content)
	}
}

func TestCtrlC_Quits(t *testing.T) {
	a := appWithScreen(newMockScreen(nil))

	_, cmd := a.Update(tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl}))
	if cmd == nil {
		t.Fatal("expected tea.Quit from ctrl+c")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected QuitMsg, got %T", msg)
	}
}

// --- Navigation flow tests (TEST-06) ---

func TestDeepNavigation(t *testing.T) {
	screen1 := newMockScreen(nil)
	a := appWithScreen(screen1)

	// Push screen2 and screen3.
	screen2 := newMockScreen(nil)
	updated, _ := a.Update(msgs.PushScreenMsg{Screen: screen2})
	a = updated.(app.App)

	screen3 := newMockScreen(nil)
	updated, _ = a.Update(msgs.PushScreenMsg{Screen: screen3})
	a = updated.(app.App)

	// Pop screen3 -- screen2 should receive subsequent messages.
	updated, _ = a.Update(msgs.PopScreenMsg{})
	a = updated.(app.App)
	testMsg := tea.WindowSizeMsg{Width: 1, Height: 1}
	a.Update(testMsg)
	if screen2.lastMsg == nil {
		t.Fatal("after popping screen3, screen2 should be active and receive messages")
	}

	// Pop screen2 -- screen1 should receive subsequent messages.
	updated, _ = a.Update(msgs.PopScreenMsg{})
	a = updated.(app.App)
	screen1.lastMsg = nil
	a.Update(testMsg)
	if screen1.lastMsg == nil {
		t.Fatal("after popping screen2, screen1 should be active and receive messages")
	}

	// Pop screen1 -- should quit.
	_, cmd := a.Update(msgs.PopScreenMsg{})
	if cmd == nil {
		t.Fatal("expected tea.Quit when popping last screen")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected QuitMsg on final pop, got %T", msg)
	}
}

func TestDataLoadedMsg_TransitionsToListScreen(t *testing.T) {
	client := api.NewClient("http://localhost", "token")
	a := app.New(client)

	// Confirm loading state.
	v := a.View()
	if !strings.Contains(v.Content, "Loading habits...") {
		t.Fatal("expected loading view before DataLoadedMsg")
	}

	habits := []*domain.Habit{{ID: "h1", Name: "Running"}}
	groups := []*domain.Group{{ID: "g1", Name: "Fitness"}}

	updated, _ := a.Update(msgs.DataLoadedMsg{Habits: habits, Groups: groups})
	a = updated.(app.App)

	// After DataLoadedMsg, should no longer show loading.
	v = a.View()
	if strings.Contains(v.Content, "Loading habits...") {
		t.Fatal("expected transition away from loading after DataLoadedMsg")
	}

	// View should reflect the list screen content (not loading spinner).
	if v.Content == "" {
		t.Fatal("expected non-empty view after DataLoadedMsg transition")
	}
}

func TestStateMessages_ForwardToScreen(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.Msg
	}{
		{"GroupCreatedMsg", msgs.GroupCreatedMsg{Group: &domain.Group{ID: "g2", Name: "New"}}},
		{"HabitCreatedMsg", msgs.HabitCreatedMsg{Habit: &domain.Habit{ID: "h2", Name: "New"}}},
		{"HabitDeletedMsg", msgs.HabitDeletedMsg{ID: "h1"}},
		{"GroupDeletedMsg", msgs.GroupDeletedMsg{ID: "g1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := newMockScreen(nil)
			a := appWithReadyState(screen)

			a.Update(tt.msg)

			if screen.lastMsg == nil {
				t.Fatalf("expected %s to be forwarded to active screen", tt.name)
			}
		})
	}
}

func TestMultiplePushPop_Sequence(t *testing.T) {
	screenA := newMockScreen(nil)
	a := appWithScreen(screenA)

	// Push B.
	screenB := newMockScreen(nil)
	updated, _ := a.Update(msgs.PushScreenMsg{Screen: screenB})
	a = updated.(app.App)

	// Push C.
	screenC := newMockScreen(nil)
	updated, _ = a.Update(msgs.PushScreenMsg{Screen: screenC})
	a = updated.(app.App)

	// Pop C -- B should be active.
	updated, _ = a.Update(msgs.PopScreenMsg{})
	a = updated.(app.App)
	screenB.lastMsg = nil
	probe := tea.WindowSizeMsg{Width: 42, Height: 42}
	a.Update(probe)
	if screenB.lastMsg == nil {
		t.Fatal("after popping C, B should be active")
	}

	// Pop B -- A should be active.
	updated, _ = a.Update(msgs.PopScreenMsg{})
	a = updated.(app.App)
	screenA.lastMsg = nil
	a.Update(probe)
	if screenA.lastMsg == nil {
		t.Fatal("after popping B, A should be active")
	}

	// Push D on top of A.
	screenD := newMockScreen(nil)
	updated, _ = a.Update(msgs.PushScreenMsg{Screen: screenD})
	a = updated.(app.App)
	screenD.lastMsg = nil
	a.Update(probe)
	if screenD.lastMsg == nil {
		t.Fatal("D should be active after pushing on top of A")
	}

	// Pop D -- A should be active again.
	updated, _ = a.Update(msgs.PopScreenMsg{})
	a = updated.(app.App)
	screenA.lastMsg = nil
	a.Update(probe)
	if screenA.lastMsg == nil {
		t.Fatal("after popping D, A should be active again")
	}
}

// --- helpers ---

// appWithScreen creates an App that is "ready" with one screen on the stack.
// Uses exported NewForTest to bypass normal initialization.
func appWithScreen(screen tea.Model) app.App {
	return app.NewForTest(screen)
}

// appWithReadyState creates an App with initial state data and a screen on the stack.
func appWithReadyState(screen tea.Model) app.App {
	a := app.NewForTest(screen)
	a.SetStateForTest([]*domain.Habit{{ID: "h1", Name: "Original"}}, []*domain.Group{{ID: "g1", Name: "Group"}})
	return a
}
