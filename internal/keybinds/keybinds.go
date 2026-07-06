package keybinds

import tea "github.com/charmbracelet/bubbletea"

func IsHardQuit(msg tea.KeyMsg) bool {
	return msg.String() == "ctrl+c"
}

func IsQuit(msg tea.KeyMsg) bool {
	return msg.String() == "q" || IsHardQuit(msg)
}
