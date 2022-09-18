package queries

import (
	"context"
	"github.com/google/uuid"
	"library/domain"
	"library/services/author"
)

type GetAuthorByIDQueryHandler struct {
	db domain.Database
}

func NewGetAuthorByIDQueryHandler(db domain.Database) *GetAuthorByIDQueryHandler {
	return &GetAuthorByIDQueryHandler{db: db}
}

func (c *GetAuthorByIDQueryHandler) Handle(_ context.Context, query *GetAuthorByIDQuery) (*GetAuthorByIDQueryResponse, error) {
	authorID, err := uuid.Parse(string(query.AuthorID))
	if err != nil {
		return nil, err
	}
	result, err := author.GetByID(c.db, authorID)
	if err != nil {
		return nil, err
	}

	return &GetAuthorByIDQueryResponse{
		AuthorID:   result.ID,
		FirstName:  result.FirstName,
		MiddleName: result.MiddleName,
		LastName:   result.LastName,
	}, nil
}
