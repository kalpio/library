package commands_test

import (
	"context"
	"library/application/authors/commands"
	"library/application/authors/events"
	domainEvents "library/domain/events"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthor_CreateCommandHandler_RaisedAuthorCreatedEvent(t *testing.T) {
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

	commandHandler := commands.NewCreateAuthorCommandHandler(nil, mckService)
	_, err := commandHandler.Handle(context.Background(),
		commands.NewCreateAuthorCommand(expectedAuthor.ID,
			expectedAuthor.FirstName,
			expectedAuthor.MiddleName,
			expectedAuthor.LastName))

	ass.NoError(err)
	mckService.AssertExpectations(t)

	notifications := domainEvents.GetEvents[*events.AuthorCreatedEvent]()
	ass.Equal(1, len(notifications))

	notification := notifications[0]
	ass.Equal(expectedAuthor.ID, notification.AuthorID)
	ass.Equal(expectedAuthor.FirstName, notification.FirstName)
	ass.Equal(expectedAuthor.MiddleName, notification.MiddleName)
	ass.Equal(expectedAuthor.LastName, notification.LastName)
}
