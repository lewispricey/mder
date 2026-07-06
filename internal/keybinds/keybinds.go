package keybinds

import tea "github.com/charmbracelet/bubbletea"

func IsQuit(msg tea.KeyMsg) bool {
	s := msg.String()
	return s == "q" || s == "ctrl+c"
}
