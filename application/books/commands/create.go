package commands

import (
	"library/domain"
)

type CreateBookCommand struct {
	BookID      domain.BookID
	Title       string
	ISBN        domain.ISBN
	Description string
	AuthorID    domain.AuthorID
}

type CreateBookCommandResponse struct {
	BookID      domain.BookID   `json:"id"`
	Title       string          `json:"title"`
	ISBN        domain.ISBN     `json:"isbn"`
	Description string          `json:"description"`
	AuthorID    domain.AuthorID `json:"author_id"`
}

func NewCreateBookCommand(id domain.BookID,
	title string,
	isbn domain.ISBN,
	description string,
	authorID domain.AuthorID) *CreateBookCommand {

	return &CreateBookCommand{
		BookID:      id,
		Title:       title,
		ISBN:        isbn,
		Description: description,
		AuthorID:    authorID,
	}
}
