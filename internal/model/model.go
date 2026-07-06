package model

import (
	"os"

	"charm.land/bubbletea/v2"
	"github.com/yourname/mded/internal/keybinds"
)

type Mode int

const (
	ViewMode Mode = iota
	EditMode
)

type fileLoadedMsg struct {
	content string
	err     error
}

type Model struct {
	mode     Mode
	filePath string
	content  string
	readErr  error
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
	return func() tea.Msg {
		data, err := os.ReadFile(m.filePath)
		return fileLoadedMsg{content: string(data), err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fileLoadedMsg:
		if msg.err != nil {
			m.readErr = msg.err
			return m, tea.Quit
		}
		m.content = msg.content
		return m, nil
	case tea.KeyMsg:
		if keybinds.IsQuit(msg) {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	var s string
	switch {
	case m.readErr != nil:
		s = "Error reading file.\n\nPress q to quit."
	case m.content == "":
		s = "Loading...\n\nPress q to quit."
	default:
		s = m.content
	}
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}
