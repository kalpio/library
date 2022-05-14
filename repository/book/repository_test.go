package book

import (
	"github.com/matryer/is"
	"gorm.io/gorm"
	"library/models"
	"library/random"
	"library/repository"
	"library/repository/testutils"
	"testing"
)

func Test_SaveNewBook(t *testing.T) {
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

	book0 := models.NewBook(title, isbn, format, author)
	book0.Content = content
	book0.Version = version
	iss := is.New(t)
	err := Save(db, book0)

	iss.NoErr(err)
	iss.True(book0.ID > 0)
	iss.Equal(book0.Title, title)
	iss.Equal(book0.ISBN, isbn)
	iss.Equal(book0.Content, content)
	iss.Equal(book0.Format, format)
	iss.Equal(book0.Version, version)
	iss.Equal(book0.AuthorID, author.ID)
	iss.Equal(book0.Author, author)
}

func Test_GetByISBN(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	iss := is.New(t)

	expected := createNewBookInDB(db, t)
	get, err := GetByISBN(db, expected.ISBN)

	iss.NoErr(err)
	iss.Equal(get.ID, expected.ID)
}

func Test_GetAll(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	iss := is.New(t)
	var expected []*models.Book
	expected = append(expected, createNewBookInDB(db, t))
	expected = append(expected, createNewBookInDB(db, t))
	expected = append(expected, createNewBookInDB(db, t))

	results, err := repository.GetAll[models.Book](db)
	iss.NoErr(err)
	iss.Equal(len(results), 3)

	for _, v := range expected {
		assertThatContainsBook(t, results, v)
	}
}

func assertThatContainsBook(t *testing.T, results []*models.Book, expected *models.Book) {
	for _, v := range results {
		if v.ID == expected.ID {
			assertThatTheyAreSameBook(t, v, expected)
			return
		}
	}
}

func assertThatTheyAreSameBook(t *testing.T, get *models.Book, expected *models.Book) {
	iss := is.New(t)
	iss.Equal(get.Title, expected.Title)
	iss.Equal(get.ISBN, expected.ISBN)
	iss.Equal(get.Content, expected.Content)
	iss.Equal(get.Format, expected.Format)
	iss.Equal(get.Version, expected.Version)
	iss.Equal(get.AuthorID, expected.AuthorID)

	iss.Equal(get.CreatedAt.UTC(), expected.CreatedAt.UTC())
	iss.Equal(get.UpdatedAt.UTC(), expected.UpdatedAt.UTC())
	if expected.DeletedAt.Valid {
		iss.True(get.DeletedAt.Valid)
		iss.Equal(get.DeletedAt.Time.UTC(), expected.DeletedAt.Time.UTC())
	}
}

func Test_Delete(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	iss := is.New(t)
	b := createNewBookInDB(db, t)
	err := Delete(db, b.ID)

	iss.NoErr(err)
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
	iss := is.New(t)
	err := Save(db, b)
	iss.NoErr(err)

	return b
}
