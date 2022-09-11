package queries

import (
	"context"
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
	result, err := author.GetByID(c.db, query.AuthorID)
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
