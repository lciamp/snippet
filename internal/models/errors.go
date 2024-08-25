package models

import (
	"errors"
)

var (
	// ErrNoRecord no matching record
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials login with incorrect email or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrDuplicateEmail signup with an email that's already in use
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
