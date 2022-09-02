package authors

import (
	"github.com/mehdihadeli/go-mediatr"
	"library/application/authors/commands"
	"library/application/authors/queries"
	"library/domain"
)

func Register(db domain.Database) error {
	var lastErr error

	createAuthorCommandHandler := commands.NewCreateAuthorCommandHandler(db)
	if err := mediatr.RegisterRequestHandler[*commands.CreateAuthorCommand, *commands.CreateAuthorCommandResponse](createAuthorCommandHandler); err != nil {
		lastErr = err
	}

	deleteAuthorCommandHandler := commands.NewDeleteAuthorCommandHandler(db)
	if err := mediatr.RegisterRequestHandler[*commands.DeleteAuthorCommand, *commands.DeleteAuthorCommandResponse](deleteAuthorCommandHandler); err != nil {
		lastErr = err
	}

	getAuthorQueryHandler := queries.NewGetAuthorQueryHandler(db)
	if err := mediatr.RegisterRequestHandler[*queries.GetAuthorQuery, *queries.GetAuthorQueryResponse](getAuthorQueryHandler); err != nil {
		lastErr = err
	}

	getAllAuthorsQueryHandler := queries.NewGetAllAuthorsQueryHandler(db)
	if err := mediatr.RegisterRequestHandler[*queries.GetAllAuthorsQuery, *queries.GetAllAuthorsQueryResponse](getAllAuthorsQueryHandler); err != nil {
		lastErr = err
	}

	return lastErr
}
