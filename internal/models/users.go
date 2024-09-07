package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	// create bcrypt has for plain txt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, password, created) 
	VALUES (?,?,?, UTC_TIMESTAMP())`

	// try to exec to insert
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
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
