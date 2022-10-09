package queries

import (
	"context"
	"github.com/google/uuid"
	"library/domain"
	"library/services/author"
)

type GetAuthorByIDQueryHandler struct {
	db domain.IDatabase
	authorSrv author.IAuthorService
}

func NewGetAuthorByIDQueryHandler(db domain.IDatabase, authorSrv author.IAuthorService) *GetAuthorByIDQueryHandler {
	return &GetAuthorByIDQueryHandler{db: db, authorSrv: authorSrv}
}

func (c *GetAuthorByIDQueryHandler) Handle(_ context.Context, query *GetAuthorByIDQuery) (*GetAuthorByIDQueryResponse, error) {
	authorID, err := uuid.Parse(string(query.AuthorID))
	if err != nil {
		return nil, err
	}
	result, err := c.authorSrv.GetByID(authorID)
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
