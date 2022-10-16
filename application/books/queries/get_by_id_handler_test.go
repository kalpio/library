package queries_test

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/stretchr/testify/assert"
	"library/application/books/bookstest"
	"library/application/books/queries"
	"library/ioc"
	"library/register"

	"library/domain"
	"testing"
)

func TestBook_QueryHandler_ReturnCorrectData(t *testing.T) {
	ass := assert.New(t)

	err := bookstest.Register()
	ass.NoError(err)

	reg, err := ioc.Get[register.IRegister[*domain.Book]]()
	ass.NoError(err)

	err = reg.Register()
	ass.NoError(err)

	mckService := new(bookstest.BookServiceMock)
	expectedBook := bookstest.CreateBook()

	mckService.
		On("GetByID",
			expectedBook.ID).
		Return(expectedBook, nil)

	query := queries.NewGetBookByIDQuery(domain.BookID(expectedBook.ID.String()))
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
