package commands

import (
	"library/domain"
	"time"
)

type CreateAuthorCommand struct {
	AuthorID   domain.AuthorID
	FirstName  string
	MiddleName string
	LastName   string
}

type CreateAuthorCommandResponse struct {
	AuthorID   domain.AuthorID `json:"id"`
	FirstName  string          `json:"first_name"`
	MiddleName string          `json:"middle_name"`
	LastName   string          `json:"last_name"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func NewCreateAuthorCommand(authorID domain.AuthorID, firstName, middleName, lastName string) *CreateAuthorCommand {
	return &CreateAuthorCommand{
		AuthorID:   authorID,
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}
}
