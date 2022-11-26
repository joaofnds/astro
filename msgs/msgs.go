package msgs

import tea "github.com/charmbracelet/bubbletea"

type Msg int

const MsgUpdateList Msg = iota

var UpdateList = Cmd(MsgUpdateList)

func Cmd(msg Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
