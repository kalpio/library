package events

import "library/domain"

type BookCreatedEvent struct {
	BookID      domain.BookID   `json:"id"`
	Title       string          `json:"title"`
	ISBN        domain.ISBN     `json:"isbn"`
	Description string          `json:"description"`
	AuthorID    domain.AuthorID `json:"author_id"`
}

func NewBookCreatedEvent(bookID domain.BookID,
	title string,
	isbn domain.ISBN,
	description string,
	authorID domain.AuthorID) *BookCreatedEvent {

	return &BookCreatedEvent{
		BookID:      bookID,
		Title:       title,
		ISBN:        isbn,
		Description: description,
		AuthorID:    authorID,
	}
}
