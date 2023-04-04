package events

import (
	"library/domain"
)

type AuthorDeletedEvent struct {
	AuthorID domain.AuthorID `json:"id"`
}
