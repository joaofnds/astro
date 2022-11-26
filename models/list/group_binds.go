package list

import "github.com/charmbracelet/bubbles/key"

type groupBinds struct {
	view key.Binding
}

func NewGroupBinds() groupBinds {
	return groupBinds{
		view: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view"),
		),
	}
}

func (k groupBinds) ToSlice() []key.Binding {
	return []key.Binding{k.view}
}
