package commands

import (
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"library/application/books/events"
	"library/domain"
	"library/random"
)

func registerEvents(ass *assert.Assertions) {
	if err := mediatr.RegisterNotificationHandler[*events.BookCreatedEvent](
		&events.BookCreatedEventHandler{}); err != nil {
		ass.NoError(err)
	}
}

type bookServiceMock struct {
	mock.Mock
}

func (b *bookServiceMock) Create(id uuid.UUID,
	title, isbn, description string,
	authorID uuid.UUID) (*domain.Book, error) {

	args := b.Called(id, title, isbn, description, authorID)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (b *bookServiceMock) Edit(id uuid.UUID,
	title, isbn, description string,
	authorID uuid.UUID) (*domain.Book, error) {

	args := b.Called(id, title, isbn, description, authorID)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (b *bookServiceMock) GetByID(id uuid.UUID) (*domain.Book, error) {
	args := b.Called(id)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (b *bookServiceMock) GetAll() ([]domain.Book, error) {
	args := b.Called()
	return args.Get(0).([]domain.Book), args.Error(1)
}

func (b *bookServiceMock) Delete(id uuid.UUID) (bool, error) {
	args := b.Called(id)
	return args.Bool(0), args.Error(1)
}

func createBook() *domain.Book {
	return domain.NewBook(uuid.New(),
		random.String(20),
		random.String(12),
		random.String(120),
		domain.NewAuthor(uuid.New(),
			random.String(20),
			random.String(20),
			random.String(20)))
}
