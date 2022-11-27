package list

import "github.com/charmbracelet/bubbles/key"

type habitBinds struct {
	checkIn    key.Binding
	add        key.Binding
	rename     key.Binding
	delete     key.Binding
	addGroup   key.Binding
	addToGroup key.Binding
	view       key.Binding
}

func NewHabitBinds() habitBinds {
	return habitBinds{
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
			key.WithKeys("D"),
			key.WithHelp("D", "delete"),
		),
		addGroup: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "add group"),
		),
		addToGroup: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "add to group"),
		),
		view: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view"),
		),
	}
}

func (k habitBinds) ShortHelp() []key.Binding {
	return []key.Binding{k.add, k.view, k.rename, k.delete, k.addGroup}
}

func (k habitBinds) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.checkIn, k.add, k.rename, k.delete, k.addGroup, k.addToGroup, k.view},
	}
}
