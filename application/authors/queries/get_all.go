package queries

import (
	"library/domain"
	"time"
)

type GetAllAuthorsQuery struct {
}

type GetAllAuthorsQueryResponse struct {
	Result []result
}

type result struct {
	AuthorID   domain.AuthorID `json:"id"`
	FirstName  string          `json:"first_name"`
	MiddleName string          `json:"middle_name"`
	LastName   string          `json:"last_name"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func NewGetAllAuthorsQuery() *GetAllAuthorsQuery {
	return &GetAllAuthorsQuery{}
}
