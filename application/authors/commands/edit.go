package commands

import "time"

type EditAuthorCommand struct {
	ID         uint
	FirstName  string
	MiddleName string
	LastName   string
}

type EditAuthorCommandResponse struct {
	AuthorID   uint      `json:"id"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

func NewEditAuthorCommand(id uint, firstName, middleName, lastName string) *EditAuthorCommand {
	return &EditAuthorCommand{
		ID:         id,
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}
}
