package commands

import (
	"github.com/google/uuid"
	"library/domain"
	"time"
)

type EditAuthorCommand struct {
	ID         domain.AuthorID
	FirstName  string
	MiddleName string
	LastName   string
}

type EditAuthorCommandResponse struct {
	AuthorID   uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

func NewEditAuthorCommand(id domain.AuthorID, firstName, middleName, lastName string) *EditAuthorCommand {
	return &EditAuthorCommand{
		ID:         id,
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}
}
