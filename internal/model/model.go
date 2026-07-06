package model

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lewispricey/mded/internal/keybinds"
)

var paneStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

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
	mode         Mode
	filePath     string
	content      string
	readErr      error
	status       string
	textarea     textarea.Model
	width        int
	height       int
	dirty        bool
	cleanContent string
	quitting     bool
}

func (m Model) Width() int            { return m.width }
func (m Model) Height() int           { return m.height }
func (m Model) FilePath() string      { return m.filePath }
func (m Model) TextareaValue() string { return m.textarea.Value() }
func (m Model) PaneWidths() (int, int) {
	if m.width < 4 {
		return 0, 0
	}
	gap := 1
	leftWidth := (m.width - gap) / 2
	rightWidth := m.width - leftWidth - gap
	return leftWidth, rightWidth
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
			leftWidth, _ := m.PaneWidths()
			if leftWidth > 0 {
				m.textarea.SetWidth(leftWidth - 2)
			}
			if m.height >= 4 {
				m.textarea.SetHeight(m.height - 3)
			}
			m.textarea.SetValue(msg.content)
			m.cleanContent = msg.content
			return m, m.textarea.Focus()
		}
		return m, nil
	case saveMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Save error: %v", msg.err)
		} else {
			m.status = "Saved"
			m.cleanContent = m.textarea.Value()
			m.dirty = false
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		leftWidth, _ := m.PaneWidths()
		if leftWidth > 0 && m.mode == EditMode && m.content != "" {
			m.textarea.SetWidth(leftWidth - 2)
		}
		if m.mode == EditMode && m.content != "" && m.height >= 4 {
			m.textarea.SetHeight(m.height - 3)
		}
		return m, nil
	case tea.KeyMsg:
		m.status = ""
		if keybinds.IsHardQuit(msg) {
			if m.mode == EditMode && m.dirty {
				if m.quitting {
					return m, tea.Quit
				}
				m.quitting = true
				return m, nil
			}
			return m, tea.Quit
		}
		m.quitting = false
		if keybinds.IsSave(msg) && m.mode == EditMode {
			return m, saveFile(m.filePath, m.textarea.Value())
		}
		if m.mode == EditMode {
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			m.dirty = m.textarea.Value() != m.cleanContent
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
		leftWidth, rightWidth := m.PaneWidths()

		leftPane := paneStyle.Width(leftWidth).Render(m.textarea.View())
		rightPane := paneStyle.Width(rightWidth).Render("— Preview —")

		layout := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

		var footer string
		if m.quitting {
			footer = "Unsaved changes! Press ctrl+c again to quit."
		} else if m.status != "" {
			footer = m.status
		} else {
			footer = quitHint
		}
		return layout + "\n" + footer
	default:
		return m.content
	}
}
