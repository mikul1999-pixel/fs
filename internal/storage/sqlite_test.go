package storage

import (
	"path/filepath"
	"testing"
)

func newTestSQLiteStorage(t *testing.T) *SQLiteStorage {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "shortcuts.db")
	s, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to create test storage: %v", err)
	}

	t.Cleanup(func() {
		_ = s.Close()
	})

	return s
}

func seedSearchData(t *testing.T, s *SQLiteStorage) {
	t.Helper()

	if err := s.AddShortcut("api", "/tmp/api"); err != nil {
		t.Fatalf("failed to add shortcut api: %v", err)
	}
	if err := s.AddShortcut("web", "/tmp/web"); err != nil {
		t.Fatalf("failed to add shortcut web: %v", err)
	}
	if err := s.AddShortcut("ops", "/tmp/ops"); err != nil {
		t.Fatalf("failed to add shortcut ops: %v", err)
	}

	if err := s.AddTags("api", []string{"go", "proj"}); err != nil {
		t.Fatalf("failed to add tags to api: %v", err)
	}
	if err := s.AddTags("web", []string{"proj", "frontend"}); err != nil {
		t.Fatalf("failed to add tags to web: %v", err)
	}
	if err := s.AddTags("ops", []string{"infra"}); err != nil {
		t.Fatalf("failed to add tags to ops: %v", err)
	}
}

func shortcutNames(shortcuts []Shortcut) []string {
	names := make([]string, 0, len(shortcuts))
	for _, sc := range shortcuts {
		names = append(names, sc.Name)
	}
	return names
}

func containsName(shortcuts []Shortcut, name string) bool {
	for _, sc := range shortcuts {
		if sc.Name == name {
			return true
		}
	}
	return false
}

func TestSearchShortcuts_TagOpOr(t *testing.T) {
	s := newTestSQLiteStorage(t)
	seedSearchData(t, s)

	results, err := s.SearchShortcuts("", []string{"go", "frontend"}, "or")
	if err != nil {
		t.Fatalf("SearchShortcuts returned error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results for or search, got %d (%v)", len(results), shortcutNames(results))
	}

	if !containsName(results, "api") || !containsName(results, "web") {
		t.Fatalf("expected results to contain api and web, got %v", shortcutNames(results))
	}
}

func TestSearchShortcuts_TagOpAnd(t *testing.T) {
	s := newTestSQLiteStorage(t)
	seedSearchData(t, s)

	results, err := s.SearchShortcuts("", []string{"proj", "go"}, "and")
	if err != nil {
		t.Fatalf("SearchShortcuts returned error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result for and search, got %d (%v)", len(results), shortcutNames(results))
	}

	if results[0].Name != "api" {
		t.Fatalf("expected api for and search, got %q", results[0].Name)
	}
}

func TestSearchShortcuts_TagOpAliases(t *testing.T) {
	s := newTestSQLiteStorage(t)
	seedSearchData(t, s)

	orResults, err := s.SearchShortcuts("", []string{"go", "frontend"}, "any")
	if err != nil {
		t.Fatalf("SearchShortcuts(any) returned error: %v", err)
	}

	andResults, err := s.SearchShortcuts("", []string{"proj", "go"}, "all")
	if err != nil {
		t.Fatalf("SearchShortcuts(all) returned error: %v", err)
	}

	if len(orResults) != 2 {
		t.Fatalf("expected 2 results for any alias, got %d", len(orResults))
	}
	if len(andResults) != 1 || andResults[0].Name != "api" {
		t.Fatalf("expected api for all alias, got %v", shortcutNames(andResults))
	}
}

func TestSearchShortcuts_InvalidTagOp(t *testing.T) {
	s := newTestSQLiteStorage(t)
	seedSearchData(t, s)

	_, err := s.SearchShortcuts("", []string{"proj"}, "xor")
	if err == nil {
		t.Fatal("expected error for invalid tag operator, got nil")
	}
}

func TestRemoveAllTags(t *testing.T) {
	s := newTestSQLiteStorage(t)

	if err := s.AddShortcut("cli", "/tmp/cli"); err != nil {
		t.Fatalf("failed to add shortcut: %v", err)
	}
	if err := s.AddTags("cli", []string{"go", "proj"}); err != nil {
		t.Fatalf("failed to add tags: %v", err)
	}

	if err := s.RemoveAllTags("cli"); err != nil {
		t.Fatalf("RemoveAllTags returned error: %v", err)
	}

	tags, err := s.GetShortcutTags("cli")
	if err != nil {
		t.Fatalf("GetShortcutTags returned error: %v", err)
	}

	if len(tags) != 0 {
		t.Fatalf("expected 0 tags after RemoveAllTags, got %d (%v)", len(tags), tags)
	}
}

func TestListShortcuts_LoadsTags(t *testing.T) {
	s := newTestSQLiteStorage(t)

	if err := s.AddShortcut("cli", "/tmp/cli"); err != nil {
		t.Fatalf("failed to add shortcut: %v", err)
	}
	if err := s.AddTags("cli", []string{"go", "proj"}); err != nil {
		t.Fatalf("failed to add tags: %v", err)
	}

	shortcuts, err := s.ListShortcuts()
	if err != nil {
		t.Fatalf("ListShortcuts returned error: %v", err)
	}

	if len(shortcuts) != 1 {
		t.Fatalf("expected 1 shortcut, got %d", len(shortcuts))
	}

	if len(shortcuts[0].Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d (%v)", len(shortcuts[0].Tags), shortcuts[0].Tags)
	}
}
