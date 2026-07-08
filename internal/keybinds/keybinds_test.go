package keybinds_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lewispricey/mded/internal/keybinds"
)

func TestIsHardQuitCtrlC(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	if !keybinds.IsHardQuit(msg) {
		t.Fatal("expected ctrl+c to be a hard quit")
	}
}

func TestIsHardQuitQ(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	if keybinds.IsHardQuit(msg) {
		t.Fatal("expected 'q' not to be a hard quit")
	}
}

func TestIsQuitQ(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	if !keybinds.IsQuit(msg) {
		t.Fatal("expected 'q' to be a quit key")
	}
}

func TestIsQuitCtrlC(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	if !keybinds.IsQuit(msg) {
		t.Fatal("expected ctrl+c to be a quit key")
	}
}

func TestIsSaveCtrlS(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyCtrlS}
	if !keybinds.IsSave(msg) {
		t.Fatal("expected ctrl+s to be a save key")
	}
}

func TestIsSaveNonSave(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	if keybinds.IsSave(msg) {
		t.Fatal("expected 's' not to be a save key")
	}
}

func TestIsQuitNonQuit(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	if keybinds.IsQuit(msg) {
		t.Fatal("expected 'a' not to be a quit key")
	}
}

func TestIsToggleModeCtrlE(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyCtrlE}
	if !keybinds.IsToggleMode(msg) {
		t.Fatal("expected ctrl+e to be a toggle key")
	}
}

func TestIsToggleModeNonToggle(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	if keybinds.IsToggleMode(msg) {
		t.Fatal("expected 'e' not to be a toggle key")
	}
}
