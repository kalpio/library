package commands

import (
	"golang.org/x/net/context"
	"library/application/books/events"
	"library/domain"
	domainEvents "library/domain/events"
	"library/services/book"
)

type DeleteBookCommandHandler struct {
	db      domain.IDatabase
	bookSrv book.IBookService
}

func NewDeleteBookCommandHandler(db domain.IDatabase, bookSrv book.IBookService) *DeleteBookCommandHandler {
	return &DeleteBookCommandHandler{db: db, bookSrv: bookSrv}
}

func (c *DeleteBookCommandHandler) Handle(ctx context.Context, command *DeleteBookCommand) (*DeleteBookCommandResponse, error) {
	err := c.bookSrv.Delete(command.BookID)
	if err != nil {
		return nil, err
	}

	response := &DeleteBookCommandResponse{Succeeded: true}
	go domainEvents.Publish(ctx, &events.BookDeletedEvent{BookID: command.BookID})

	return response, nil
}
