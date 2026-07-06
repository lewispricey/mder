package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourname/mded/internal/cmd"
)

func runCommand(t *testing.T, args ...string) error {
	t.Helper()
	root := cmd.NewRootCmd()
	root.SetArgs(args)
	return root.Execute()
}

func TestViewMode(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.md")
	os.WriteFile(f, []byte("# hi"), 0644)

	err := runCommand(t, f, "--view")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestEditMode(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.md")
	os.WriteFile(f, []byte("# hi"), 0644)

	err := runCommand(t, f, "--edit")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDefaultModeIsEdit(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.md")
	os.WriteFile(f, []byte("# hi"), 0644)

	err := runCommand(t, f)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestMissingFile(t *testing.T) {
	err := runCommand(t, "/nonexistent/path.md", "--view")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestConflictingFlags(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.md")
	os.WriteFile(f, []byte("# hi"), 0644)

	err := runCommand(t, f, "--view", "--edit")
	if err == nil {
		t.Fatal("expected error for conflicting flags, got nil")
	}
}

func TestNoArgs(t *testing.T) {
	err := runCommand(t)
	if err == nil {
		t.Fatal("expected error for no args, got nil")
	}
}
