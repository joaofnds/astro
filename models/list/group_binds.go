package list

import "github.com/charmbracelet/bubbles/key"

type groupBinds struct {
	add      key.Binding
	view     key.Binding
	delete   key.Binding
	addGroup key.Binding
}

func NewGroupBinds() groupBinds {
	return groupBinds{
		add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		view: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view"),
		),
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		addGroup: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "add group"),
		),
	}
}

func (k groupBinds) ShortHelp() []key.Binding {
	return []key.Binding{k.add, k.view, k.delete, k.addGroup}
}

func (k groupBinds) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.add, k.view, k.delete, k.addGroup}}
}
