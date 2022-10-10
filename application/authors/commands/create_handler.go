package commands

import (
	"context"
	"library/application/authors/events"
	"library/domain"
	domain_events "library/domain/events"
	"library/services/author"
)

type CreateAuthorCommandHandler struct {
	db        domain.IDatabase
	authorSrv author.IAuthorService
}

func NewCreateAuthorCommandHandler(db domain.IDatabase, authorSrv author.IAuthorService) *CreateAuthorCommandHandler {
	return &CreateAuthorCommandHandler{db: db, authorSrv: authorSrv}
}

func (c *CreateAuthorCommandHandler) Handle(ctx context.Context, command *CreateAuthorCommand) (*CreateAuthorCommandResponse, error) {
	model, err := c.authorSrv.Create(command.AuthorID, command.FirstName, command.MiddleName, command.LastName)
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
	}

	domain_events.Publish(ctx, &events.AuthorCreatedEvent{
		AuthorID:   model.ID,
		FirstName:  model.FirstName,
		MiddleName: model.MiddleName,
		LastName:   model.LastName,
	})

	return response, nil
}
