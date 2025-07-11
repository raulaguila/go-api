package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	BaseInt struct {
		ID        uint      `gorm:"primarykey"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	BaseUUID struct {
		ID        uuid.UUID `gorm:"primarykey;type:uuid;default:uuid_generate_v4()"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}
)
