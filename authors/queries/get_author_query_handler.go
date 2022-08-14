package queries

import (
	"context"
	"library/services/author"

	"gorm.io/gorm"
)

type GetAuthorQueryHandler struct {
	db *gorm.DB
}

func NewGetAuthorQueryHandler(db *gorm.DB) *GetAuthorQueryHandler {
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
