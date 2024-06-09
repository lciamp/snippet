package models

import (
	"database/sql"
	"time"
)

// Snippet create snippet type
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel type that warps db connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert function to insert into db
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// Get function to get a snippet from db
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// Latest function that returns the 10 latest
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
