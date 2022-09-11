package commands

import (
	"context"
	"library/domain"
	"library/services/author"
)

type DeleteAuthorCommandHandler struct {
	db domain.Database
}

func NewDeleteAuthorCommandHandler(db domain.Database) *DeleteAuthorCommandHandler {
	return &DeleteAuthorCommandHandler{db: db}
}

func (c *DeleteAuthorCommandHandler) Handle(ctx context.Context, command *DeleteAuthorCommand) (*DeleteAuthorCommandResponse, error) {
	succeeded, err := author.Delete(c.db, command.AuthorID)
	if err != nil {
		return nil, err
	}

	response := &DeleteAuthorCommandResponse{Succeeded: succeeded}

	return response, nil
}
