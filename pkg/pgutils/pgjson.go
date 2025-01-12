package pgutils

import (
	"encoding/json"
	"errors"
	"fmt"

	"database/sql/driver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JSONB json.RawMessage

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSONB(result)

	return err
}

func (j *JSONB) MarshalJSON() ([]byte, error) {
	v, err := j.Value()
	return v.([]byte), err
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	result := json.RawMessage{}
	err := json.Unmarshal(data, &result)
	*j = JSONB(result)

	return err
}

func (j *JSONB) Value() (driver.Value, error) {
	if len(*j) == 0 {
		return nil, nil
	}
	return json.RawMessage(*j).MarshalJSON()
}

func (j *JSONB) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}
