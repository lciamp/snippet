package models

import (
	"database/sql"
	"time"
)

// User struct types align with db columns
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// UserModel model struct to connect to db
type UserModel struct {
	DB *sql.DB
}

// Insert method
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate method
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists method
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
