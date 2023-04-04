package commands

import (
	"context"
	"library/application/authors/events"
	"library/domain"
	"library/services/author"

	domainEvents "library/domain/events"
)

type EditAuthorCommandHandler struct {
	db        domain.IDatabase
	authorSrv author.IAuthorService
}

func NewEditAuthorCommandHandler(db domain.IDatabase, authorSrv author.IAuthorService) *EditAuthorCommandHandler {
	return &EditAuthorCommandHandler{db: db, authorSrv: authorSrv}
}

func (c *EditAuthorCommandHandler) Handle(ctx context.Context, command *EditAuthorCommand) (*EditAuthorCommandResponse, error) {
	model, err := c.authorSrv.Edit(command.AuthorID.UUID(), command.FirstName, command.MiddleName, command.LastName)
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

	domainEvents.Publish(ctx, &events.AuthorEditedEvent{
		AuthorID:   model.ID,
		FirstName:  model.FirstName,
		MiddleName: model.MiddleName,
		LastName:   model.LastName,
	})

	return response, nil
}
