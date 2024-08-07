package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if the FieldErrors map is empty
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
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
