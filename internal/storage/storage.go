package storage

type Storage interface {
	// Shortcut operations
	AddShortcut(name, path string) error
	GetShortcut(name string) (*Shortcut, error)
	ListShortcuts() ([]Shortcut, error)
	DeleteShortcut(name string) error
	
	// Tag operations
	AddTag(shortcutName string, tags []string) error
	
	// Close the database
	Close() error
}