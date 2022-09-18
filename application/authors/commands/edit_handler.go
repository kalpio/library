package commands

import (
	"context"
	"github.com/google/uuid"
	"library/domain"
	"library/services/author"
)

type EditAuthorCommandHandler struct {
	db domain.Database
}

func NewEditAuthorCommandHandler(db domain.Database) *EditAuthorCommandHandler {
	return &EditAuthorCommandHandler{db: db}
}

func (c *EditAuthorCommandHandler) Handle(_ context.Context, command *EditAuthorCommand) (*EditAuthorCommandResponse, error) {
	authorID, err := uuid.Parse(string(command.ID))
	if err != nil {
		return nil, err
	}

	model, err := author.Edit(c.db, authorID, command.FirstName, command.MiddleName, command.LastName)
	if err != nil {
		return nil, err
	}

	response := &EditAuthorCommandResponse{
		AuthorID:   model.ID,
		FirstName:  model.FirstName,
		MiddleName: model.MiddleName,
		LastName:   model.LastName,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
		DeletedAt:  model.DeletedAt.Time.UTC(),
	}

	return response, nil
}
