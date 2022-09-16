package queries

import "time"

type GetAllAuthorsQuery struct {
}

type GetAllAuthorsQueryResponse struct {
	Result []result
}

type result struct {
	AuthorID   uint      `json:"id"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewGetAllAuthorsQuery() *GetAllAuthorsQuery {
	return &GetAllAuthorsQuery{}
}
