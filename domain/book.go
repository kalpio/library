package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type BookID string

func (b BookID) String() string {
	return string(b)
}

func (b BookID) IsNil() bool {
	return b == BookID(uuid.Nil.String())
}

func (b BookID) IsEmpty() bool {
	return b == BookID(EmptyUUID().String())
}

func (b BookID) UUID() uuid.UUID {
	return uuid.MustParse(b.String())
}

func NewBookID() BookID {
	return BookID(uuid.NewString())
}

type Book struct {
	Entity[BookID]
	Title       string   `gorm:"column:title" json:"title"`
	ISBN        string   `gorm:"uniqueIndex;size:13" json:"isbn"`
	Description string   `json:"description"`
	AuthorID    AuthorID `json:"author_id"`
	Author      *Author  `json:"-"`
}

func NewBook(id BookID, title, isbn, description string, author *Author) *Book {
	return &Book{
		Entity:      Entity[BookID]{ID: id},
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

func (b Book) GetID() uuid.UUID {
	return b.ID.UUID()
}

func (b Book) validateAuthorID(_ interface{}) error {
	if b.AuthorID.IsNil() {
		return validation.NewError("author_id", "author_id is null")
	}
	if b.AuthorID.IsEmpty() {
		return validation.NewError("author_id", "author_id id is empty")
	}
	if b.Author != nil && b.AuthorID != b.Author.ID {
		return validation.NewError("author_id", "author_id is not equal to author.id")
	}

	return nil
}

func (b Book) validateID(_ interface{}) error {
	if b.ID.IsNil() {
		return validation.NewError("id", "id is null")
	}
	if b.ID.IsEmpty() {
		return validation.NewError("id", "id is empty")
	}

	return nil
}
