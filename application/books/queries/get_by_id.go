package queries

import (
	"github.com/google/uuid"
	"library/domain"
)

type GetBookByIDQuery struct {
	BookID domain.BookID
}

type GetBookByIDQueryResponse struct {
	BookID      uuid.UUID   `json:"id"`
	Title       string      `json:"title"`
	ISBN        domain.ISBN `json:"isbn"`
	Description string      `json:"description"`
	AuthorID    uuid.UUID   `json:"author_id"`
}

func NewGetBookByIDQuery(bookID domain.BookID) *GetBookByIDQuery {
	return &GetBookByIDQuery{BookID: bookID}
}
