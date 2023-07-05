package queries

import (
	"library/domain"
	"time"
)

type GetAllBooksQuery struct {
	Page int
	Size int
}

func NewGetAllBooksQuery() *GetAllBooksQuery {
	return &GetAllBooksQuery{}
}

type GetAllBooksQueryResponse struct {
	Result []result
}

type result struct {
	BookID      domain.BookID   `json:"id"`
	Title       string          `json:"title"`
	ISBN        string          `json:"isbn"`
	Description string          `json:"description"`
	AuthorID    domain.AuthorID `json:"author_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
