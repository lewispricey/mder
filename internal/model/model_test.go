package model_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lewispricey/mded/internal/model"
)

func newTestModel() model.Model {
	return model.New(model.ViewMode, "/tmp/test.md")
}

func TestQuitKey(t *testing.T) {
	m := newTestModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Fatal("expected a command for 'q' key")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg for 'q' key")
	}
}

func TestQuitCtrlC(t *testing.T) {
	m := newTestModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected a command for ctrl+c")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg for ctrl+c")
	}
}

func TestNonQuitKey(t *testing.T) {
	m := newTestModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if cmd != nil {
		t.Fatal("expected no command for non-quit key")
	}
}

func TestFileReadSuccess(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.md")
	want := "# Hello\n\nWorld"
	os.WriteFile(f, []byte(want), 0644)

	m := model.New(model.ViewMode, f)
	cmd := m.Init()
	msg := cmd()
	m2, _ := m.Update(msg)
	if !strings.Contains(m2.View(), want) {
		t.Fatalf("expected view to contain %q, got %q", want, m2.View())
	}
}

func TestWindowResize(t *testing.T) {
	m := newTestModel()
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	got := m2.(model.Model)
	if got.Width() != 120 {
		t.Fatalf("expected width 120, got %d", got.Width())
	}
	if got.Height() != 40 {
		t.Fatalf("expected height 40, got %d", got.Height())
	}
}

func TestFileReadError(t *testing.T) {
	m := model.New(model.ViewMode, "/nonexistent/path.md")
	cmd := m.Init()
	msg := cmd()
	m2, cmd2 := m.Update(msg)
	if cmd2 == nil {
		t.Fatal("expected quit command on read error")
	}
	if _, ok := cmd2().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg on read error")
	}
	if !strings.Contains(m2.View(), "Error") {
		t.Fatalf("expected error message in view, got %q", m2.View())
	}
}

// Edit mode tests

func loadInEditMode(t *testing.T, content string) (model.Model, tea.Cmd) {
	t.Helper()
	dir := t.TempDir()
	f := filepath.Join(dir, "test.md")
	os.WriteFile(f, []byte(content), 0644)

	m := model.New(model.EditMode, f)
	cmd := m.Init()
	msg := cmd()
	m2, cmd2 := m.Update(msg)
	return m2.(model.Model), cmd2
}

func TestEditModeTextareaInit(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	if !strings.Contains(m.View(), "hello") {
		t.Fatalf("expected view to contain initial content, got %q", m.View())
	}
}

func TestEditModeTyping(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	got := m2.(model.Model)
	if !strings.Contains(got.TextareaValue(), "hellox") {
		t.Fatalf("expected textarea value to contain typed char, got %q", got.TextareaValue())
	}
}

func TestEditModeCtrlCQuits(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected a command for ctrl+c in edit mode")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg for ctrl+c in edit mode")
	}
}

func TestEditModeQTypes(t *testing.T) {
	m, _ := loadInEditMode(t, "hel")
	m2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	got := m2.(model.Model)
	if !strings.Contains(got.TextareaValue(), "helq") {
		t.Fatalf("expected 'q' to be typed, got %q", got.TextareaValue())
	}
	if cmd != nil {
		if _, ok := cmd().(tea.QuitMsg); ok {
			t.Fatal("did not expect QuitMsg when typing in edit mode")
		}
	}
}
