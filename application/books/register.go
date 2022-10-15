package books

import (
	"github.com/mehdihadeli/go-mediatr"
	"library/application/books/commands"
	"library/application/books/events"
	"library/domain"
	"library/ioc"
	"library/services/author"
	"library/services/book"
)

func Register(db domain.IDatabase) error {
	var lastErr error
	authorSrv, err := ioc.Get[author.IAuthorService]()
	if err != nil {
		return err
	}
	bookSrv := book.NewBookService(db, authorSrv)

	createBookCommandHandler := commands.NewCreateBookCommandHandler(db, bookSrv)
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

	return lastErr
}
