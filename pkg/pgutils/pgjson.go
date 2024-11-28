package pgutils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONB is a type representing a map with string keys and values of any type, designed to handle JSONB column data.
type JSONB map[string]any

// Value converts the JSONB object into a driver.Value, which is a byte slice representation of the JSON object.
func (a *JSONB) Value() (driver.Value, error) {
	return json.Marshal(*a)
}

// Scan implements the sql.Scanner interface for the JSONB type.
// It attempts to convert a database value to the JSONB format.
// If the value cannot be asserted as a []byte, it returns an error.
// Upon successful assertion, it unmarshals the byte array into the JSONB type.
func (a *JSONB) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
