package model

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
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

type saveMsg struct {
	err error
}

type Model struct {
	mode     Mode
	filePath string
	content  string
	readErr  error
	status   string
	textarea textarea.Model
	width    int
	height   int
}

func (m Model) Width() int            { return m.width }
func (m Model) Height() int           { return m.height }
func (m Model) FilePath() string      { return m.filePath }
func (m Model) TextareaValue() string { return m.textarea.Value() }

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

func saveFile(path, content string) tea.Cmd {
	return func() tea.Msg {
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			err = fmt.Errorf("save %s: %w", path, err)
		}
		return saveMsg{err: err}
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
		if m.mode == EditMode {
			m.textarea = textarea.New()
			m.textarea.SetValue(msg.content)
			return m, m.textarea.Focus()
		}
		return m, nil
	case saveMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Save error: %v", msg.err)
		} else {
			m.status = "Saved"
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		m.status = ""
		if keybinds.IsHardQuit(msg) {
			return m, tea.Quit
		}
		if keybinds.IsSave(msg) && m.mode == EditMode {
			return m, saveFile(m.filePath, m.textarea.Value())
		}
		if m.mode == EditMode {
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
		if keybinds.IsQuit(msg) {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	var quitHint string
	if m.mode == EditMode {
		quitHint = "Press ctrl+c to quit."
	} else {
		quitHint = "Press q to quit."
	}

	switch {
	case m.readErr != nil:
		return fmt.Sprintf("Error: %v\n\n%s", m.readErr, quitHint)
	case m.content == "":
		return "Loading...\n\n" + quitHint
	case m.mode == EditMode:
		v := m.textarea.View()
		if m.status != "" {
			v += "\n" + m.status
		}
		return v
	default:
		return m.content
	}
}
