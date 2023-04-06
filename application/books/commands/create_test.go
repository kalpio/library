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

func TestBook_CreateCommandHandler_RaisedBookCreatedEvent(t *testing.T) {
	ass := assert.New(t)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	mckService := bookSrv.(*bookstest.BookServiceMock)
	expectedBook := bookstest.CreateBook()

	mckService.
		On("Create",
			expectedBook.ID,
			expectedBook.Title,
			expectedBook.ISBN,
			expectedBook.Description,
			expectedBook.AuthorID).
		Return(expectedBook, nil)

	commandHandler := commands.NewCreateBookCommandHandler(nil, mckService)
	_, err = commandHandler.Handle(context.Background(),
		commands.NewCreateBookCommand(expectedBook.ID,
			expectedBook.Title,
			expectedBook.ISBN,
			expectedBook.Description,
			domain.AuthorID(expectedBook.AuthorID.String())))

	ass.NoError(err)
	mckService.AssertExpectations(t)

	tests.Wait(
		func() bool {
			return len(domainEvents.GetEvents[*events.BookCreatedEvent]()) > 0
		},
		1*time.Second)

	notifications := domainEvents.GetEvents[*events.BookCreatedEvent]()
	ass.Len(notifications, 1)

	ass.Equal(expectedBook.ID, notifications[0].BookID)
	ass.Equal(expectedBook.Title, notifications[0].Title)
	ass.Equal(expectedBook.ISBN, notifications[0].ISBN)
	ass.Equal(expectedBook.Description, notifications[0].Description)
	ass.Equal(expectedBook.AuthorID, notifications[0].AuthorID)
}
