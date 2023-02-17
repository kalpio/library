package commands

import (
	"context"
	"library/application/authors/events"
	"library/domain"
	domainEvents "library/domain/events"
	"library/services/author"

	"github.com/google/uuid"
)

type DeleteAuthorCommandHandler struct {
	db        domain.IDatabase
	authorSrv author.IAuthorService
}

func NewDeleteAuthorCommandHandler(db domain.IDatabase, authorSrv author.IAuthorService) *DeleteAuthorCommandHandler {
	return &DeleteAuthorCommandHandler{db: db, authorSrv: authorSrv}
}

func (c *DeleteAuthorCommandHandler) Handle(ctx context.Context, command *DeleteAuthorCommand) (*DeleteAuthorCommandResponse, error) {
	authorID, err := uuid.Parse(string(command.AuthorID))
	if err != nil {
		return nil, err
	}

	err = c.authorSrv.Delete(authorID)
	if err != nil {
		return &DeleteAuthorCommandResponse{Succeeded: false}, err
	}

	response := &DeleteAuthorCommandResponse{Succeeded: true}
	domainEvents.Publish(ctx, &events.AuthorDeletedEvent{AuthorID: authorID})

	return response, nil
}
