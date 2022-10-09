package commands

import (
	"context"
	"library/application/authors/events"
	"library/domain"
	domain_events "library/domain/events"
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

	succeeded, err := c.authorSrv.Delete(authorID)
	if err != nil {
		return nil, err
	}

	response := &DeleteAuthorCommandResponse{Succeeded: succeeded}
	if succeeded {
		domain_events.Publish(ctx, &events.AuthorDeletedEvent{AuthorID: authorID})
	}

	return response, nil
}
