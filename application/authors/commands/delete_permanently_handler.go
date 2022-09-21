package commands

import (
	"context"
	"github.com/google/uuid"
	"library/domain"
	"library/services/author"
)

type DeletePermanentlyCommandHandler struct {
	db domain.Database
}

func NewDeletePermanentlyCommandHandler(db domain.Database) *DeletePermanentlyCommandHandler {
	return &DeletePermanentlyCommandHandler{db: db}
}

func (c *DeletePermanentlyCommandHandler) Handle(_ context.Context, command *DeletePermanentlyCommand) (*DeletePermanentlyCommandResponse, error) {
	authorID, err := uuid.Parse(string(command.AuthorID))
	if err != nil {
		return nil, err
	}

	succeeded, err := author.DeletePermanently(c.db, authorID)
	if err != nil {
		return nil, err
	}

	response := &DeletePermanentlyCommandResponse{Succeeded: succeeded}

	return response, nil
}
