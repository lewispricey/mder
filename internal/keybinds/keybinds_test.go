package keybinds_test

import (
	"testing"

	"charm.land/bubbletea/v2"
	"github.com/yourname/mded/internal/keybinds"
)

func TestIsQuitQ(t *testing.T) {
	msg := tea.KeyPressMsg{Text: "q", Code: 'q'}
	if !keybinds.IsQuit(msg) {
		t.Fatal("expected 'q' to be a quit key")
	}
}

func TestIsQuitCtrlC(t *testing.T) {
	msg := tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
	if !keybinds.IsQuit(msg) {
		t.Fatal("expected ctrl+c to be a quit key")
	}
}

func TestIsQuitNonQuit(t *testing.T) {
	msg := tea.KeyPressMsg{Text: "a", Code: 'a'}
	if keybinds.IsQuit(msg) {
		t.Fatal("expected 'a' not to be a quit key")
	}
}
