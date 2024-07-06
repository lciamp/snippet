package validator

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
