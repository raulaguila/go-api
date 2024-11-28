package domain

import (
	"time"
)

// Base provides common fields for other structs, including ID, CreatedAt, and UpdatedAt timestamps.
type Base struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
