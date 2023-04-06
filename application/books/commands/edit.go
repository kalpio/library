package commands

import "library/domain"

type EditBookCommand struct {
	BookID      domain.BookID
	Title       string
	ISBN        string
	Description string
	AuthorID    domain.AuthorID
}

func NewEditBookCommand(id domain.BookID,
	title, isbn, description string,
	authorID domain.AuthorID) *EditBookCommand {

	return &EditBookCommand{
		BookID:      id,
		Title:       title,
		ISBN:        isbn,
		Description: description,
		AuthorID:    authorID,
	}
}

type EditBookCommandResponse struct {
	BookID      domain.BookID   `json:"id"`
	Title       string          `json:"title"`
	ISBN        string          `json:"isbn"`
	Description string          `json:"description"`
	AuthorID    domain.AuthorID `json:"author_id"`
}
