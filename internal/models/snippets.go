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
	// sql stmt to insert
	stmt := `INSERT INTO snippets (title, content, created, expires) 
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// exec the stmt to insert
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// use LastInsertId to get the id of what we inserted
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// LastInsertId returns a type int64, convert to normal int
	return int(id), nil
}

// Get function to get a snippet from db
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// Latest function that returns the 10 latest
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
