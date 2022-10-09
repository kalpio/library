package commands

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"library/application/authors/events"
	"library/domain"
	domain_events "library/domain/events"
	"testing"
)

func TestAuthor_Delete(t *testing.T) {
	ass := assert.New(t)
	registerEvents(ass)
	mck := new(authorServiceMock)

	authorID := uuid.New()
	expectedAuthorID := domain.AuthorID(authorID.String())
	mck.
		On("Delete", authorID).
		Return(true, nil)

	commandHandler := NewDeleteAuthorCommandHandler(nil, mck)
	response, err := commandHandler.Handle(context.Background(), NewDeleteAuthorCommand(expectedAuthorID))

	ass.NoError(err)
	mck.AssertExpectations(t)
	ass.True(response.Succeeded)

	notifications := domain_events.GetEvents(&events.AuthorDeletedEvent{})
	ass.Equal(1, len(notifications))

	notification := notifications[0]
	ass.Equal(authorID, notification.AuthorID)
}
