package queries_test

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"
	"library/application/books/bookstest"
	"library/application/books/queries"
	"library/ioc"
	"library/services/book"
	"testing"
)

func TestBook_QueryHandler_ReturnCorrectData(t *testing.T) {
	ass := assert.New(t)

	err := bookstest.Initialize()
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	mckService := bookSrv.(*bookstest.BookServiceMock)
	expectedBook := bookstest.CreateBook()

	mckService.
		On("GetByID",
			expectedBook.ID).
		Return(expectedBook, nil)

	query := queries.NewGetBookByIDQuery(expectedBook.ID)
	response, err := mediatr.Send[
		*queries.GetBookByIDQuery,
		*queries.GetBookByIDQueryResponse](
		context.Background(),
		query)

	ass.NoError(err)
	ass.NotNil(response)
	ass.Equal(expectedBook.ID, response.BookID)
	ass.Equal(expectedBook.ISBN, response.ISBN)
	ass.Equal(expectedBook.Title, response.Title)
	ass.Equal(expectedBook.Description, response.Description)
	ass.Equal(expectedBook.AuthorID, response.AuthorID)
}
