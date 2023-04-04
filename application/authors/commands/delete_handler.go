package commands

import (
	"context"
	"library/application/authors/events"
	"library/domain"
	domainEvents "library/domain/events"
	"library/services/author"
)

type DeleteAuthorCommandHandler struct {
	db        domain.IDatabase
	authorSrv author.IAuthorService
}

func NewDeleteAuthorCommandHandler(db domain.IDatabase, authorSrv author.IAuthorService) *DeleteAuthorCommandHandler {
	return &DeleteAuthorCommandHandler{db: db, authorSrv: authorSrv}
}

func (c *DeleteAuthorCommandHandler) Handle(ctx context.Context, command *DeleteAuthorCommand) (*DeleteAuthorCommandResponse, error) {
	err := c.authorSrv.Delete(command.AuthorID)
	if err != nil {
		return nil, err
	}

	response := &DeleteAuthorCommandResponse{Succeeded: true}
	domainEvents.Publish(ctx, &events.AuthorDeletedEvent{AuthorID: command.AuthorID})

	return response, nil
}
