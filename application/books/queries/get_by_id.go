package queries

import (
	"library/domain"
)

type GetBookByIDQuery struct {
	BookID domain.BookID
}

type GetBookByIDQueryResponse struct {
	BookID      domain.BookID   `json:"id"`
	Title       string          `json:"title"`
	ISBN        string          `json:"isbn"`
	Description string          `json:"description"`
	AuthorID    domain.AuthorID `json:"author_id"`
}

func NewGetBookByIDQuery(bookID domain.BookID) *GetBookByIDQuery {
	return &GetBookByIDQuery{BookID: bookID}
}
