package msgs

import tea "charm.land/bubbletea/v2"

type Msg int

const MsgUpdateList Msg = iota

var UpdateList = Cmd(MsgUpdateList)

func Cmd(msg Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
