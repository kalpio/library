package commands

import (
	"context"
	"github.com/pkg/errors"
	"library/application/books/events"
	"library/domain"
	domainEvents "library/domain/events"
	"library/services/book"
)

type EditBookCommandHandler struct {
	db      domain.IDatabase
	bookSrv book.IBookService
}

func NewEditBookCommandHandler(db domain.IDatabase,
	bookSrv book.IBookService) *EditBookCommandHandler {

	return &EditBookCommandHandler{
		db:      db,
		bookSrv: bookSrv,
	}
}

func (c *EditBookCommandHandler) Handle(
	ctx context.Context,
	command *EditBookCommand) (*EditBookCommandResponse, error) {

	model, err := c.bookSrv.Edit(command.BookID,
		command.Title,
		command.ISBN,
		command.Description,
		command.AuthorID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to edit book")
	}

	response := &EditBookCommandResponse{
		BookID:      domain.BookID(model.ID.String()),
		Title:       model.Title,
		ISBN:        model.ISBN,
		Description: model.Description,
		AuthorID:    domain.AuthorID(model.AuthorID.String()),
	}

	go domainEvents.Publish(ctx, events.NewBookEditedEvent(
		model.ID,
		model.Title,
		model.ISBN,
		model.Description,
		model.AuthorID))

	return response, nil
}
