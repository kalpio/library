package events

import "github.com/google/uuid"

type AuthorDeletedEvent struct {
	AuthorID uuid.UUID `json:"id"`
}
