package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ID interface {
	String() string
	IsNil() bool
	IsEmpty() bool
	UUID() uuid.UUID
}

type Entity[T ID] struct {
	ID        T         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *Entity[T]) BeforeCreate(db *gorm.DB) error {
	db.Set("id", e.ID.String())

	return nil
}

func (e *Entity[T]) SetID(id T) {
	e.ID = id
}

func EmptyUUID() uuid.UUID {
	return uuid.MustParse("00000000-0000-0000-0000-000000000000")
}
