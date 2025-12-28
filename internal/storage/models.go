package storage

import "time"

type Shortcut struct {
	ID        int
	Name      string
	Path      string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}