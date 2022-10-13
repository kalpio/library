package commands_test

import (
	"context"
	"library/application/authors/commands"
	"library/application/authors/events"
	"library/domain"
	domainEvents "library/domain/events"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthor_EditCommandHandler_Raised_AuthorEditedEvent(t *testing.T) {
	ass := assert.New(t)
	registerEvents(ass)

	mckService := new(authorServiceMock)
	payloadAuthorID := uuid.New()
	editedAuthor := createAuthor()
	editedAuthor.ID = payloadAuthorID

	mckService.
		On("Edit",
			payloadAuthorID,
			editedAuthor.FirstName,
			editedAuthor.MiddleName,
			editedAuthor.LastName).
		Return(editedAuthor, nil)

	commandHandler := commands.NewEditAuthorCommandHandler(nil, mckService)
	_, err := commandHandler.Handle(context.Background(),
		commands.NewEditAuthorCommand(domain.AuthorID(payloadAuthorID.String()),
			editedAuthor.FirstName,
			editedAuthor.MiddleName,
			editedAuthor.LastName))

	ass.NoError(err)
	mckService.AssertExpectations(t)

	notifications := domainEvents.GetEvents(&events.AuthorEditedEvent{})
	ass.Equal(1, len(notifications))

	notification := notifications[0]
	ass.Equal(payloadAuthorID, notification.AuthorID)
	ass.Equal(editedAuthor.FirstName, notification.FirstName)
	ass.Equal(editedAuthor.MiddleName, notification.MiddleName)
	ass.Equal(editedAuthor.LastName, notification.LastName)
}
