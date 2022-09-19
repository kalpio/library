package commands

import (
	"context"
	"github.com/google/uuid"
	"library/domain"
	"library/services/author"
)

type CreateAuthorCommandHandler struct {
	db domain.Database
}

func NewCreateAuthorCommandHandler(db domain.Database) *CreateAuthorCommandHandler {
	return &CreateAuthorCommandHandler{db: db}
}

func (c *CreateAuthorCommandHandler) Handle(_ context.Context, command *CreateAuthorCommand) (*CreateAuthorCommandResponse, error) {
	model, err := author.Create(c.db, uuid.New(), command.FirstName, command.MiddleName, command.LastName)
	if err != nil {
		return nil, err
	}

	response := &CreateAuthorCommandResponse{
		AuthorID:   model.ID,
		FirstName:  model.FirstName,
		MiddleName: model.MiddleName,
		LastName:   model.LastName,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
		DeletedAt:  model.DeletedAt,
	}

	return response, nil
}
