package commands

import (
	"context"
	"github.com/google/uuid"
	"library/domain"
	"library/services/author"
)

type DeleteAuthorCommandHandler struct {
	db domain.Database
}

func NewDeleteAuthorCommandHandler(db domain.Database) *DeleteAuthorCommandHandler {
	return &DeleteAuthorCommandHandler{db: db}
}

func (c *DeleteAuthorCommandHandler) Handle(_ context.Context, command *DeleteAuthorCommand) (*DeleteAuthorCommandResponse, error) {
	authorID, err := uuid.Parse(string(command.AuthorID))
	if err != nil {
		return nil, err
	}

	succeeded, err := author.Delete(c.db, authorID)
	if err != nil {
		return nil, err
	}

	response := &DeleteAuthorCommandResponse{Succeeded: succeeded}

	return response, nil
}
