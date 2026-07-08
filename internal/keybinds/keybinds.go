package keybinds

import tea "github.com/charmbracelet/bubbletea"

func IsHardQuit(msg tea.KeyMsg) bool {
	return msg.String() == "ctrl+c"
}

func IsSave(msg tea.KeyMsg) bool {
	return msg.String() == "ctrl+s"
}

func IsToggleMode(msg tea.KeyMsg) bool {
	return msg.String() == "ctrl+e"
}

func IsQuit(msg tea.KeyMsg) bool {
	return msg.String() == "q" || IsHardQuit(msg)
}
