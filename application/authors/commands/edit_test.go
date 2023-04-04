package commands_test

import (
	"context"
	"library/application/authors/commands"
	"library/application/authors/events"
	"library/domain"
	domainEvents "library/domain/events"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthor_EditCommandHandler_Raised_AuthorEditedEvent(t *testing.T) {
	ass := assert.New(t)
	registerEvents(ass)

	mckService := new(authorServiceMock)
	expectedID := domain.NewAuthorID()
	editedAuthor := createAuthor()
	editedAuthor.SetID(expectedID)

	mckService.
		On("Edit",
			expectedID.UUID(),
			editedAuthor.FirstName,
			editedAuthor.MiddleName,
			editedAuthor.LastName).
		Return(editedAuthor, nil)

	commandHandler := commands.NewEditAuthorCommandHandler(nil, mckService)
	_, err := commandHandler.Handle(context.Background(),
		commands.NewEditAuthorCommand(expectedID,
			editedAuthor.FirstName,
			editedAuthor.MiddleName,
			editedAuthor.LastName))

	ass.NoError(err)
	mckService.AssertExpectations(t)

	notifications := domainEvents.GetEvents[*events.AuthorEditedEvent]()
	ass.Equal(1, len(notifications))

	notification := notifications[0]
	ass.Equal(expectedID, notification.AuthorID)
	ass.Equal(editedAuthor.FirstName, notification.FirstName)
	ass.Equal(editedAuthor.MiddleName, notification.MiddleName)
	ass.Equal(editedAuthor.LastName, notification.LastName)
}
