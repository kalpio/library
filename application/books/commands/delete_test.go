package commands_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"library/application/books/bookstest"
	"library/application/books/commands"
	"library/application/books/events"
	"library/domain"
	domainEvents "library/domain/events"
	"library/ioc"
	"library/services/book"
	"library/tests"
	"testing"
	"time"
)

func TestBook_DeleteCommandHandler_RaisedBookDeletedEvent(t *testing.T) {
	ass := assert.New(t)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	mckService := bookSrv.(*bookstest.BookServiceMock)
	expectedID := domain.NewBookID()
	mckService.
		On("Delete", expectedID).
		Return(nil)

	commandHandler := commands.NewDeleteBookCommandHandler(nil, mckService)
	response, err := commandHandler.Handle(context.Background(),
		commands.NewDeleteBookCommand(expectedID))
	ass.NoError(err)
	mckService.AssertExpectations(t)
	ass.True(response.Succeeded)

	tests.Wait(
		func() bool {
			return len(domainEvents.GetEvents[*events.BookDeletedEvent]()) > 0
		},
		1*time.Second)

	notifications := domainEvents.GetEvents[*events.BookDeletedEvent]()
	ass.Equal(1, len(notifications))

	notification := notifications[0]
	ass.Equal(expectedID, notification.BookID)
}
