package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

// EmailRx use regex.MustCompile() function to check format of email addresses
// returns a pointer to a 'compiled' regex type or panics
var EmailRx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid returns true if the FieldErrors map is empty
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddNonFieldErrors adds to the NonFieldErrors slice
func (v *Validator) AddNonFieldErrors(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// AddFieldErrors adds and error msg to the FieldErrors map (if it doesn't already exist)
func (v *Validator) AddFieldErrors(key, message string) {

	// initialize map if it doesn't already exist
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	// if key doesn't exist in map, create it
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField adds an error to the FieldErrors map only if it fails validation
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldErrors(key, message)
	}
}

// NotBlank returns true if field is not a blank string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if the value is less than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue return true if the value is in a list of permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// MinChars returns true if the value is less than n characters
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches returns true if a value matches a compiled regular expression
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
