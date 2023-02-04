package authors

import (
	"github.com/pkg/errors"
	"library/application/authors/commands"
	"library/application/authors/events"
	"library/application/authors/queries"
	"library/domain"
	"library/ioc"
	"library/register"
	"library/services/author"

	"github.com/mehdihadeli/go-mediatr"
)

type authorRegister struct {
}

func NewAuthorRegister() register.IRegister[*domain.Author] {
	return &authorRegister{}
}

func (r *authorRegister) Register() error {
	database, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return errors.Wrap(err, "register [author]: failed to get database service")
	}

	authorSrv, err := ioc.Get[author.IAuthorService]()
	if err != nil {
		return errors.Wrap(err, "register [author]: failed to get author service")
	}

	var lastErr error

	createAuthorCommandHandler := commands.NewCreateAuthorCommandHandler(database, authorSrv)
	if err := mediatr.RegisterRequestHandler[
		*commands.CreateAuthorCommand,
		*commands.CreateAuthorCommandResponse](
		createAuthorCommandHandler); err != nil {
		lastErr = err
	}

	deleteAuthorCommandHandler := commands.NewDeleteAuthorCommandHandler(database, authorSrv)
	if err := mediatr.RegisterRequestHandler[
		*commands.DeleteAuthorCommand,
		*commands.DeleteAuthorCommandResponse](
		deleteAuthorCommandHandler); err != nil {
		lastErr = err
	}

	editAuthorCommandHandler := commands.NewEditAuthorCommandHandler(database, authorSrv)
	if err := mediatr.RegisterRequestHandler[
		*commands.EditAuthorCommand,
		*commands.EditAuthorCommandResponse](
		editAuthorCommandHandler); err != nil {
		lastErr = err
	}

	getAuthorByIDQueryHandler := queries.NewGetAuthorByIDQueryHandler(database, authorSrv)
	if err := mediatr.RegisterRequestHandler[
		*queries.GetAuthorByIDQuery,
		*queries.GetAuthorByIDQueryResponse](
		getAuthorByIDQueryHandler); err != nil {
		lastErr = err
	}

	getAllAuthorsQueryHandler := queries.NewGetAllAuthorsQueryHandler(database, authorSrv)
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
