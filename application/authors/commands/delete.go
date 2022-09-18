package commands

import (
	"library/domain"
)

type DeleteAuthorCommand struct {
	AuthorID domain.AuthorID
}

func NewDeleteAuthorCommand(authorID domain.AuthorID) *DeleteAuthorCommand {
	return &DeleteAuthorCommand{AuthorID: authorID}
}

type DeleteAuthorCommandResponse struct {
	Succeeded bool
}
