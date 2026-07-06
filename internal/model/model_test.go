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

func TestSaveSuccess(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	got := m2.(model.Model)

	_, saveCmd := got.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	if saveCmd == nil {
		t.Fatal("expected save command for ctrl+s")
	}
	msg := saveCmd()
	m3, _ := got.Update(msg)
	got2 := m3.(model.Model)

	data, err := os.ReadFile(got2.FilePath())
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello!" {
		t.Fatalf("expected file content 'hello!', got %q", string(data))
	}
	if !strings.Contains(got2.View(), "Saved") {
		t.Fatalf("expected 'Saved' status after save, got %q", got2.View())
	}
}

func TestSaveError(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "readonly.md")
	os.WriteFile(f, []byte("hello"), 0444)

	m := model.New(model.EditMode, f)
	cmd := m.Init()
	msg := cmd()
	m2, _ := m.Update(msg)
	got := m2.(model.Model)

	_, saveCmd := got.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	if saveCmd == nil {
		t.Fatal("expected save command")
	}
	saveMsgVal := saveCmd()
	m3, _ := got.Update(saveMsgVal)
	got2 := m3.(model.Model)

	if !strings.Contains(got2.View(), "Save error") {
		t.Fatalf("expected save error in view, got %q", got2.View())
	}

	data, err := os.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("expected original content preserved, got %q", string(data))
	}
}

func TestStatusClearedOnKeystroke(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	_, saveCmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	msg := saveCmd()
	m2, _ := m.Update(msg)
	got := m2.(model.Model)

	if !strings.Contains(got.View(), "Saved") {
		t.Fatal("expected 'Saved' status")
	}
	m3, _ := got.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	got2 := m3.(model.Model)
	if strings.Contains(got2.View(), "Saved") {
		t.Fatal("status should clear on keystroke")
	}
}

func TestCtrlSIgnoredInViewMode(t *testing.T) {
	m := newTestModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	if cmd != nil {
		t.Fatal("expected no command for ctrl+s in view mode")
	}
}

func TestUnsavedChangedBlocksQuit(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	got := m2.(model.Model)

	m3, cmd := got.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	got2 := m3.(model.Model)
	if cmd != nil {
		t.Fatal("expected no quit command on first ctrl+c with unsaved changes")
	}
	if !strings.Contains(got2.View(), "Unsaved changes") {
		t.Fatalf("expected quitting prompt in view, got %q", got2.View())
	}

	m4, cmd2 := got2.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	_ = m4
	if cmd2 == nil {
		t.Fatal("expected a command on second ctrl+c")
	}
	if _, ok := cmd2().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg on second ctrl+c")
	}
}

func TestCleanFileExitsImmediately(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected a command for ctrl+c on clean file")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg for ctrl+c on clean file")
	}
}

func TestKeypressResetsQuitting(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	got := m2.(model.Model)

	m3, _ := got.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	got2 := m3.(model.Model)

	m4, _ := got2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	got3 := m4.(model.Model)

	_, cmd := got3.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd != nil {
		t.Fatal("expected first ctrl+c still blocked after keypress resets quitting")
	}
}

func TestSaveClearsDirty(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	got := m2.(model.Model)

	_, saveCmd := got.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	msg := saveCmd()
	m3, _ := got.Update(msg)
	got2 := m3.(model.Model)

	_, cmd := got2.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected quit command after save")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg after save")
	}
}

func TestUndoRestoreClearsDirty(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	got := m2.(model.Model)
	m3, _ := got.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	got2 := m3.(model.Model)

	_, cmd := got2.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected quit command when content restored to original")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg when content restored to original")
	}
}

func TestQuittingSaveClearsQuitting(t *testing.T) {
	m, _ := loadInEditMode(t, "hello")
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	got := m2.(model.Model)

	m3, _ := got.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	got2 := m3.(model.Model)

	_, saveCmd := got2.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	if saveCmd == nil {
		t.Fatal("expected save command")
	}
	msg := saveCmd()
	m4, _ := got2.Update(msg)
	got3 := m4.(model.Model)

	_, cmd := got3.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected quit command after save during quitting state")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg after save during quitting state")
	}
}

func TestSplitLayoutPaneWidths(t *testing.T) {
	for _, width := range []int{80, 120, 200} {
		m, _ := loadInEditMode(t, "hello")
		m2, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 40})
		got := m2.(model.Model)
		left, right := got.PaneWidths()
		diff := left - right
		if diff < -1 || diff > 1 {
			t.Errorf("width %d: pane widths differ by %d (left=%d, right=%d)",
				width, diff, left, right)
		}
	}
}

func TestSplitLayoutBorders(t *testing.T) {
	m, _ := loadInEditMode(t, "# Hello\n\nWorld")
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	got := m2.(model.Model)
	v := got.View()
	if !strings.Contains(v, "╭") {
		t.Fatal("expected top-left border corner in split layout view")
	}
	if !strings.Contains(v, "— Preview —") {
		t.Fatal("expected placeholder text in split layout view")
	}
}

func TestSplitLayoutResize(t *testing.T) {
	m, _ := loadInEditMode(t, "# Hello\n\nWorld")
	m2, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 40})
	got := m2.(model.Model)
	v80 := got.View()

	m3, _ := got.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	got2 := m3.(model.Model)
	v120 := got2.View()

	if v80 == v120 {
		t.Fatal("expected view to differ after resize from 80 to 120")
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
