package commands_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"library/application/books/commands"
	"library/application/books/events"
	"library/domain"
	domainEvents "library/domain/events"
	"testing"
)

func TestBook_CreateCommandHandler_RaisedBookCreatedEvent(t *testing.T) {
	ass := assert.New(t)
	registerEvents(ass)

	mckService := new(bookServiceMock)
	expectedBook := createBook()

	mckService.
		On("Create",
			expectedBook.ID,
			expectedBook.Title,
			expectedBook.ISBN,
			expectedBook.Description,
			expectedBook.AuthorID).
		Return(expectedBook, nil)

	commandHandler := commands.NewCreateBookCommandHandler(nil, mckService)
	_, err := commandHandler.Handle(context.Background(),
		commands.NewCreateBookCommand(domain.BookID(expectedBook.ID.String()),
			expectedBook.Title,
			expectedBook.ISBN,
			expectedBook.Description,
			domain.AuthorID(expectedBook.AuthorID.String())))

	ass.NoError(err)
	mckService.AssertExpectations(t)

	notifications := domainEvents.GetEvents(&events.BookCreatedEvent{})
	ass.Len(notifications, 1)

	ass.Equal(expectedBook.ID, uuid.MustParse(string(notifications[0].BookID)))
	ass.Equal(expectedBook.Title, notifications[0].Title)
	ass.Equal(expectedBook.ISBN, notifications[0].ISBN)
	ass.Equal(expectedBook.Description, notifications[0].Description)
	ass.Equal(expectedBook.AuthorID, uuid.MustParse(string(notifications[0].AuthorID)))
}
