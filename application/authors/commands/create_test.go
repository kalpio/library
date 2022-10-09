package commands

import (
	"context"
	"github.com/stretchr/testify/assert"
	"library/application/authors/events"
	domain_events "library/domain/events"
	"testing"
)

func TestAuthor_Create(t *testing.T) {
	ass := assert.New(t)
	registerEvents(ass)

	mckService := new(authorServiceMock)
	expectedAuthor := createAuthor()

	mckService.
		On("Create",
			expectedAuthor.ID,
			expectedAuthor.FirstName,
			expectedAuthor.MiddleName,
			expectedAuthor.LastName).
		Return(expectedAuthor, nil)

	commandHandler := NewCreateAuthorCommandHandler(nil, mckService)
	_, err := commandHandler.Handle(context.Background(),
		NewCreateAuthorCommand(expectedAuthor.ID,
			expectedAuthor.FirstName,
			expectedAuthor.MiddleName,
			expectedAuthor.LastName))

	ass.NoError(err)
	mckService.AssertExpectations(t)
	notifications := domain_events.GetEvents(&events.AuthorCreatedEvent{})
	ass.Equal(1, len(notifications))
	notification := notifications[0]
	ass.Equal(expectedAuthor.ID, notification.AuthorID)
	ass.Equal(expectedAuthor.FirstName, notification.FirstName)
	ass.Equal(expectedAuthor.MiddleName, notification.MiddleName)
	ass.Equal(expectedAuthor.LastName, notification.LastName)
}
