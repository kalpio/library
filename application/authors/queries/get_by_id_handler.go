package queries

import (
	"context"
	"library/domain"
	"library/services/author"
)

type GetAuthorByIDQueryHandler struct {
	db        domain.IDatabase
	authorSrv author.IAuthorService
}

func NewGetAuthorByIDQueryHandler(db domain.IDatabase, authorSrv author.IAuthorService) *GetAuthorByIDQueryHandler {
	return &GetAuthorByIDQueryHandler{db: db, authorSrv: authorSrv}
}

func (c *GetAuthorByIDQueryHandler) Handle(_ context.Context, query *GetAuthorByIDQuery) (*GetAuthorByIDQueryResponse, error) {
	result, err := c.authorSrv.GetByID(query.AuthorID.UUID())
	if err != nil {
		return nil, err
	}

	return &GetAuthorByIDQueryResponse{
		AuthorID:   result.ID,
		FirstName:  result.FirstName,
		MiddleName: result.MiddleName,
		LastName:   result.LastName,
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}, nil
}
