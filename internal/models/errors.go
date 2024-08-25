package models

import (
	"errors"
)

var (
	// ErrNoRecord no matching record
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials login with incorrect email or password
	ErrInvalidCredentials = err.New("models: invalid credentials")

	// ErrDuplicateEmail signup with an email that's already in use
	ErrDuplicateEmail = err.New("models: duplicate email")
)
