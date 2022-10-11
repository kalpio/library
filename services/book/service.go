package book

import (
	"github.com/sirupsen/logrus"
	"library/domain"

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
	db domain.IDatabase
}

func NewBookService(db domain.IDatabase) IBookService {
	return &bookService{db: db}
}

func (b *bookService) Create(id uuid.UUID,
	title, isbn, description string,
	authorID uuid.UUID) (*domain.Book, error) {

	logrus.Fatal("undefined")
	return nil, nil
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