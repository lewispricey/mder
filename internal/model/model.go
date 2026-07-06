package model

import (
	"charm.land/bubbletea/v2"
	"github.com/yourname/mded/internal/keybinds"
)

type Mode int

const (
	ViewMode Mode = iota
	EditMode
)

type Model struct {
	mode     Mode
	filePath string
	width    int
	height   int
}

func New(mode Mode, filePath string) Model {
	return Model{
		mode:     mode,
		filePath: filePath,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if keybinds.IsQuit(msg) {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	v := tea.NewView("mded — press q to quit")
	v.AltScreen = true
	return v
}
