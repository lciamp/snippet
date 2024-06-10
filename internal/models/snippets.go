package models

import (
	"database/sql"
	"errors"
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
	// stmt to get snippet
	stmt := `SELECT id, title, content, created, expires FROM snippets 
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// returns pointer to row object
	row := m.DB.QueryRow(stmt, id)

	// new s Snippet struct
	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// if no rows Scan will return ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	// all good, return the snippet
	return s, nil
}

// Latest function that returns the 10 latest
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
