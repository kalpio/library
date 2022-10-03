package queries

import (
	"context"
	"library/domain"
	"library/services/author"
)

type GetAllAuthorsQueryHandler struct {
	db domain.Database
}

func NewGetAllAuthorsQueryHandler(db domain.Database) *GetAllAuthorsQueryHandler {
	return &GetAllAuthorsQueryHandler{db: db}
}

func (c *GetAllAuthorsQueryHandler) Handle(_ context.Context, _ *GetAllAuthorsQuery) (*GetAllAuthorsQueryResponse, error) {
	var (
		res []domain.Author
		err error
	)
	if res, err = author.GetAll(c.db); err != nil {
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
