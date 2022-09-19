package commands

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type CreateAuthorCommand struct {
	FirstName  string
	MiddleName string
	LastName   string
}

type CreateAuthorCommandResponse struct {
	AuthorID   uuid.UUID      `json:"id"`
	FirstName  string         `json:"first_name"`
	MiddleName string         `json:"middle_name"`
	LastName   string         `json:"last_name"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

func NewCreateAuthorCommand(firstName, middleName, lastName string) *CreateAuthorCommand {
	return &CreateAuthorCommand{
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}
}
