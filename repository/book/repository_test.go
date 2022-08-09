package book

import (
	"library/models"
	"library/random"
	"library/repository"
	"library/repository/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSaveNewBook(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	title := random.RandomString(100)
	isbn := random.RandomString(13)
	content := []byte(random.RandomString(256))
	format := random.RandomString(3)
	version := random.RandomString(4)
	author := &models.Author{
		FirstName:  random.RandomString(10),
		MiddleName: random.RandomString(10),
		LastName:   random.RandomString(10),
	}

	book := models.NewBook(title, isbn, format, author)
	book.Content = content
	book.Version = version
	ass := assert.New(t)
	result, err := repository.Save(db, *book)

	ass.NoError(err)
	ass.Greater(result.ID, uint(0))
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
	got, err := repository.GetByColumns[models.Book](db, columnValue)

	ass.NoError(err)
	ass.Equal(got.ID, expected.ID)
}

func TestGetAll(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	var expected []*models.Book
	expected = append(expected, createNewBookInDB(db, t))
	expected = append(expected, createNewBookInDB(db, t))
	expected = append(expected, createNewBookInDB(db, t))

	results, err := repository.GetAll[models.Book](db)
	ass.NoError(err)
	ass.Equal(len(results), 3)

	for _, v := range expected {
		assertThatContainsBook(ass, results, v)
	}
}

func assertThatContainsBook(ass *assert.Assertions, results []models.Book, expected *models.Book) {
	for _, v := range results {
		if v.ID == expected.ID {
			assertThatTheyAreSameBook(ass, v, expected)
			return
		}
	}
}

func assertThatTheyAreSameBook(ass *assert.Assertions, get models.Book, expected *models.Book) {
	ass.Equal(get.Title, expected.Title)
	ass.Equal(get.ISBN, expected.ISBN)
	ass.Equal(get.Content, expected.Content)
	ass.Equal(get.Format, expected.Format)
	ass.Equal(get.Version, expected.Version)
	ass.Equal(get.AuthorID, expected.AuthorID)

	ass.Equal(get.CreatedAt.UTC(), expected.CreatedAt.UTC())
	ass.Equal(get.UpdatedAt.UTC(), expected.UpdatedAt.UTC())

	if expected.DeletedAt.Valid {
		ass.True(get.DeletedAt.Valid)
		ass.Equal(get.DeletedAt.Time.UTC(), expected.DeletedAt.Time.UTC())
	}
}

func TestDelete(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	ass := assert.New(t)

	book := createNewBookInDB(db, t)
	rowsAffected, err := repository.Delete[models.Book](db, book.ID)

	ass.NoError(err)
	ass.Greater(rowsAffected, int64(0))
}

func createNewBookInDB(db *gorm.DB, t *testing.T) *models.Book {
	b := models.NewBook(
		random.RandomString(100),
		random.RandomString(13),
		random.RandomString(3),
		models.NewAuthor(
			random.RandomString(6),
			random.RandomString(6),
			random.RandomString(6)))

	ass := assert.New(t)
	result, err := repository.Save(db, *b)
	ass.NoError(err)

	return &result
}
