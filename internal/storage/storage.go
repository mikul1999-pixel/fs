package storage

type Storage interface {
	// Shortcut operations
	AddShortcut(name, path string) error
	GetShortcut(name string) (*Shortcut, error)
	ListShortcuts() ([]Shortcut, error)
	DeleteShortcut(name string) error
	
	// Tag operations
	AddTags(shortcutName string, tags []string) error
	RemoveTags(shortcutName string, tags []string) error
	GetShortcutTags(shortcutName string) ([]string, error)
	SearchShortcuts(query string, tags []string) ([]Shortcut, error)
	
	// Close the database
	Close() error
}