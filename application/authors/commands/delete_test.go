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

func TestAuthor_DeleteCommandHandler_RaisedAuthorDeletedEvent(t *testing.T) {
	ass := assert.New(t)
	registerEvents(ass)
	mckService := new(authorServiceMock)

	expectedID := domain.NewAuthorID()
	mckService.
		On("Delete", expectedID.UUID()).
		Return(nil)

	commandHandler := commands.NewDeleteAuthorCommandHandler(nil, mckService)
	response, err := commandHandler.Handle(context.Background(),
		commands.NewDeleteAuthorCommand(expectedID))

	ass.NoError(err)
	mckService.AssertExpectations(t)
	ass.True(response.Succeeded)

	notifications := domainEvents.GetEvents[*events.AuthorDeletedEvent]()
	ass.Equal(1, len(notifications))

	notification := notifications[0]
	ass.Equal(expectedID, notification.AuthorID)
}
