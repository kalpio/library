package commands

import (
	"context"
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

	// TODO: don't swallow the error
	model, _ := c.bookSrv.Create(command.BookID,
		command.Title,
		command.ISBN,
		command.Description,
		command.AuthorID)

	response := &CreateBookCommandResponse{
		BookID:      domain.BookID(model.ID.String()),
		Title:       model.Title,
		ISBN:        model.ISBN,
		Description: model.Description,
		AuthorID:    domain.AuthorID(model.AuthorID.String()),
	}

	domainEvents.Publish(ctx, events.NewBookCreatedEvent(
		model.ID,
		model.Title,
		model.ISBN,
		model.Description,
		model.AuthorID))

	return response, nil
}
