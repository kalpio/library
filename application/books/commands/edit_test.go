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

func TestBook_EditCommandHandler_RaisedBookEditEvent(t *testing.T) {
	ass := assert.New(t)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	mckService := bookSrv.(*bookstest.BookServiceMock)
	expectedID := domain.NewBookID()
	editedBook := bookstest.CreateBook()
	editedBook.SetID(expectedID)

	mckService.
		On("Edit",
			expectedID,
			editedBook.Title,
			editedBook.ISBN,
			editedBook.Description,
			editedBook.AuthorID).
		Return(editedBook, nil)

	commandHandler := commands.NewEditBookCommandHandler(nil, mckService)
	_, err = commandHandler.Handle(context.Background(),
		commands.NewEditBookCommand(expectedID,
			editedBook.Title,
			editedBook.ISBN,
			editedBook.Description,
			editedBook.AuthorID))

	ass.NoError(err)
	mckService.AssertExpectations(t)

	tests.Wait(
		func() bool {
			return len(domainEvents.GetEvents[*events.BookEditedEvent]()) > 0
		},
		1*time.Second)

	notifications := domainEvents.GetEvents[*events.BookEditedEvent]()
	ass.Equal(1, len(notifications))

	notification := notifications[0]
	ass.Equal(expectedID, notification.BookID)
	ass.Equal(editedBook.Title, notification.Title)
	ass.Equal(editedBook.ISBN, notification.ISBN)
	ass.Equal(editedBook.AuthorID, notification.AuthorID)
}
