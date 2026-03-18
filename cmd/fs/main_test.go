package main

import (
	"os"
	"path/filepath"
	"strings"
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

func TestRenderInitScript_IncludesErrorHandlingForJumpAndFind(t *testing.T) {
	script := renderInitScript("f", "ff")

	required := []string{
		"if [ $# -ne 1 ]; then",
		"path=\"$(fs go \"$1\")\" || return $?",
		"if [ ! -d \"$path\" ]; then",
		"path=$(fs find \"$@\" </dev/tty)",
		"local status=$?",
		"if [ $status -ne 0 ]; then",
	}

	for _, needle := range required {
		if !strings.Contains(script, needle) {
			t.Fatalf("expected init script to contain %q", needle)
		}
	}
}

func TestRenderInitScript_UsesCustomFunctionNames(t *testing.T) {
	script := renderInitScript("go", "search")

	if !strings.Contains(script, "go() {") {
		t.Fatal("expected script to define custom jump function")
	}

	if !strings.Contains(script, "Usage: go <shortcut>") {
		t.Fatal("expected script to include custom jump usage")
	}

	if !strings.Contains(script, "search() {") {
		t.Fatal("expected script to define custom find function")
	}
}
