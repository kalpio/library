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
	database domain.IDatabase
	bookSrv  book.IBookService
}

func NewBookRegister() register.IRegister[*domain.Book] {
	return &bookRegister{}
}

func (r *bookRegister) Register() error {
	var lastErr error
	if err := r.resolveServices(); err != nil {
		lastErr = err
	}

	if err := r.registerCreateCommand(); err != nil {
		lastErr = err
	}

	if err := r.registerEditCommand(); err != nil {
		lastErr = err
	}

	if err := r.registerDeleteCommand(); err != nil {
		lastErr = err
	}

	if err := r.registerGetByIDQuery(); err != nil {
		lastErr = err
	}

	if err := r.registerGetAllQuery(); err != nil {
		lastErr = err
	}

	return errors.Wrap(lastErr, "register [book]: failed to register mediatr")
}

func (r *bookRegister) resolveServices() error {
	database, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return errors.Wrap(err, "register [book]: failed to get database service")
	}
	r.database = database

	bookSrv, err := ioc.Get[book.IBookService]()
	if err != nil {
		return errors.Wrap(err, "register [book]: failed to get book service")
	}
	r.bookSrv = bookSrv

	return nil
}

func (r *bookRegister) registerCreateCommand() error {
	var lastErr error

	createBookCommandHandler := commands.NewCreateBookCommandHandler(r.database, r.bookSrv)
	if err := mediatr.RegisterRequestHandler[
		*commands.CreateBookCommand,
		*commands.CreateBookCommandResponse](
		createBookCommandHandler); err != nil {
		lastErr = errors.Wrap(err, "register [book]: failed to register create book command")
	}

	if err := mediatr.RegisterNotificationHandler[*events.BookCreatedEvent](
		&events.BookCreatedEventHandler{}); err != nil {
		lastErr = errors.Wrap(err, "register [book]: failed to register book created event")
	}

	return lastErr
}

func (r *bookRegister) registerEditCommand() error {
	var lastErr error

	editBookCommandHandler := commands.NewEditBookCommandHandler(r.database, r.bookSrv)
	if err := mediatr.RegisterRequestHandler[
		*commands.EditBookCommand,
		*commands.EditBookCommandResponse](
		editBookCommandHandler); err != nil {
		lastErr = errors.Wrap(err, "register [book]: failed to register edit book command")
	}

	if err := mediatr.RegisterNotificationHandler[*events.BookEditedEvent](
		&events.BookEditedEventHandler{}); err != nil {
		lastErr = errors.Wrap(err, "register [book]: failed to register book edited event")
	}

	return lastErr
}

func (r *bookRegister) registerDeleteCommand() error {
	var lastErr error

	deleteBookCommandHandler := commands.NewDeleteBookCommandHandler(r.database, r.bookSrv)
	if err := mediatr.RegisterRequestHandler[
		*commands.DeleteBookCommand,
		*commands.DeleteBookCommandResponse](
		deleteBookCommandHandler); err != nil {
		lastErr = errors.Wrap(err, "register [book]: failed to register delete book command")
	}

	if err := mediatr.RegisterNotificationHandler[*events.BookDeletedEvent](
		&events.BookDeletedEventHandler{}); err != nil {
		lastErr = errors.Wrap(err, "register [book]: failed to register book deleted event")
	}

	return lastErr
}

func (r *bookRegister) registerGetByIDQuery() error {
	getBookQueryHandler := queries.NewGetBookByIDQueryHandler(r.database, r.bookSrv)
	if err := mediatr.RegisterRequestHandler[
		*queries.GetBookByIDQuery,
		*queries.GetBookByIDQueryResponse](
		getBookQueryHandler); err != nil {
		return errors.Wrap(err, "register [book]: failed to register get book by id query")
	}

	return nil
}

func (r *bookRegister) registerGetAllQuery() error {
	getAllBooksQueryHandler := queries.NewGetAllBooksQueryHandler(r.database, r.bookSrv)
	if err := mediatr.RegisterRequestHandler[
		*queries.GetAllBooksQuery,
		*queries.GetAllBooksQueryResponse](
		getAllBooksQueryHandler); err != nil {
		return errors.Wrap(err, "register [book]: failed to register get all books query")
	}

	return nil
}
