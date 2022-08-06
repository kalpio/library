package author

import (
	"library/models"
	"library/random"
	"library/repository"
	"library/repository/testutils"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

func TestSaveNewAuthor(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	firstName := random.RandomString(10)
	middleName := random.RandomString(10)
	lastName := random.RandomString(10)
	author0 := models.NewAuthor(firstName, middleName, lastName)

	ass := assert.New(t)
	result, err := repository.Save(db, *author0)

	ass.NoError(err)
	ass.Greater(result.ID, uint(0))
	ass.Equal(result.FirstName, firstName)
	ass.Equal(result.MiddleName, middleName)
	ass.Equal(result.LastName, lastName)
}

func TestTryAddExistingAuthor(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	firstName := random.RandomString(10)
	middleName := random.RandomString(10)
	lastName := random.RandomString(10)

	author0 := models.NewAuthor(firstName, middleName, lastName)
	result0, err := repository.Save(db, *author0)

	ass.NoError(err)
	ass.Equal(result0.FirstName, firstName)
	ass.Equal(result0.MiddleName, middleName)
	ass.Equal(result0.LastName, lastName)

	author1 := models.NewAuthor(firstName, middleName, lastName)
	result1, err := repository.Save(db, *author1)
	ass.Error(err)
	ass.Equal(result1.ID, uint(0))
}

func TestTryAddEmptyFirstName(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	firstName := ""
	middleName := random.RandomString(10)
	lastName := random.RandomString(10)

	author := models.NewAuthor(firstName, middleName, lastName)
	result, err := repository.Save(db, *author)
	ass.Error(err)
	ass.Equal(result.ID, uint(0))
}

func TestTryAddEmptyLastName(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	firstName := random.RandomString(10)
	middleName := random.RandomString(10)
	lastName := ""

	author := models.NewAuthor(firstName, middleName, lastName)
	result, err := repository.Save(db, *author)
	ass.Error(err)
	ass.Equal(result.ID, uint(0))
}

func TestGetByID(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	expect := createNewAuthorInDB(db, t)
	got, err := repository.GetByID[models.Author](db, expect.ID)

	ass.NoError(err)
	assertThatTheyAreSameAuthor(ass, got, expect)
}

func TestGetAll(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	var expects []models.Author
	expects = append(expects, createNewAuthorInDB(db, t))
	expects = append(expects, createNewAuthorInDB(db, t))
	expects = append(expects, createNewAuthorInDB(db, t))

	results, err := repository.GetAll[models.Author](db)

	ass.NoError(err)
	ass.Equal(len(expects), 3)

	for _, v := range expects {
		assertThatContainsAuthor(ass, results, v)
	}
}

func assertThatContainsAuthor(ass *assert.Assertions, results []models.Author, expect models.Author) {
	for _, v := range results {
		if v.ID == expect.ID {
			assertThatTheyAreSameAuthor(ass, v, expect)
			return
		}
	}
}

func assertThatTheyAreSameAuthor(ass *assert.Assertions, got models.Author, expect models.Author) {
	ass.Equal(got.FirstName, expect.FirstName)
	ass.Equal(got.MiddleName, expect.MiddleName)
	ass.Equal(got.LastName, expect.LastName)
	ass.Equal(got.CreatedAt.UTC(), expect.CreatedAt.UTC())
	ass.Equal(got.UpdatedAt.UTC(), expect.UpdatedAt.UTC())

	if expect.DeletedAt.Valid {
		ass.True(got.DeletedAt.Valid)
		ass.Equal(got.DeletedAt.Time.UTC(), expect.DeletedAt.Time.UTC())
	}
}

func TestDelete(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	assert := assert.New(t)

	a := createNewAuthorInDB(db, t)
	err := repository.Delete[models.Author](db, a.ID)

	assert.NoError(err)
}

func createNewAuthorInDB(db *gorm.DB, t *testing.T) models.Author {
	a := models.NewAuthor(
		random.RandomString(6),
		random.RandomString(6),
		random.RandomString(6))
	assert := assert.New(t)
	result, err := repository.Save(db, *a)
	assert.NoError(err)

	return result
}
