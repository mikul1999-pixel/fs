package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikul1999-pixel/fs/internal/storage"
)

func testShortcuts() []storage.Shortcut {
	return []storage.Shortcut{
		{Name: "api", Path: "/tmp/api", Tags: []string{"go", "proj"}},
		{Name: "web", Path: "/tmp/web", Tags: []string{"frontend"}},
	}
}

func TestModelUpdate_NavigationAndSelection(t *testing.T) {
	m := InitialModel(testShortcuts())

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m1 := updated.(model)
	if m1.cursor != 1 {
		t.Fatalf("expected cursor at index 1, got %d", m1.cursor)
	}

	updated, _ = m1.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := updated.(model)
	if m2.selected == nil {
		t.Fatal("expected selected shortcut after enter")
	}
	if m2.selected.Name != "web" {
		t.Fatalf("expected selected shortcut 'web', got %q", m2.selected.Name)
	}
}

func TestModelUpdate_Quit(t *testing.T) {
	m := InitialModel(testShortcuts())

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m1 := updated.(model)

	if !m1.quitting {
		t.Fatal("expected model to be quitting after q")
	}
}

func TestModelView_ShowsShortcutAndTags(t *testing.T) {
	m := InitialModel(testShortcuts())
	view := m.View()

	if !strings.Contains(view, "api -> /tmp/api") {
		t.Fatalf("expected view to contain api shortcut, got %q", view)
	}
	if !strings.Contains(view, "[go, proj]") {
		t.Fatalf("expected view to contain tags, got %q", view)
	}
}
