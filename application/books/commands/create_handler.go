package commands

import (
	"context"
	"github.com/pkg/errors"
	"library/application/books/events"
	"library/domain"
	domainEvents "library/domain/events"
	"library/services/book"
)

type CreateBookCommandHandler struct {
	db      domain.IDatabase
	bookSrv book.IBookService
}

func NewCreateBookCommandHandler(db domain.IDatabase,
	bookSrv book.IBookService) *CreateBookCommandHandler {

	return &CreateBookCommandHandler{
		db:      db,
		bookSrv: bookSrv,
	}
}

func (c *CreateBookCommandHandler) Handle(
	ctx context.Context,
	command *CreateBookCommand) (*CreateBookCommandResponse, error) {

	model, err := c.bookSrv.Create(command.BookID,
		command.Title,
		command.ISBN,
		command.Description,
		command.AuthorID)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create book")
	}

	response := &CreateBookCommandResponse{
		BookID:      domain.BookID(model.ID.String()),
		Title:       model.Title,
		ISBN:        model.ISBN,
		Description: model.Description,
		AuthorID:    domain.AuthorID(model.AuthorID.String()),
	}

	go domainEvents.Publish(ctx, events.NewBookCreatedEvent(
		model.ID,
		model.Title,
		model.ISBN,
		model.Description,
		model.AuthorID))

	return response, nil
}
