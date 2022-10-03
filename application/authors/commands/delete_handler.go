package commands

import (
	"context"
	"library/application/authors/events"
	"library/domain"
	"library/services/author"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

type DeleteAuthorCommandHandler struct {
	db domain.Database
}

func NewDeleteAuthorCommandHandler(db domain.Database) *DeleteAuthorCommandHandler {
	return &DeleteAuthorCommandHandler{db: db}
}

func (c *DeleteAuthorCommandHandler) Handle(ctx context.Context, command *DeleteAuthorCommand) (*DeleteAuthorCommandResponse, error) {
	authorID, err := uuid.Parse(string(command.AuthorID))
	if err != nil {
		return nil, err
	}

	succeeded, err := author.Delete(c.db, authorID)
	if err != nil {
		return nil, err
	}

	response := &DeleteAuthorCommandResponse{Succeeded: succeeded}
	if succeeded {
		mediatr.Publish(ctx, &events.AuthorDeletedEvent{AuthorID: authorID})
	}

	return response, nil
}
