package author

import (
	"gorm.io/gorm"
	"library/models"
	"library/random"
	"library/repository"
	"library/repository/testutils"
	"testing"

	"github.com/matryer/is"
)

func Test_SaveNewAuthor(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	firstName := random.RandomString(10)
	middleName := random.RandomString(10)
	lastName := random.RandomString(10)
	author0 := models.NewAuthor(firstName, middleName, lastName)

	iss := is.New(t)
	result, err := repository.Save(db, *author0)

	iss.NoErr(err)
	iss.True(result.ID > 0)
	iss.Equal(result.FirstName, firstName)
	iss.Equal(result.MiddleName, middleName)
	iss.Equal(result.LastName, lastName)
}

func Test_TryAddExistingAuthor(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	iss := is.New(t)
	firstName := random.RandomString(10)
	middleName := random.RandomString(10)
	lastName := random.RandomString(10)

	author0 := models.NewAuthor(firstName, middleName, lastName)
	result0, err := repository.Save(db, *author0)

	iss.NoErr(err)
	iss.Equal(result0.FirstName, firstName)
	iss.Equal(result0.MiddleName, middleName)
	iss.Equal(result0.LastName, lastName)

	author1 := models.NewAuthor(firstName, middleName, lastName)
	result1, err := repository.Save(db, *author1)
	iss.True(err != nil)
	iss.True(result1.ID == 0)
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
	var expects []models.Author
	expects = append(expects, createNewAuthorInDB(db, t))
	expects = append(expects, createNewAuthorInDB(db, t))
	expects = append(expects, createNewAuthorInDB(db, t))

	results, err := repository.GetAll[models.Author](db)

	iss.NoErr(err)
	iss.Equal(len(expects), 3)

	for _, v := range expects {
		assertThatContainsAuthor(t, results, v)
	}
}

func assertThatContainsAuthor(t *testing.T, results []models.Author, expect models.Author) {
	for _, v := range results {
		if v.ID == expect.ID {
			assertThatTheyAreSameAuthor(t, v, expect)
			return
		}
	}
}

func assertThatTheyAreSameAuthor(t *testing.T, got models.Author, expect models.Author) {
	iss := is.New(t)
	iss.Equal(got.FirstName, expect.FirstName)
	iss.Equal(got.MiddleName, expect.MiddleName)
	iss.Equal(got.LastName, expect.LastName)
	iss.Equal(got.CreatedAt.UTC(), expect.CreatedAt.UTC())
	iss.Equal(got.UpdatedAt.UTC(), expect.UpdatedAt.UTC())
	if expect.DeletedAt.Valid {
		iss.True(got.DeletedAt.Valid)
		iss.Equal(got.DeletedAt.Time.UTC(), expect.DeletedAt.Time.UTC())
	}
}

func Test_Delete(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	iss := is.New(t)

	a := createNewAuthorInDB(db, t)
	err := repository.Delete[models.Author](db, a.ID)

	iss.NoErr(err)
}

func createNewAuthorInDB(db *gorm.DB, t *testing.T) models.Author {
	a := models.NewAuthor(
		random.RandomString(6),
		random.RandomString(6),
		random.RandomString(6))
	iss := is.New(t)
	result, err := repository.Save(db, *a)
	iss.NoErr(err)

	return result
}
