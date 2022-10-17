package show

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	CheckIn key.Binding
	Delete  key.Binding
	Help    key.Binding
	Quit    key.Binding
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Down, k.Up, k.Right, k.CheckIn, k.Delete, k.Help, k.Quit}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.CheckIn, k.Delete, k.Help, k.Quit},
	}
}

func NewKeymap() keymap {
	return keymap{
		CheckIn: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "check in"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete activity"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		Up: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "prev day"),
		),
		Down: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "next day"),
		),
		Left: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "prev week"),
		),
		Right: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "next week"),
		),
	}
}
