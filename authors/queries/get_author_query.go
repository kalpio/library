package queries

type GetAuthorQuery struct {
	AuthorID uint
}

type GetAuthorQueryResponse struct {
	AuthorID   uint
	FirstName  string
	MiddleName string
	LastName   string
}

func NewGetAuthorQuery(authorID uint) *GetAuthorQuery {
	return &GetAuthorQuery{AuthorID: authorID}
}
