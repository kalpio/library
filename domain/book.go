package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type BookID string

type Book struct {
	Entity
	Title       string    `gorm:"column:title" json:"title"`
	ISBN        string    `gorm:"uniqueIndex;size:13" json:"isbn"`
	Description string    `json:"description"`
	AuthorID    uuid.UUID `json:"author_id"`
	Author      *Author   `json:"-"`
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
		validation.Field(&b.ID, validation.By(b.validateID)),
		validation.Field(&b.Title, validation.Required),
		validation.Field(&b.ISBN, validation.Required),
		validation.Field(&b.ISBN, validation.Length(13, 13)),
		validation.Field(&b.AuthorID, validation.By(b.validateAuthorID)),
	)
}

func (b Book) validateAuthorID(_ interface{}) error {
	if b.AuthorID == uuid.Nil {
		return validation.NewError("author_id", "author_id is null")
	}
	if b.AuthorID == EmptyUUID() {
		return validation.NewError("author_id", "author_id id is empty")
	}
	if b.Author != nil && b.AuthorID != b.Author.ID {
		return validation.NewError("author_id", "author_id is not equal to author.id")
	}

	return nil
}

func (b Book) validateID(_ interface{}) error {
	if b.ID == uuid.Nil {
		return validation.NewError("id", "id is null")
	}
	if b.ID == EmptyUUID() {
		return validation.NewError("id", "id is empty")
	}

	return nil
}

func (b Book) GetID() uuid.UUID {
	return b.ID
}
