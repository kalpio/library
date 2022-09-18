package queries

import (
	"github.com/google/uuid"
	"library/domain"
)

type GetAuthorByIDQuery struct {
	AuthorID domain.AuthorID
}

type GetAuthorByIDQueryResponse struct {
	AuthorID   uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
}

func NewGetAuthorByIDQuery(authorID domain.AuthorID) *GetAuthorByIDQuery {
	return &GetAuthorByIDQuery{AuthorID: authorID}
}
