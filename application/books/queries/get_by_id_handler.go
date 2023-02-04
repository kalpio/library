package queries

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"library/domain"
	"library/services/book"
)

type GetBookByIDQueryHandler struct {
	db      domain.IDatabase
	bookSrv book.IBookService
}

func NewGetBookByIDQueryHandler(db domain.IDatabase, bookSrv book.IBookService) *GetBookByIDQueryHandler {
	return &GetBookByIDQueryHandler{
		db:      db,
		bookSrv: bookSrv,
	}
}

func (c *GetBookByIDQueryHandler) Handle(_ context.Context, query *GetBookByIDQuery) (*GetBookByIDQueryResponse, error) {
	bookID, err := uuid.Parse(string(query.BookID))
	if err != nil {
		return nil, fmt.Errorf("book queries: cannot parse book ID: %w", err)
	}

	result, err := c.bookSrv.GetByID(bookID)
	if err != nil {
		return nil, fmt.Errorf("book queries: cannot get book from service: %w", err)
	}

	return &GetBookByIDQueryResponse{
		BookID:      result.ID,
		Title:       result.Title,
		ISBN:        result.ISBN,
		Description: result.Description,
		AuthorID:    result.AuthorID,
	}, nil
}
