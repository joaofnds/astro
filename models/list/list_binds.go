package list

import "github.com/charmbracelet/bubbles/key"

type listBinds struct {
	add      key.Binding
	addGroup key.Binding
}

func NewListBinds() listBinds {
	return listBinds{
		add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		addGroup: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "add group"),
		),
	}
}

func (k listBinds) ToSlice() []key.Binding {
	return []key.Binding{k.add, k.addGroup}
}
