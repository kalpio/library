package queries

import (
	"context"
	"library/domain"
	"library/services/author"
)

type GetAllAuthorsQueryHandler struct {
	db        domain.IDatabase
	authorSrv author.IAuthorService
}

func NewGetAllAuthorsQueryHandler(db domain.IDatabase, authorSrv author.IAuthorService) *GetAllAuthorsQueryHandler {
	return &GetAllAuthorsQueryHandler{db: db, authorSrv: authorSrv}
}

func (c *GetAllAuthorsQueryHandler) Handle(_ context.Context, _ *GetAllAuthorsQuery) (*GetAllAuthorsQueryResponse, error) {
	var (
		res []domain.Author
		err error
	)
	if res, err = c.authorSrv.GetAll(); err != nil {
		return nil, err
	}

	var results []result
	for _, r := range res {
		results = append(results, result{
			AuthorID:   r.ID,
			FirstName:  r.FirstName,
			MiddleName: r.MiddleName,
			LastName:   r.LastName,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		})
	}

	return &GetAllAuthorsQueryResponse{
		Result: results,
	}, nil
}
