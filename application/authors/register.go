package authors

import (
	"library/application/authors/commands"
	"library/application/authors/queries"
	"library/domain"

	"github.com/mehdihadeli/go-mediatr"
)

func Register(db domain.Database) error {
	var lastErr error

	createAuthorCommandHandler := commands.NewCreateAuthorCommandHandler(db)
	if err := mediatr.RegisterRequestHandler[
		*commands.CreateAuthorCommand,
		*commands.CreateAuthorCommandResponse](
		createAuthorCommandHandler); err != nil {
		lastErr = err
	}

	deleteAuthorCommandHandler := commands.NewDeleteAuthorCommandHandler(db)
	if err := mediatr.RegisterRequestHandler[
		*commands.DeleteAuthorCommand,
		*commands.DeleteAuthorCommandResponse](
		deleteAuthorCommandHandler); err != nil {
		lastErr = err
	}

	editAuthorCommandHandler := commands.NewEditAuthorCommandHandler(db)
	if err := mediatr.RegisterRequestHandler[
		*commands.EditAuthorCommand,
		*commands.EditAuthorCommandResponse](
		editAuthorCommandHandler); err != nil {
		lastErr = err
	}

	deletePermanentlyCommandHandler := commands.NewDeletePermanentlyCommandHandler(db)
	if err := mediatr.RegisterRequestHandler[
		*commands.DeletePermanentlyCommand,
		*commands.DeletePermanentlyCommandResponse](
		deletePermanentlyCommandHandler); err != nil {
		lastErr = err
	}

	getAuthorByIDQueryHandler := queries.NewGetAuthorByIDQueryHandler(db)
	if err := mediatr.RegisterRequestHandler[
		*queries.GetAuthorByIDQuery,
		*queries.GetAuthorByIDQueryResponse](
		getAuthorByIDQueryHandler); err != nil {
		lastErr = err
	}

	getAllAuthorsQueryHandler := queries.NewGetAllAuthorsQueryHandler(db)
	if err := mediatr.RegisterRequestHandler[
		*queries.GetAllAuthorsQuery,
		*queries.GetAllAuthorsQueryResponse](
		getAllAuthorsQueryHandler); err != nil {
		lastErr = err
	}

	return lastErr
}
