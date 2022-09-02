package queries

import (
	"context"
	"library/domain"
	"library/services/author"
)

type GetAuthorQueryHandler struct {
	db domain.Database
}

func NewGetAuthorQueryHandler(db domain.Database) *GetAuthorQueryHandler {
	return &GetAuthorQueryHandler{db: db}
}

func (c *GetAuthorQueryHandler) Handle(ctx context.Context, query *GetAuthorQuery) (*GetAuthorQueryResponse, error) {
	result, err := author.GetByID(c.db, query.AuthorID)
	if err != nil {
		return nil, err
	}

	return &GetAuthorQueryResponse{
		AuthorID:   result.ID,
		FirstName:  result.FirstName,
		MiddleName: result.MiddleName,
		LastName:   result.LastName,
	}, nil
}
