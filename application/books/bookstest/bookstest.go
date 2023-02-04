package bookstest

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"library/application/books"
	"library/domain"
	"library/ioc"
	"library/random"
	"library/services/book"
)

type fakeDatabase struct {
}

func (fdb *fakeDatabase) GetDB() *gorm.DB {
	return nil
}

func Initialize() error {
	var lastErr error
	if err := ioc.AddSingleton[domain.IDatabase](new(fakeDatabase)); err != nil {
		lastErr = err
	}

	if err := ioc.AddSingleton[book.IBookService](new(BookServiceMock)); err != nil {
		lastErr = err
	}

	bookRegister := books.NewBookRegister()
	if err := bookRegister.Register(); err != nil {
		lastErr = err
	}

	return lastErr
}

type BookServiceMock struct {
	mock.Mock
}

func (b *BookServiceMock) Create(id uuid.UUID,
	title, isbn, description string,
	authorID uuid.UUID) (*domain.Book, error) {

	args := b.Called(id, title, isbn, description, authorID)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (b *BookServiceMock) Edit(id uuid.UUID,
	title, isbn, description string,
	authorID uuid.UUID) (*domain.Book, error) {

	args := b.Called(id, title, isbn, description, authorID)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (b *BookServiceMock) GetByID(id uuid.UUID) (*domain.Book, error) {
	args := b.Called(id)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (b *BookServiceMock) GetAll() ([]domain.Book, error) {
	args := b.Called()
	return args.Get(0).([]domain.Book), args.Error(1)
}

func (b *BookServiceMock) Delete(id uuid.UUID) (bool, error) {
	args := b.Called(id)
	return args.Bool(0), args.Error(1)
}

func CreateBook() *domain.Book {
	return domain.NewBook(uuid.New(),
		random.String(20),
		random.String(12),
		random.String(120),
		domain.NewAuthor(uuid.New(),
			random.String(20),
			random.String(20),
			random.String(20)))
}
