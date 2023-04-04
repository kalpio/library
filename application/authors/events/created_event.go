package events

import (
	"library/domain"
)

type AuthorCreatedEvent struct {
	AuthorID   domain.AuthorID `json:"id"`
	FirstName  string          `json:"first_name"`
	MiddleName string          `json:"middle_name"`
	LastName   string          `json:"last_name"`
}
