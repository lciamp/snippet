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
	// new s Snippet struct
	var s Snippet

	err := m.DB.QueryRow(`SELECT id, title, content, created, expires 
							FROM snippets 
							WHERE expires > UTC_TIMESTAMP() 
							AND id = ?`, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

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
	stmt := `SELECT id, title, content, created, expires 
					FROM snippets
					WHERE expires > UTC_TIMESTAMP()
					ORDER BY id
					DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// defer to make sure closed before return
	defer rows.Close()

	// slice of snippet structs
	var snippets []Snippet

	// rows.Nest to iterate
	for rows.Next() {
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// add snippet to snippets slice
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// all good return snippets slice
	return snippets, nil
}
