package commands

import (
	"context"
	"library/domain"
	"library/services/author"
)

type CreateAuthorCommandHandler struct {
	db domain.Database
}

func NewCreateAuthorCommandHandler(db domain.Database) *CreateAuthorCommandHandler {
	return &CreateAuthorCommandHandler{db: db}
}

func (c *CreateAuthorCommandHandler) Handle(ctx context.Context, command *CreateAuthorCommand) (*CreateAuthorCommandResponse, error) {
	model, err := author.Create(c.db, command.FirstName, command.MiddleName, command.LastName)
	if err != nil {
		return nil, err
	}

	response := &CreateAuthorCommandResponse{AuthorID: model.ID}

	return response, nil
}
