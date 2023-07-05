package commands_test

import (
	log "github.com/sirupsen/logrus"
	"library/application/authors/events"
	"library/domain"
	"library/random"
	"time"

	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/mock"
)

func init() {
	if err := registerEvents(); err != nil {
		log.Fatalln(err)
	}
}
func registerEvents() error {
	if err := mediatr.RegisterNotificationHandler[*events.AuthorCreatedEvent](&events.AuthorCreatedEventHandler{}); err != nil {
		return err
	}

	if err := mediatr.RegisterNotificationHandler[*events.AuthorEditedEvent](&events.AuthorEditedEventHandler{}); err != nil {
		return err
	}

	if err := mediatr.RegisterNotificationHandler[*events.AuthorDeletedEvent](&events.AuthorDeletedEventHandler{}); err != nil {
		return err
	}

	return nil
}

type authorServiceMock struct {
	mock.Mock
}

func (a *authorServiceMock) Create(id domain.AuthorID, firstName, middleName, lastName string) (*domain.Author, error) {
	args := a.Called(id, firstName, middleName, lastName)
	return args.Get(0).(*domain.Author), args.Error(1)
}

func (a *authorServiceMock) Edit(id domain.AuthorID, firstName, middleName, lastName string) (*domain.Author, error) {
	args := a.Called(id, firstName, middleName, lastName)
	return args.Get(0).(*domain.Author), args.Error(1)
}

func (a *authorServiceMock) GetByID(id domain.AuthorID) (*domain.Author, error) {
	args := a.Called(id)
	return args.Get(0).(*domain.Author), args.Error(1)
}

func (a *authorServiceMock) GetAll(page, size int) ([]domain.Author, error) {
	args := a.Called(page, size)
	return args.Get(0).([]domain.Author), args.Error(1)
}

func (a *authorServiceMock) Delete(id domain.AuthorID) error {
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
