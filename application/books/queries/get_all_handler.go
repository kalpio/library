package queries

import (
	"golang.org/x/net/context"
	"library/domain"
	"library/services/book"
)

type GetAllBooksQueryHandler struct {
	db      domain.IDatabase
	bookSrv book.IBookService
}

func NewGetAllBooksQueryHandler(db domain.IDatabase, bookSrv book.IBookService) *GetAllBooksQueryHandler {
	return &GetAllBooksQueryHandler{db: db, bookSrv: bookSrv}
}

func (c *GetAllBooksQueryHandler) Handle(_ context.Context, query *GetAllBooksQuery) (*GetAllBooksQueryResponse, error) {
	var (
		res []domain.Book
		err error
	)
	if res, err = c.bookSrv.GetAll(query.Page, query.Size); err != nil {
		return nil, err
	}

	var results []result
	for _, r := range res {
		results = append(results, result{
			BookID:      r.ID,
			Title:       r.Title,
			ISBN:        r.ISBN,
			Description: r.Description,
			AuthorID:    r.AuthorID,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		})
	}

	return &GetAllBooksQueryResponse{
		Result: results,
	}, nil
}
