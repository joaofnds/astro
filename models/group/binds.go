package group

import "github.com/charmbracelet/bubbles/key"

type binds struct {
	checkIn key.Binding
	rename  key.Binding
	delete  key.Binding
	view    key.Binding
	quit    key.Binding
	up      key.Binding
	down    key.Binding
	left    key.Binding
	right   key.Binding
	tab     key.Binding
}

func newBinds() binds {
	return binds{
		up: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "down"),
		),
		left: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "left"),
		),
		right: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "right"),
		),
		checkIn: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "check in"),
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
		quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q/esc", "quit"),
		),
		tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "switch"),
		),
	}
}

func (k binds) ToSlice() []key.Binding {
	return []key.Binding{k.view, k.rename, k.delete, k.tab, k.quit}
}
