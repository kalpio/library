package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Entity struct {
	ID        uuid.UUID      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (e *Entity) BeforeCreate(db *gorm.DB) error {
	db.Set("id", e.ID.String())

	return nil
}

func EmptyUUID() uuid.UUID {
	return uuid.MustParse("00000000-0000-0000-0000-000000000000")
}
