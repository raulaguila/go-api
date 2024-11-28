package validator

import (
	"fmt"
)

// ErrValidator is a variable used to reference an instance of ValidateError, indicating validation issues within a structure.
var ErrValidator *ValidateError

// ValidateError represents an error during validation, storing details about the failed field and condition.
type ValidateError struct {
	field string
	tag   string
	param string
	value any
}

// Error constructs and returns a string representation of the validation error encapsulated by ValidateError.
func (m *ValidateError) Error() string {
	if m.param != "" {
		return fmt.Sprintf("value: '%v' in the '%s' field does not meet the requirement: %s[%s]", m.value, m.field, m.tag, m.param)
	}
	return fmt.Sprintf("value: '%v' in the '%s' field does not meet the requirement: %s", m.value, m.field, m.tag)
}
