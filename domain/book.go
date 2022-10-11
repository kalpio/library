package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type BookID string

type Book struct {
	Entity
	Title       string
	ISBN        string `gorm:"uniqueIndex;size:13"`
	Description string
	AuthorID    uuid.UUID
	Author      *Author
}

func NewBook(id uuid.UUID, title, isbn, description string, author *Author) *Book {
	return &Book{
		Entity:      Entity{ID: id},
		Title:       title,
		ISBN:        isbn,
		Description: description,
		AuthorID:    author.ID,
		Author:      author,
	}
}

func (b Book) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.ID, validation.Required),
		validation.Field(&b.Title, validation.Required),
		validation.Field(&b.ISBN, validation.Required),
		validation.Field(&b.Author, validation.Required))
}

func (b Book) GetID() uuid.UUID {
	return b.ID
}
