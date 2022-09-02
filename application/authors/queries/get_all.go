package queries

type GetAllAuthorsQuery struct {
}

type GetAllAuthorsQueryResponse struct {
	Result []result
}

type result struct {
	AuthorID   uint
	FirstName  string
	MiddleName string
	LastName   string
}

func NewGetAllAuthorsQuery() *GetAllAuthorsQuery {
	return &GetAllAuthorsQuery{}
}
