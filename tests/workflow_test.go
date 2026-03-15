package tests

import (
	"path/filepath"
	"testing"

	"github.com/mikul1999-pixel/fs/internal/storage"
)

func TestAliasTagFallbackWorkflow(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "shortcuts.db")
	s, err := storage.NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })

	if err := s.AddShortcut("cli", "/tmp/cli"); err != nil {
		t.Fatalf("failed to add cli: %v", err)
	}
	if err := s.AddShortcut("api", "/tmp/api"); err != nil {
		t.Fatalf("failed to add api: %v", err)
	}
	if err := s.AddTags("cli", []string{"proj", "go"}); err != nil {
		t.Fatalf("failed to tag cli: %v", err)
	}
	if err := s.AddTags("api", []string{"proj", "backend"}); err != nil {
		t.Fatalf("failed to tag api: %v", err)
	}

	results, err := s.SearchShortcuts("", []string{"proj", "go"}, "and")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(results) != 1 || results[0].Name != "cli" {
		t.Fatalf("expected only cli from fallback search, got %v", results)
	}
}
