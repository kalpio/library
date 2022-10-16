package books

import (
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"
	"library/application/books/commands"
	"library/application/books/events"
	"library/application/books/queries"
	"library/domain"
	"library/ioc"
	"library/register"
	"library/services/book"
)

type bookRegister struct {
}

func NewBookRegister() register.IRegister[*domain.Book] {
	return &bookRegister{}
}

func (r *bookRegister) Register() error {
	database, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return errors.Wrap(err, "register [book]: failed to get database service")
	}
	bookSrv, err := ioc.Get[book.IBookService]()
	if err != nil {
		return errors.Wrap(err, "register [book]: failed to get book service")
	}

	var lastErr error

	createBookCommandHandler := commands.NewCreateBookCommandHandler(database, bookSrv)
	if err := mediatr.RegisterRequestHandler[
		*commands.CreateBookCommand,
		*commands.CreateBookCommandResponse](
		createBookCommandHandler); err != nil {
		lastErr = err
	}

	if err := mediatr.RegisterNotificationHandler[*events.BookCreatedEvent](
		&events.BookCreatedEventHandler{}); err != nil {
		lastErr = err
	}

	getBookQueryHandler := queries.NewGetBookByIDQueryHandler(database, bookSrv)
	if err := mediatr.RegisterRequestHandler[
		*queries.GetBookByIDQuery,
		*queries.GetBookByIDQueryResponse](
		getBookQueryHandler); err != nil {
		lastErr = err
	}

	if lastErr != nil {
		return errors.Wrap(lastErr, "register [book]: failed to register mediatr")
	}

	return nil
}
