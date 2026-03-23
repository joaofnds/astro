package textinput_test

import (
	"astro/models/textinput"
	"astro/msgs"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestView_ContainsPrompt(t *testing.T) {
	m := textinput.New("Enter habit name", "", "habit", "h1", 80)
	v := m.View()
	if !strings.Contains(v.Content, "Enter habit name") {
		t.Fatalf("expected View to contain prompt text, got %q", v.Content)
	}
}

func TestView_ContainsInitialValue(t *testing.T) {
	m := textinput.New("Edit", "existing value", "edit", "e1", 80)
	v := m.View()
	if !strings.Contains(v.Content, "existing value") {
		t.Fatalf("expected View to contain initial value, got %q", v.Content)
	}
}

func TestKeyEsc_ReturnsPopScreen(t *testing.T) {
	m := textinput.New("prompt", "", "key", "id", 80)
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

func TestKeyCtrlC_ReturnsPopScreen(t *testing.T) {
	m := textinput.New("prompt", "", "key", "id", 80)
	keyCtrlC := tea.KeyPressMsg(tea.Key{Code: 'c', Mod: tea.ModCtrl})
	_, cmd := m.Update(keyCtrlC)
	if cmd == nil {
		t.Fatal("expected non-nil cmd from ctrl+c")
	}
	msg := cmd()
	if _, ok := msg.(msgs.PopScreenMsg); !ok {
		t.Fatalf("expected PopScreenMsg from ctrl+c, got %T", msg)
	}
}

func TestKeyEnter_WithValue_ReturnsPopWithSubmit(t *testing.T) {
	m := textinput.New("prompt", "hello world", "habit", "h1", 80)
	keyEnter := tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter})
	_, cmd := m.Update(keyEnter)
	if cmd == nil {
		t.Fatal("expected non-nil cmd from enter with value")
	}

	msg := cmd()
	popMsg, ok := msg.(msgs.PopScreenMsg)
	if !ok {
		t.Fatalf("expected PopScreenMsg from enter, got %T", msg)
	}
	if popMsg.Cmd == nil {
		t.Fatal("expected PopScreenMsg to have follow-up Cmd with Submit")
	}

	submitMsg := popMsg.Cmd()
	submit, ok := submitMsg.(textinput.Submit)
	if !ok {
		t.Fatalf("expected Submit message from Cmd, got %T", submitMsg)
	}
	if submit.Key != "habit" {
		t.Fatalf("expected Key='habit', got %q", submit.Key)
	}
	if submit.ID != "h1" {
		t.Fatalf("expected ID='h1', got %q", submit.ID)
	}
	if submit.Value != "hello world" {
		t.Fatalf("expected Value='hello world', got %q", submit.Value)
	}
}

func TestKeyEnter_Empty_DoesNotPop(t *testing.T) {
	m := textinput.New("prompt", "", "key", "id", 80)
	keyEnter := tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter})
	_, cmd := m.Update(keyEnter)

	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(msgs.PopScreenMsg); ok {
			t.Fatal("expected no PopScreenMsg from enter with empty value")
		}
	}
}

func TestInit_ReturnsNonNilCmd(t *testing.T) {
	m := textinput.New("prompt", "", "key", "id", 80)
	cmd := m.Init()
	if cmd == nil {
		t.Fatal("expected non-nil Cmd from Init (textarea.Blink)")
	}
}
