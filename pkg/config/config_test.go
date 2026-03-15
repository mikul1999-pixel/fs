package config

import (
	"path/filepath"
	"testing"
)

func TestGetDBPath_UsesXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/fs-config")

	got := GetDBPath()
	want := filepath.Join("/tmp/fs-config", "fs", "shortcuts.db")

	if got != want {
		t.Fatalf("unexpected db path. got=%q want=%q", got, want)
	}
}

func TestGetDBPath_FallsBackToHomeConfig(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "/tmp/fs-home")

	got := GetDBPath()
	want := filepath.Join("/tmp/fs-home", ".config", "fs", "shortcuts.db")

	if got != want {
		t.Fatalf("unexpected fallback db path. got=%q want=%q", got, want)
	}
}
