package commands_test

import (
	"library/application/authors/events"
	"library/domain"
	"library/random"
	"time"

	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func registerEvents(ass *assert.Assertions) {
	if err := mediatr.RegisterNotificationHandler[*events.AuthorCreatedEvent](&events.AuthorCreatedEventHandler{}); err != nil {
		ass.NoError(err)
	}

	if err := mediatr.RegisterNotificationHandler[*events.AuthorEditedEvent](&events.AuthorEditedEventHandler{}); err != nil {
		ass.NoError(err)
	}

	if err := mediatr.RegisterNotificationHandler[*events.AuthorDeletedEvent](&events.AuthorDeletedEventHandler{}); err != nil {
		ass.NoError(err)
	}
}

type authorServiceMock struct {
	mock.Mock
}

func (a *authorServiceMock) Create(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	args := a.Called(id, firstName, middleName, lastName)
	return args.Get(0).(*domain.Author), args.Error(1)
}

func (a *authorServiceMock) Edit(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	args := a.Called(id, firstName, middleName, lastName)
	return args.Get(0).(*domain.Author), args.Error(1)
}

func (a *authorServiceMock) GetByID(id uuid.UUID) (*domain.Author, error) {
	args := a.Called(id)
	return args.Get(0).(*domain.Author), args.Error(1)
}

func (a *authorServiceMock) GetAll() ([]domain.Author, error) {
	args := a.Called()
	return args.Get(0).([]domain.Author), args.Error(1)
}

func (a *authorServiceMock) Delete(id uuid.UUID) error {
	args := a.Called(id)
	return args.Error(0)
}

func createAuthor() *domain.Author {
	return &domain.Author{
		Entity: domain.Entity[domain.AuthorID]{
			ID:        domain.NewAuthorID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		FirstName:  random.String(20),
		MiddleName: random.String(20),
		LastName:   random.String(20),
		Books:      nil,
	}
}
