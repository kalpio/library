package commands

import (
	"context"
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
	model, err := author.Edit(c.db, command.ID, command.FirstName,
		command.MiddleName, command.LastName)
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
