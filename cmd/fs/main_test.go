package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath_EmptyInput(t *testing.T) {
	_, err := expandPath("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestExpandPath_TildeExpansion(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home dir: %v", err)
	}

	got, err := expandPath("~/my-project")
	if err != nil {
		t.Fatalf("expandPath returned error: %v", err)
	}

	want := filepath.Join(home, "my-project")
	if got != want {
		t.Fatalf("unexpected expanded path. got=%q want=%q", got, want)
	}
}

func TestExpandPath_RelativeToAbsolute(t *testing.T) {
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldCwd)
	}()

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}

	got, err := expandPath("./repo")
	if err != nil {
		t.Fatalf("expandPath returned error: %v", err)
	}

	want := filepath.Join(tmp, "repo")
	if got != want {
		t.Fatalf("unexpected absolute path. got=%q want=%q", got, want)
	}
}
