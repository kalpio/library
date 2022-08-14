package commands

import (
	"library/services/author"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type CreateAuthorCommandHandler struct {
	db *gorm.DB
}

func NewCreateAuthorCommandHandler(db *gorm.DB) *CreateAuthorCommandHandler {
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
