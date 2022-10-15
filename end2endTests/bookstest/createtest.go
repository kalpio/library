package bookstest

import (
	"github.com/stretchr/testify/assert"
	"library/application"
	"library/domain"
	"library/random"
	"net/http"
	"testing"
)

var a application.App

func SetApp(app application.App) {
	a = app
}

func PostNewBook(t *testing.T) {
	ass := assert.New(t)
	bookAuthor, err := postAuthor(a, postAuthorDto{
		FirstName:  random.String(20),
		MiddleName: random.String(20),
		LastName:   random.String(20),
	})

	ass.NoError(err)

	bookDto := generateBookDto(bookAuthor.ID)
	rr := postBookRR(a, bookDto)

	ass.NotNil(rr)
	ass.Equal(http.StatusCreated, rr.Code)

	bookResponse, err := getBookFromResponseBody(rr.Body)
	ass.NoError(err)
	assertCreatedBook(ass, bookDto, bookResponse)
}

func assertCreatedBook(ass *assert.Assertions, expected postBookDto, actual *domain.Book) {
	ass.Equal(expected.ID, actual.ID.String())
	ass.Equal(expected.Title, actual.Title)
	ass.Equal(expected.ISBN, actual.ISBN)
	ass.Equal(expected.Description, actual.Description)
	ass.Equal(expected.AuthorID, actual.AuthorID.String())
}
