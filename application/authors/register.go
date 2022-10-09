package authors

import (
	"library/application/authors/commands"
	"library/application/authors/events"
	"library/application/authors/queries"
	"library/domain"
	"library/services/author"

	"github.com/mehdihadeli/go-mediatr"
)

func Register(db domain.IDatabase) error {
	var (lastErr error
		authorSrv = author.NewAuthorService(db)
	)
	createAuthorCommandHandler := commands.NewCreateAuthorCommandHandler(db, authorSrv)
	if err := mediatr.RegisterRequestHandler[
		*commands.CreateAuthorCommand,
		*commands.CreateAuthorCommandResponse](
		createAuthorCommandHandler); err != nil {
		lastErr = err
	}

	deleteAuthorCommandHandler := commands.NewDeleteAuthorCommandHandler(db, authorSrv)
	if err := mediatr.RegisterRequestHandler[
		*commands.DeleteAuthorCommand,
		*commands.DeleteAuthorCommandResponse](
		deleteAuthorCommandHandler); err != nil {
		lastErr = err
	}

	editAuthorCommandHandler := commands.NewEditAuthorCommandHandler(db, authorSrv)
	if err := mediatr.RegisterRequestHandler[
		*commands.EditAuthorCommand,
		*commands.EditAuthorCommandResponse](
		editAuthorCommandHandler); err != nil {
		lastErr = err
	}

	getAuthorByIDQueryHandler := queries.NewGetAuthorByIDQueryHandler(db, authorSrv)
	if err := mediatr.RegisterRequestHandler[
		*queries.GetAuthorByIDQuery,
		*queries.GetAuthorByIDQueryResponse](
		getAuthorByIDQueryHandler); err != nil {
		lastErr = err
	}

	getAllAuthorsQueryHandler := queries.NewGetAllAuthorsQueryHandler(db, authorSrv)
	if err := mediatr.RegisterRequestHandler[
		*queries.GetAllAuthorsQuery,
		*queries.GetAllAuthorsQueryResponse](
		getAllAuthorsQueryHandler); err != nil {
		lastErr = err
	}

	if err := mediatr.RegisterNotificationHandler[*events.AuthorCreatedEvent](
		&events.AuthorCreatedEventHandler{}); err != nil {
		lastErr = err
	}

	if err := mediatr.RegisterNotificationHandler[*events.AuthorDeletedEvent](
		&events.AuthorDeletedEventHandler{}); err != nil {
		lastErr = err
	}

	if err := mediatr.RegisterNotificationHandler[*events.AuthorEditedEvent](
		&events.AuthorEditedEventHandler{}); err != nil {
		lastErr = err
	}

	return lastErr
}
