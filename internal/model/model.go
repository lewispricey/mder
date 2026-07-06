package model

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lewispricey/mded/internal/keybinds"
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

func (m Model) Width() int  { return m.width }
func (m Model) Height() int { return m.height }

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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if keybinds.IsQuit(msg) {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch {
	case m.readErr != nil:
		return fmt.Sprintf("Error: %v\n\nPress q to quit.", m.readErr)
	case m.content == "":
		return "Loading...\n\nPress q to quit."
	default:
		return m.content
	}
}
