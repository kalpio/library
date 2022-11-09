package commands

import (
	"context"
	"library/application/books/events"
	"library/domain"
	domainEvents "library/domain/events"
	"library/services/book"

	"github.com/google/uuid"
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

	var (
		bookID   uuid.UUID
		authorID uuid.UUID
		model    *domain.Book
		err      error
	)
	bookID, err = uuid.Parse(string(command.BookID))
	if err != nil {
		return nil, err
	}

	authorID, err = uuid.Parse(string(command.AuthorID))
	model, err = c.bookSrv.Create(bookID,
		command.Title,
		command.ISBN,
		command.Description,
		authorID)

	response := &CreateBookCommandResponse{
		BookID:      domain.BookID(model.ID.String()),
		Title:       model.Title,
		ISBN:        model.ISBN,
		Description: model.Description,
		AuthorID:    domain.AuthorID(model.AuthorID.String()),
	}

	domainEvents.Publish(ctx, events.NewBookCreatedEvent(
		domain.BookID(model.ID.String()),
		model.Title,
		model.ISBN,
		model.Description,
		domain.AuthorID(model.AuthorID.String())))

	return response, nil
}
