package commands

import (
	"library/domain"
)

type DeletePermanentlyCommand struct {
	AuthorID domain.AuthorID
}

func NewDeletePermanentlyCommand(authorID domain.AuthorID) *DeletePermanentlyCommand {
	return &DeletePermanentlyCommand{AuthorID: authorID}
}

type DeletePermanentlyCommandResponse struct {
	Succeeded bool
}
