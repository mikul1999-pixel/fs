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
	m := InitialModel(testShortcuts(), SelectorOptions{NoColor: true})

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
	m := InitialModel(testShortcuts(), SelectorOptions{NoColor: true})

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m1 := updated.(model)

	if !m1.quitting {
		t.Fatal("expected model to be quitting after q")
	}
}

func TestModelView_ShowsShortcutAndTags(t *testing.T) {
	m := InitialModel(testShortcuts(), SelectorOptions{NoColor: true})
	view := m.View()

	if !strings.Contains(view, "api -> /tmp/api") {
		t.Fatalf("expected view to contain api shortcut, got %q", view)
	}
	if !strings.Contains(view, "[go, proj]") {
		t.Fatalf("expected view to contain tags, got %q", view)
	}
}

func TestModelView_HighlightsQueryAndMatchedTagWithColor(t *testing.T) {
	m := InitialModel(testShortcuts(), SelectorOptions{
		Query:      "ap",
		FilterTags: []string{"go"},
	})

	view := m.View()

	if !strings.Contains(view, "api") {
		t.Fatalf("expected view to contain shortcut name, got %q", view)
	}
	if !strings.Contains(view, "go") || !strings.Contains(view, "proj") {
		t.Fatalf("expected view to contain rendered tags, got %q", view)
	}
}

func TestQueryTokens_DedupesAndSortsByLength(t *testing.T) {
	tokens := queryTokens("api proj api p")

	if len(tokens) != 3 {
		t.Fatalf("expected 3 unique tokens, got %d (%v)", len(tokens), tokens)
	}

	if tokens[0] != "proj" {
		t.Fatalf("expected longest token first, got %v", tokens)
	}
}
