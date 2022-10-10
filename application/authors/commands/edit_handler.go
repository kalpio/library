package commands

import (
	"context"
	"library/application/authors/events"
	"library/domain"
	"library/services/author"

	domain_events "library/domain/events"

	"github.com/google/uuid"
)

type EditAuthorCommandHandler struct {
	db        domain.IDatabase
	authorSrv author.IAuthorService
}

func NewEditAuthorCommandHandler(db domain.IDatabase, authorSrv author.IAuthorService) *EditAuthorCommandHandler {
	return &EditAuthorCommandHandler{db: db, authorSrv: authorSrv}
}

func (c *EditAuthorCommandHandler) Handle(ctx context.Context, command *EditAuthorCommand) (*EditAuthorCommandResponse, error) {
	authorID, err := uuid.Parse(string(command.ID))
	if err != nil {
		return nil, err
	}

	model, err := c.authorSrv.Edit(authorID, command.FirstName, command.MiddleName, command.LastName)
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
	}

	domain_events.Publish(ctx, &events.AuthorEditedEvent{
		AuthorID:   model.ID,
		FirstName:  model.FirstName,
		MiddleName: model.MiddleName,
		LastName:   model.LastName,
	})

	return response, nil
}
