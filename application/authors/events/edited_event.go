package events

import "github.com/google/uuid"

type AuthorEditedEvent struct {
	AuthorID   uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
}
