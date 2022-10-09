package commands

import (
	"github.com/google/uuid"
	"time"
)

type CreateAuthorCommand struct {
	AuthorID   uuid.UUID
	FirstName  string
	MiddleName string
	LastName   string
}

type CreateAuthorCommandResponse struct {
	AuthorID   uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewCreateAuthorCommand(authorID uuid.UUID, firstName, middleName, lastName string) *CreateAuthorCommand {
	return &CreateAuthorCommand{
		AuthorID:   authorID,
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}
}
