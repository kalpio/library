package book

import (
	"library/domain"
	"library/infrastructure/repository"
	"library/infrastructure/repository/testutils"
	"library/random"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveNewBook(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	id := domain.NewBookID()
	title := random.String(100)
	isbn := random.String(13)
	format := random.String(3)
	author := domain.NewAuthor(
		domain.NewAuthorID(),
		random.String(10),
		random.String(10),
		random.String(10))

	book := domain.NewBook(id, title, isbn, format, author)
	ass := assert.New(t)
	result, err := repository.Save(*book)

	ass.NoError(err)
	ass.Equal(result.ID, book.ID)
	ass.Equal(result.Title, title)
	ass.Equal(result.ISBN, isbn)
	ass.Equal(result.AuthorID, author.ID)
	ass.Equal(result.Author, author)
}

func TestGetByISBN(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	ass := assert.New(t)

	expected := createNewBookInDB(t)
	columnValue := map[string]interface{}{"ISBN": expected.ISBN}
	got, err := repository.GetByColumns[domain.Book](columnValue)

	ass.NoError(err)
	ass.Equal(got.ID, expected.ID)
}

func TestGetAll(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	var expected []*domain.Book
	expected = append(expected, createNewBookInDB(t))
	expected = append(expected, createNewBookInDB(t))
	expected = append(expected, createNewBookInDB(t))

	page := 1
	pageSize := 50
	results, err := repository.GetAll[domain.Book](page, pageSize)
	ass.NoError(err)
	ass.Equal(len(results), 3)

	for _, v := range expected {
		assertThatContainsBook(ass, results, v)
	}
}

func assertThatContainsBook(ass *assert.Assertions, results []domain.Book, expected *domain.Book) {
	for _, v := range results {
		if v.ID == expected.ID {
			assertThatTheyAreSameBook(ass, v, expected)
			return
		}
	}
}

func assertThatTheyAreSameBook(ass *assert.Assertions, get domain.Book, expected *domain.Book) {
	ass.Equal(get.Title, expected.Title)
	ass.Equal(get.ISBN, expected.ISBN)
	ass.Equal(get.AuthorID, expected.AuthorID)

	ass.Equal(get.CreatedAt.UTC(), expected.CreatedAt.UTC())
	ass.Equal(get.UpdatedAt.UTC(), expected.UpdatedAt.UTC())
}

func TestDelete(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	ass := assert.New(t)

	book := createNewBookInDB(t)
	err := repository.Delete[domain.Book](book.ID.UUID())

	ass.NoError(err)
}

func createNewBookInDB(t *testing.T) *domain.Book {
	b := domain.NewBook(
		domain.NewBookID(),
		random.String(100),
		random.String(13),
		random.String(3),
		domain.NewAuthor(
			domain.NewAuthorID(),
			random.String(6),
			random.String(6),
			random.String(6)))

	ass := assert.New(t)
	result, err := repository.Save(*b)
	ass.NoError(err)

	return &result
}
