package queries

type GetAuthorByIDQuery struct {
	AuthorID uint
}

type GetAuthorByIDQueryResponse struct {
	AuthorID   uint   `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

func NewGetAuthorByIDQuery(authorID uint) *GetAuthorByIDQuery {
	return &GetAuthorByIDQuery{AuthorID: authorID}
}
