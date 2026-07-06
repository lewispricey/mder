package keybinds

import tea "charm.land/bubbletea/v2"

func IsQuit(msg tea.KeyMsg) bool {
	s := msg.String()
	return s == "q" || s == "ctrl+c"
}
