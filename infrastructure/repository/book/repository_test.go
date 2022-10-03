package book

import (
	"github.com/google/uuid"
	"library/domain"
	"library/infrastructure/repository"
	"library/infrastructure/repository/testutils"
	"library/random"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveNewBook(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	id := uuid.New()
	title := random.String(100)
	isbn := random.String(13)
	content := []byte(random.String(256))
	format := random.String(3)
	version := random.String(4)
	author := domain.NewAuthor(
		uuid.New(),
		random.String(10),
		random.String(10),
		random.String(10))

	book := domain.NewBook(id, title, isbn, format, author)
	book.Content = content
	book.Version = version
	ass := assert.New(t)
	result, err := repository.Save(db, *book)

	ass.NoError(err)
	ass.Equal(result.ID, book.ID)
	ass.Equal(result.Title, title)
	ass.Equal(result.ISBN, isbn)
	ass.Equal(result.Content, content)
	ass.Equal(result.Format, format)
	ass.Equal(result.Version, version)
	ass.Equal(result.AuthorID, author.ID)
	ass.Equal(result.Author, author)
}

func TestGetByISBN(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	ass := assert.New(t)

	expected := createNewBookInDB(db, t)
	columnValue := map[string]interface{}{"ISBN": expected.ISBN}
	got, err := repository.GetByColumns[domain.Book](db, columnValue)

	ass.NoError(err)
	ass.Equal(got.ID, expected.ID)
}

func TestGetAll(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	var expected []*domain.Book
	expected = append(expected, createNewBookInDB(db, t))
	expected = append(expected, createNewBookInDB(db, t))
	expected = append(expected, createNewBookInDB(db, t))

	results, err := repository.GetAll[domain.Book](db)
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
	ass.Equal(get.Content, expected.Content)
	ass.Equal(get.Format, expected.Format)
	ass.Equal(get.Version, expected.Version)
	ass.Equal(get.AuthorID, expected.AuthorID)

	ass.Equal(get.CreatedAt.UTC(), expected.CreatedAt.UTC())
	ass.Equal(get.UpdatedAt.UTC(), expected.UpdatedAt.UTC())
}

func TestDelete(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	ass := assert.New(t)

	book := createNewBookInDB(db, t)
	rowsAffected, err := repository.Delete[domain.Book](db, book.ID)

	ass.NoError(err)
	ass.Greater(rowsAffected, int64(0))
}

func createNewBookInDB(db domain.Database, t *testing.T) *domain.Book {
	b := domain.NewBook(
		uuid.New(),
		random.String(100),
		random.String(13),
		random.String(3),
		domain.NewAuthor(
			uuid.New(),
			random.String(6),
			random.String(6),
			random.String(6)))

	ass := assert.New(t)
	result, err := repository.Save(db, *b)
	ass.NoError(err)

	return &result
}
