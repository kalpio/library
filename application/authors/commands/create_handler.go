package commands

import (
	"context"
	"library/application/authors/events"
	"library/domain"
	"library/services/author"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

type CreateAuthorCommandHandler struct {
	db domain.Database
}

func NewCreateAuthorCommandHandler(db domain.Database) *CreateAuthorCommandHandler {
	return &CreateAuthorCommandHandler{db: db}
}

func (c *CreateAuthorCommandHandler) Handle(ctx context.Context, command *CreateAuthorCommand) (*CreateAuthorCommandResponse, error) {
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
	}

	mediatr.Publish(ctx, &events.AuthorCreatedEvent{
		AuthorID:   model.ID,
		FirstName:  model.FirstName,
		MiddleName: model.MiddleName,
		LastName:   model.LastName,
	})

	return response, nil
}
