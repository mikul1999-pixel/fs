package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db *sql.DB
}

// Create a new SQLite storage instance
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	// Create config directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &SQLiteStorage{db: db}

	// Initialize tables
	if err := storage.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return storage, nil
}

func (s *SQLiteStorage) initTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS shortcuts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		path TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS shortcut_tags (
		shortcut_id INTEGER,
		tag_id INTEGER,
		FOREIGN KEY (shortcut_id) REFERENCES shortcuts(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
		PRIMARY KEY (shortcut_id, tag_id)
	);
	`

	_, err := s.db.Exec(schema)
	return err
}

func (s *SQLiteStorage) AddShortcut(name, path string) error {
	_, err := s.db.Exec(
		"INSERT INTO shortcuts (name, path) VALUES (?, ?)",
		name, path,
	)
	if err != nil {
		return fmt.Errorf("failed to add shortcut: %w", err)
	}
	return nil
}

func (s *SQLiteStorage) GetShortcut(name string) (*Shortcut, error) {
	var sc Shortcut
	err := s.db.QueryRow(
		"SELECT id, name, path, created_at, updated_at FROM shortcuts WHERE name = ?",
		name,
	).Scan(&sc.ID, &sc.Name, &sc.Path, &sc.CreatedAt, &sc.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shortcut '%s' not found", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get shortcut: %w", err)
	}

	return &sc, nil
}

func (s *SQLiteStorage) ListShortcuts() ([]Shortcut, error) {
	rows, err := s.db.Query(
		"SELECT id, name, path, created_at, updated_at FROM shortcuts ORDER BY name",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list shortcuts: %w", err)
	}
	defer rows.Close()

	var shortcuts []Shortcut
	for rows.Next() {
		var sc Shortcut
		if err := rows.Scan(&sc.ID, &sc.Name, &sc.Path, &sc.CreatedAt, &sc.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shortcut: %w", err)
		}
		shortcuts = append(shortcuts, sc)
	}

	return shortcuts, nil
}

func (s *SQLiteStorage) DeleteShortcut(name string) error {
	result, err := s.db.Exec("DELETE FROM shortcuts WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("failed to delete shortcut: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("shortcut '%s' not found", name)
	}

	return nil
}

func (s *SQLiteStorage) AddTag(shortcutName string, tags []string) error {
	// We'll implement this later
	return fmt.Errorf("not implemented yet")
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}