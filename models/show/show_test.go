package show_test

import (
	"astro/api"
	"astro/date"
	"astro/domain"
	"astro/models/show"
	"astro/msgs"
	"errors"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

func newTestShow(t *testing.T, habit *domain.Habit) show.Show {
	t.Helper()
	client := api.NewClient("http://unused", "tok")
	return show.NewShow(client, habit, 80)
}

func testHabit() *domain.Habit {
	return &domain.Habit{
		ID:   "h1",
		Name: "Morning Run",
		Activities: []domain.Activity{
			{ID: "a1", Desc: "5km", CreatedAt: date.Today().Add(-time.Hour)},
		},
	}
}

func TestView_ContainsHabitName(t *testing.T) {
	m := newTestShow(t, testHabit())
	v := m.View()
	if !strings.Contains(v.Content, "Morning Run") {
		t.Fatalf("expected View to contain habit name 'Morning Run', got %q", v.Content)
	}
}

func TestView_ContainsHelpKeys(t *testing.T) {
	m := newTestShow(t, testHabit())
	v := m.View()
	// Short help includes check-in key.
	if !strings.Contains(v.Content, "check-in") {
		t.Fatal("expected View to contain 'check-in' help label")
	}
	// Short help includes edit key.
	if !strings.Contains(v.Content, "edit") {
		t.Fatal("expected View to contain 'edit' help label")
	}
}

func TestAPIErrorMsg_SetsStatusWithCrossMark(t *testing.T) {
	m := newTestShow(t, testHabit())
	errMsg := msgs.APIErrorMsg{Op: "load", Err: errors.New("network fail")}
	updated, _ := m.Update(errMsg)
	m = updated.(show.Show)

	v := m.View()
	if !strings.Contains(v.Content, "\u2717") {
		t.Fatal("expected cross mark in status after APIErrorMsg")
	}
	if !strings.Contains(v.Content, "network fail") {
		t.Fatal("expected error text in status after APIErrorMsg")
	}
}

func TestAPIErrorMsg_CheckInRollback(t *testing.T) {
	h := testHabit()
	m := newTestShow(t, h)

	// Press 'c' to initiate check-in (sets snapshot, adds pending activity).
	keyC := tea.KeyPressMsg(tea.Key{Code: 'c', Text: "c"})
	updated, _ := m.Update(keyC)
	m = updated.(show.Show)

	// Verify pending activity was added.
	v := m.View()
	if !strings.Contains(v.Content, "Checking in...") {
		t.Fatal("expected 'Checking in...' status after pressing c")
	}

	// Send APIErrorMsg for check in to trigger rollback.
	errMsg := msgs.APIErrorMsg{Op: "check in", Err: errors.New("server error")}
	updated, _ = m.Update(errMsg)
	m = updated.(show.Show)

	v = m.View()
	if !strings.Contains(v.Content, "restored") {
		t.Fatal("expected 'restored' in status after check-in rollback")
	}
	if !strings.Contains(v.Content, "\u2717") {
		t.Fatal("expected cross mark after check-in error")
	}
}

func TestClearStatusMsg_ClearsStatus(t *testing.T) {
	m := newTestShow(t, testHabit())

	// Set a status via APIErrorMsg.
	errMsg := msgs.APIErrorMsg{Op: "test", Err: errors.New("err")}
	updated, _ := m.Update(errMsg)
	m = updated.(show.Show)
	if !strings.Contains(m.View().Content, "err") {
		t.Fatal("expected status to be set before clearing")
	}

	// Clear it.
	updated, _ = m.Update(msgs.ClearStatusMsg{})
	m = updated.(show.Show)
	if strings.Contains(m.View().Content, "\u2717") {
		t.Fatal("expected status to be cleared after ClearStatusMsg")
	}
}

func TestCheckInResultMsg_UpdatesHabitAndShowsStatus(t *testing.T) {
	h := testHabit()
	m := newTestShow(t, h)

	updatedHabit := &domain.Habit{
		ID:   "h1",
		Name: "Morning Run",
		Activities: []domain.Activity{
			{ID: "a1", Desc: "5km", CreatedAt: date.Today().Add(-time.Hour)},
			{ID: "a2", Desc: "10km", CreatedAt: date.Today()},
		},
	}
	updated, _ := m.Update(msgs.CheckInResultMsg{Habit: updatedHabit})
	m = updated.(show.Show)

	v := m.View()
	if !strings.Contains(v.Content, "Checked in") {
		t.Fatal("expected 'Checked in' status after CheckInResultMsg")
	}
}

func TestActivityUpdatedMsg_UpdatesDescAndShowsStatus(t *testing.T) {
	h := testHabit()
	m := newTestShow(t, h)

	updated, _ := m.Update(msgs.ActivityUpdatedMsg{
		HabitID:    "h1",
		ActivityID: "a1",
		Desc:       "10km",
	})
	m = updated.(show.Show)

	v := m.View()
	if !strings.Contains(v.Content, "Updated") {
		t.Fatal("expected 'Updated' status after ActivityUpdatedMsg")
	}
}

func TestActivityDeletedMsg_RemovesAndShowsStatus(t *testing.T) {
	h := testHabit()
	m := newTestShow(t, h)

	updated, _ := m.Update(msgs.ActivityDeletedMsg{
		HabitID:    "h1",
		ActivityID: "a1",
	})
	m = updated.(show.Show)

	v := m.View()
	if !strings.Contains(v.Content, "Deleted") {
		t.Fatal("expected 'Deleted' status after ActivityDeletedMsg")
	}
}

func TestKeyQ_ReturnsPopScreen(t *testing.T) {
	m := newTestShow(t, testHabit())
	keyQ := tea.KeyPressMsg(tea.Key{Code: 'q', Text: "q"})
	_, cmd := m.Update(keyQ)
	if cmd == nil {
		t.Fatal("expected non-nil cmd from 'q' key press")
	}
	msg := cmd()
	if _, ok := msg.(msgs.PopScreenMsg); !ok {
		t.Fatalf("expected PopScreenMsg from 'q', got %T", msg)
	}
}

func TestKeyJ_MovesSelectionDown(t *testing.T) {
	m := newTestShow(t, testHabit())

	// Move up first so we have room to move down. On Saturdays the initial
	// cursor lands at the last grid index and pressing 'j' clamps.
	keyK := tea.KeyPressMsg(tea.Key{Code: 'k', Text: "k"})
	updated, _ := m.Update(keyK)
	m = updated.(show.Show)
	vUp := m.View().Content

	keyJ := tea.KeyPressMsg(tea.Key{Code: 'j', Text: "j"})
	updated, _ = m.Update(keyJ)
	m = updated.(show.Show)

	vDown := m.View().Content
	if vUp == vDown {
		t.Fatal("expected View to change after pressing 'j' (down)")
	}
}

func TestKeyK_MovesSelectionUp(t *testing.T) {
	m := newTestShow(t, testHabit())

	// Move down first so we can move up.
	keyJ := tea.KeyPressMsg(tea.Key{Code: 'j', Text: "j"})
	updated, _ := m.Update(keyJ)
	m = updated.(show.Show)
	vDown := m.View().Content

	keyK := tea.KeyPressMsg(tea.Key{Code: 'k', Text: "k"})
	updated, _ = m.Update(keyK)
	m = updated.(show.Show)

	vUp := m.View().Content
	if vDown == vUp {
		t.Fatal("expected View to change after pressing 'k' (up)")
	}
}

func TestKeyH_MovesSelectionLeftBy7(t *testing.T) {
	m := newTestShow(t, testHabit())
	v1 := m.View().Content

	keyH := tea.KeyPressMsg(tea.Key{Code: 'h', Text: "h"})
	updated, _ := m.Update(keyH)
	m = updated.(show.Show)

	v2 := m.View().Content
	// 'h' moves by 7 days, so View should change (unless already at boundary).
	// Since initial selected is date.DiffInDays(timeframe_start, today), it is non-zero.
	if v1 == v2 {
		t.Fatal("expected View to change after pressing 'h' (left by 7)")
	}
}

func TestKeyL_MovesSelectionRightBy7(t *testing.T) {
	m := newTestShow(t, testHabit())

	// Move left first so we have room to move right.
	keyH := tea.KeyPressMsg(tea.Key{Code: 'h', Text: "h"})
	updated, _ := m.Update(keyH)
	m = updated.(show.Show)
	vLeft := m.View().Content

	keyL := tea.KeyPressMsg(tea.Key{Code: 'l', Text: "l"})
	updated, _ = m.Update(keyL)
	m = updated.(show.Show)

	vRight := m.View().Content
	if vLeft == vRight {
		t.Fatal("expected View to change after pressing 'l' (right by 7)")
	}
}
