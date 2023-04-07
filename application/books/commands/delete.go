package commands

import "library/domain"

type DeleteBookCommand struct {
	BookID domain.BookID
}

func NewDeleteBookCommand(bookID domain.BookID) *DeleteBookCommand {
	return &DeleteBookCommand{BookID: bookID}
}

type DeleteBookCommandResponse struct {
	Succeeded bool
}
