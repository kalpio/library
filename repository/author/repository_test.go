package author

import (
	"library/models"
	"library/random"
	"library/repository"
	"library/repository/testutils"
	"testing"

	"github.com/matryer/is"
	"gorm.io/gorm"
)

func Test_SaveNewAuthor(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	firstName := random.RandomString(10)
	middleName := random.RandomString(10)
	lastName := random.RandomString(10)
	author0 := models.NewAuthor(firstName, middleName, lastName)

	iss := is.New(t)
	err := Save(db, author0)

	iss.NoErr(err)
	iss.True(author0.ID > 0)
	iss.Equal(author0.FirstName, firstName)
	iss.Equal(author0.MiddleName, middleName)
	iss.Equal(author0.LastName, lastName)
}

func Test_TryAddExistingAuthor(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	iss := is.New(t)
	firstName := random.RandomString(10)
	middleName := random.RandomString(10)
	lastName := random.RandomString(10)

	author0 := models.NewAuthor(firstName, middleName, lastName)
	err := Save(db, author0)
	iss.NoErr(err)
	iss.Equal(author0.FirstName, firstName)
	iss.Equal(author0.MiddleName, middleName)
	iss.Equal(author0.LastName, lastName)

	author1 := models.NewAuthor(firstName, middleName, lastName)
	err = Save(db, author1)
	iss.True(err != nil)
	iss.True(author1.ID == 0)
}

func TestGetByID(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	iss := is.New(t)
	expect := createNewAuthorInDB(db, t)
	got, err := repository.GetByID[models.Author](db, expect.ID)

	iss.NoErr(err)
	assertThatTheyAreSameAuthor(t, got, expect)
}

func Test_GetAll(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	iss := is.New(t)
	var expects []*models.Author
	expects = append(expects, createNewAuthorInDB(db, t))
	expects = append(expects, createNewAuthorInDB(db, t))
	expects = append(expects, createNewAuthorInDB(db, t))

	results, err := GetAll(db)

	iss.NoErr(err)
	iss.Equal(len(expects), 3)

	for _, v := range expects {
		assertThatContainsAuthor(t, results, v)
	}
}

func assertThatContainsAuthor(t *testing.T, results []*models.Author, expect *models.Author) {
	for _, v := range results {
		if v.ID == expect.ID {
			assertThatTheyAreSameAuthor(t, v, expect)
			return
		}
	}
}

func assertThatTheyAreSameAuthor(t *testing.T, get *models.Author, expect *models.Author) {
	iss := is.New(t)
	iss.Equal(get.FirstName, expect.FirstName)
	iss.Equal(get.MiddleName, expect.MiddleName)
	iss.Equal(get.LastName, expect.LastName)
	iss.Equal(get.CreatedAt.UTC(), expect.CreatedAt.UTC())
	iss.Equal(get.UpdatedAt.UTC(), expect.UpdatedAt.UTC())
	if expect.DeletedAt.Valid {
		iss.True(get.DeletedAt.Valid)
		iss.Equal(get.DeletedAt.Time.UTC(), expect.DeletedAt.Time.UTC())
	}
}

func Test_Delete(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	iss := is.New(t)

	a := createNewAuthorInDB(db, t)
	err := Delete(db, a.ID)

	iss.NoErr(err)
}

func createNewAuthorInDB(db *gorm.DB, t *testing.T) *models.Author {
	a := models.NewAuthor(
		random.RandomString(6),
		random.RandomString(6),
		random.RandomString(6))
	iss := is.New(t)
	err := Save(db, a)
	iss.NoErr(err)

	return a
}
