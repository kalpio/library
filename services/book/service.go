package book

import (
	"errors"
	"fmt"
	"library/domain"
	"library/infrastructure/repository"
	"library/services/author"
	"strings"

	"github.com/google/uuid"
)

type IBookService interface {
	Create(id uuid.UUID,
		title, isbn, description string,
		authorID uuid.UUID) (*domain.Book, error)

	Edit(id uuid.UUID,
		title, isbn, description string,
		authorID uuid.UUID) (*domain.Book, error)

	GetByID(id uuid.UUID) (*domain.Book, error)

	GetAll() ([]domain.Book, error)

	Delete(id uuid.UUID) (bool, error)
}

type bookService struct {
	db        domain.IDatabase
	authorSrv author.IAuthorService
}

func newBookService(db domain.IDatabase, authorSrv author.IAuthorService) IBookService {
	return &bookService{db: db, authorSrv: authorSrv}
}

func (b *bookService) Create(id uuid.UUID,
	title, isbn, description string,
	authorID uuid.UUID) (*domain.Book, error) {

	err := b.exists(isbn)
	if err != nil {
		return nil, err
	}

	var bookAuthor *domain.Author
	bookAuthor, err = b.getAuthor(authorID)
	if err != nil {
		return nil, err
	}

	model := domain.NewBook(id, title, isbn, description, bookAuthor)

	var result domain.Book
	result, err = repository.Save(b.db, *model)
	if err != nil {
		return nil, fmt.Errorf("book service: could not save book: %w", err)
	}

	return &result, nil
}

var ErrAlreadyExists = errors.New("book service: book already exists")

func (b *bookService) exists(isbn string) error {
	var (
		err   error
		value domain.Book
	)
	value, err = repository.GetByColumns[domain.Book](b.db, map[string]interface{}{"isbn": isbn})
	if err != nil {
		return fmt.Errorf("book service: an error during check book exists %w", err)
	}

	var exists = strings.Compare(isbn, value.ISBN) == 0
	if exists {
		return ErrAlreadyExists
	}

	return nil
}

func (b *bookService) getAuthor(id uuid.UUID) (*domain.Author, error) {
	bookAuthor, err := b.authorSrv.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("book service: could not find book author: %w", err)
	}

	return bookAuthor, nil
}

func (b *bookService) Edit(id uuid.UUID,
	title, isbn, description string,
	authorID uuid.UUID) (*domain.Book, error) {

	return nil, nil
}

func (b *bookService) GetByID(id uuid.UUID) (*domain.Book, error) {
	return nil, nil
}

func (b *bookService) GetAll() ([]domain.Book, error) {
	return nil, nil
}

func (b *bookService) Delete(id uuid.UUID) (bool, error) {
	return false, nil
}
