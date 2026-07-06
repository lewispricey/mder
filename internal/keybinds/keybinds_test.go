package keybinds_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lewispricey/mded/internal/keybinds"
)

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

func TestIsQuitNonQuit(t *testing.T) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	if keybinds.IsQuit(msg) {
		t.Fatal("expected 'a' not to be a quit key")
	}
}
