package show

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	CheckIn  key.Binding
	VCheckIn key.Binding
	Edit     key.Binding
	Delete   key.Binding
	Help     key.Binding
	Quit     key.Binding
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Down, k.Up, k.Right, k.CheckIn, k.Edit, k.Delete}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Down, k.Up, k.Right},
		{k.CheckIn, k.VCheckIn, k.Edit, k.Delete},
		{k.Help, k.Quit},
	}
}

func NewKeymap() keymap {
	return keymap{
		CheckIn: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "check-in"),
		),
		VCheckIn: key.NewBinding(
			key.WithKeys("C"),
			key.WithHelp("C", "verbose check-in"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
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
			key.WithHelp("q", "go back"),
		),
		Up: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "↑"),
		),
		Down: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "↓"),
		),
		Left: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "←"),
		),
		Right: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "→"),
		),
	}
}
