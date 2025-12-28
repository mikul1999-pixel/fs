package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		
		// Get tags
		sc.Tags, err = s.GetShortcutTags(sc.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get tags: %w", err)
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

func (s *SQLiteStorage) AddTags(shortcutName string, tags []string) error {
	// Get shortcut ID
	var shortcutID int
	err := s.db.QueryRow("SELECT id FROM shortcuts WHERE name = ?", shortcutName).Scan(&shortcutID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("shortcut '%s' not found", shortcutName)
	}
	if err != nil {
		return fmt.Errorf("failed to find shortcut: %w", err)
	}

	// Add each tag
	for _, tag := range tags {
		// Insert tag if it doesn't exist
		_, err := s.db.Exec("INSERT OR IGNORE INTO tags (name) VALUES (?)", tag)
		if err != nil {
			return fmt.Errorf("failed to insert tag: %w", err)
		}

		// Get tag ID
		var tagID int
		err = s.db.QueryRow("SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("failed to get tag ID: %w", err)
		}

		// Link shortcut and tag
		_, err = s.db.Exec(
			"INSERT OR IGNORE INTO shortcut_tags (shortcut_id, tag_id) VALUES (?, ?)",
			shortcutID, tagID,
		)
		if err != nil {
			return fmt.Errorf("failed to link tag: %w", err)
		}
	}

	return nil
}

func (s *SQLiteStorage) RemoveTags(shortcutName string, tags []string) error {
	// Get shortcut ID
	var shortcutID int
	err := s.db.QueryRow("SELECT id FROM shortcuts WHERE name = ?", shortcutName).Scan(&shortcutID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("shortcut '%s' not found", shortcutName)
	}
	if err != nil {
		return fmt.Errorf("failed to find shortcut: %w", err)
	}

	// Remove each tag
	for _, tag := range tags {
		_, err := s.db.Exec(`
			DELETE FROM shortcut_tags 
			WHERE shortcut_id = ? 
			AND tag_id = (SELECT id FROM tags WHERE name = ?)
		`, shortcutID, tag)
		if err != nil {
			return fmt.Errorf("failed to remove tag: %w", err)
		}
	}

	return nil
}

func (s *SQLiteStorage) GetShortcutTags(shortcutName string) ([]string, error) {
	rows, err := s.db.Query(`
		SELECT t.name 
		FROM tags t
		JOIN shortcut_tags st ON t.id = st.tag_id
		JOIN shortcuts s ON s.id = st.shortcut_id
		WHERE s.name = ?
		ORDER BY t.name
	`, shortcutName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *SQLiteStorage) SearchShortcuts(query string, tags []string) ([]Shortcut, error) {
	sqlQuery := `
		SELECT DISTINCT s.id, s.name, s.path, s.created_at, s.updated_at
		FROM shortcuts s
	`

	var conditions []string
	var args []interface{}

	// Add tag filtering
	if len(tags) > 0 {
		sqlQuery += `
			JOIN shortcut_tags st ON s.id = st.shortcut_id
			JOIN tags t ON t.id = st.tag_id
		`
		
		// Create placeholders for tags
		placeholders := make([]string, len(tags))
		for i, tag := range tags {
			placeholders[i] = "?"
			args = append(args, tag)
		}
		conditions = append(conditions, fmt.Sprintf("t.name IN (%s)", strings.Join(placeholders, ",")))
	}

	// Add text search
	if query != "" {
		conditions = append(conditions, "(s.name LIKE ? OR s.path LIKE ?)")
		searchPattern := "%" + query + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// Combine conditions
	if len(conditions) > 0 {
		sqlQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	sqlQuery += " ORDER BY s.name"

	// Execute query
	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search shortcuts: %w", err)
	}
	defer rows.Close()

	var shortcuts []Shortcut
	for rows.Next() {
		var sc Shortcut
		if err := rows.Scan(&sc.ID, &sc.Name, &sc.Path, &sc.CreatedAt, &sc.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan shortcut: %w", err)
		}
		
		// Get tags for this shortcut
		sc.Tags, err = s.GetShortcutTags(sc.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get tags for shortcut: %w", err)
		}
		
		shortcuts = append(shortcuts, sc)
	}

	return shortcuts, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}