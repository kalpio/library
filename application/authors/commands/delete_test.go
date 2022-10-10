package commands

import (
	"context"
	"library/application/authors/events"
	"library/domain"
	domain_events "library/domain/events"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthor_DeleteCommandHandler_RaisedAuthorDeletedEvent(t *testing.T) {
	ass := assert.New(t)
	registerEvents(ass)
	mckSevice := new(authorServiceMock)

	authorID := uuid.New()
	expectedAuthorID := domain.AuthorID(authorID.String())
	mckSevice.
		On("Delete", authorID).
		Return(true, nil)

	commandHandler := NewDeleteAuthorCommandHandler(nil, mckSevice)
	response, err := commandHandler.Handle(context.Background(), NewDeleteAuthorCommand(expectedAuthorID))

	ass.NoError(err)
	mckSevice.AssertExpectations(t)
	ass.True(response.Succeeded)

	notifications := domain_events.GetEvents(&events.AuthorDeletedEvent{})
	ass.Equal(1, len(notifications))

	notification := notifications[0]
	ass.Equal(authorID, notification.AuthorID)
}
