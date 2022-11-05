package list

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	checkIn key.Binding
	add     key.Binding
	rename  key.Binding
	delete  key.Binding
	view    key.Binding
}

func NewKeymap() keymap {
	return keymap{
		checkIn: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "check in"),
		),
		add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		rename: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename"),
		),
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		view: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view"),
		),
	}
}

func (k keymap) ToSlice() []key.Binding {
	return []key.Binding{k.add, k.delete, k.view}
}
