package model_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"charm.land/bubbletea/v2"
	"github.com/yourname/mded/internal/model"
)

func newTestModel() model.Model {
	return model.New(model.ViewMode, "/tmp/test.md")
}

func TestQuitKey(t *testing.T) {
	m := newTestModel()
	_, cmd := m.Update(tea.KeyPressMsg{Text: "q", Code: 'q'})
	if cmd == nil {
		t.Fatal("expected a command for 'q' key")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg for 'q' key")
	}
}

func TestQuitCtrlC(t *testing.T) {
	m := newTestModel()
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if cmd == nil {
		t.Fatal("expected a command for ctrl+c")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg for ctrl+c")
	}
}

func TestNonQuitKey(t *testing.T) {
	m := newTestModel()
	_, cmd := m.Update(tea.KeyPressMsg{Text: "a", Code: 'a'})
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
	if !strings.Contains(m2.View().Content, want) {
		t.Fatalf("expected view to contain %q, got %q", want, m2.View().Content)
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
	if !strings.Contains(m2.View().Content, "Error") {
		t.Fatalf("expected error message in view, got %q", m2.View().Content)
	}
}
