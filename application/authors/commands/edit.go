package commands

import (
	"library/domain"
	"time"
)

type EditAuthorCommand struct {
	AuthorID   domain.AuthorID
	FirstName  string
	MiddleName string
	LastName   string
}

type EditAuthorCommandResponse struct {
	AuthorID   domain.AuthorID `json:"id"`
	FirstName  string          `json:"first_name"`
	MiddleName string          `json:"middle_name"`
	LastName   string          `json:"last_name"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func NewEditAuthorCommand(id domain.AuthorID, firstName, middleName, lastName string) *EditAuthorCommand {
	return &EditAuthorCommand{
		AuthorID:   id,
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}
}
