package events

import "library/domain"

type BookEditedEvent struct {
	BookID      domain.BookID   `json:"id"`
	Title       string          `json:"title"`
	ISBN        string          `json:"isbn"`
	Description string          `json:"description"`
	AuthorID    domain.AuthorID `json:"author_id"`
}

func NewBookEditedEvent(id domain.BookID,
	title, isbn, description string,
	authorID domain.AuthorID) *BookEditedEvent {

	return &BookEditedEvent{
		BookID:      id,
		Title:       title,
		ISBN:        isbn,
		Description: description,
		AuthorID:    authorID,
	}
}
