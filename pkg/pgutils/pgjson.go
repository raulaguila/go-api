package pgutils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONB Interface for JSONB Field of yourTableName Table
type JSONB map[string]any

// Value Marshal
func (a *JSONB) Value() (driver.Value, error) {
	return json.Marshal(*a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
