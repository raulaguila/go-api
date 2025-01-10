package validator

import (
	"fmt"
)

type ValidateError struct {
	field string
	tag   string
	param string
	value any
}

func (m *ValidateError) Error() string {
	return fmt.Sprintf(`{"value":"%v","field":"%v","tag":"%v","param":"%v"}`, m.value, m.field, m.tag, m.param)
}
