package validator

import "fmt"

var ErrValidator *ErrorValidator

type ErrorValidator struct {
	field string
	tag   string
	param string
	value interface{}
}

func (m *ErrorValidator) Error() string {
	if m.param != "" {
		return fmt.Sprintf("value: '%v' in the '%s' field does not meet the requirement: %s[%s]", m.value, m.field, m.tag, m.param)
	}
	return fmt.Sprintf("value: '%v' in the '%s' field does not meet the requirement: %s", m.value, m.field, m.tag)
}
