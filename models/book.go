package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title    string
	ISBN     string `gorm:"uniqueIndex;size:13"`
	Content  []byte
	Format   string
	Version  string
	AuthorID uint
	Author   *Author
}

func NewBook(title, isbn, format string, author *Author) *Book {
	return &Book{
		Title:    title,
		ISBN:     isbn,
		Content:  []byte{},
		Format:   format,
		Version:  "",
		AuthorID: author.ID,
		Author:   author,
	}
}

func (b Book) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Title, validation.Required),
		validation.Field(&b.ISBN, validation.Required),
		validation.Field(&b.Author, validation.Required))
}
