package msgs

import tea "charm.land/bubbletea/v2"

// Deprecated: Msg is the old message type. Use typed messages instead.
// This exists only for backward compatibility during the migration.
type Msg int

// Deprecated: MsgUpdateList is the old update-list signal.
const MsgUpdateList Msg = iota

// Deprecated: UpdateList is the old update-list command.
var UpdateList = Cmd(MsgUpdateList)

// Deprecated: Cmd wraps an old Msg into a tea.Cmd.
func Cmd(msg Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
