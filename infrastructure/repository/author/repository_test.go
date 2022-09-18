package author

import (
	"github.com/google/uuid"
	"library/domain"
	"library/infrastructure/repository"
	"library/infrastructure/repository/testutils"
	"library/random"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveNewAuthor(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	id := uuid.New()
	firstName := random.String(10)
	middleName := random.String(10)
	lastName := random.String(10)
	author0 := domain.NewAuthor(id, firstName, middleName, lastName)

	ass := assert.New(t)
	result, err := repository.Save(db, *author0)

	ass.NoError(err)
	ass.Equal(result.ID, id)
	ass.Equal(result.FirstName, firstName)
	ass.Equal(result.MiddleName, middleName)
	ass.Equal(result.LastName, lastName)
}

func TestTryAddExistingAuthor(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	id := uuid.New()
	firstName := random.String(10)
	middleName := random.String(10)
	lastName := random.String(10)

	author0 := domain.NewAuthor(id, firstName, middleName, lastName)
	result0, err := repository.Save(db, *author0)

	ass.NoError(err)
	ass.Equal(result0.ID, id)
	ass.Equal(result0.FirstName, firstName)
	ass.Equal(result0.MiddleName, middleName)
	ass.Equal(result0.LastName, lastName)

	author1 := domain.NewAuthor(id, firstName, middleName, lastName)
	result1, err := repository.Save(db, *author1)
	ass.Error(err)
	ass.Equal(result1.ID, id)
}

func TestTryAddEmptyFirstName(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	id := uuid.New()
	firstName := ""
	middleName := random.String(10)
	lastName := random.String(10)

	author := domain.NewAuthor(id, firstName, middleName, lastName)
	result, err := repository.Save(db, *author)
	ass.Error(err)
	ass.Equal(result.ID, domain.EmptyUUID())
}

func TestTryAddEmptyLastName(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	id := uuid.New()
	firstName := random.String(10)
	middleName := random.String(10)
	lastName := ""

	author := domain.NewAuthor(id, firstName, middleName, lastName)
	result, err := repository.Save(db, *author)
	ass.Error(err)
	ass.Equal(result.ID, domain.EmptyUUID())
}

func TestGetByID(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	expect := createNewAuthorInDB(db, t)
	got, err := repository.GetByID[domain.Author](db, expect.ID)

	ass.NoError(err)
	assertThatTheyAreSameAuthor(ass, got, expect)
}

func TestGetAll(t *testing.T) {
	db, afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	var expects []domain.Author
	expects = append(expects, createNewAuthorInDB(db, t))
	expects = append(expects, createNewAuthorInDB(db, t))
	expects = append(expects, createNewAuthorInDB(db, t))

	results, err := repository.GetAll[domain.Author](db)

	ass.NoError(err)
	ass.Equal(len(expects), 3)

	for _, v := range expects {
		assertThatContainsAuthor(ass, results, v)
	}
}

func assertThatContainsAuthor(ass *assert.Assertions, results []domain.Author, expect domain.Author) {
	for _, v := range results {
		if v.ID == expect.ID {
			assertThatTheyAreSameAuthor(ass, v, expect)
			return
		}
	}
}

func assertThatTheyAreSameAuthor(ass *assert.Assertions, got domain.Author, expect domain.Author) {
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
	ass := assert.New(t)

	a := createNewAuthorInDB(db, t)
	rowsAffected, err := repository.Delete[domain.Author](db, a.ID)

	ass.NoError(err)
	ass.Greater(rowsAffected, int64(0))
}

func createNewAuthorInDB(db domain.Database, t *testing.T) domain.Author {
	a := domain.NewAuthor(
		uuid.New(),
		random.String(6),
		random.String(6),
		random.String(6))
	ass := assert.New(t)
	result, err := repository.Save(db, *a)
	ass.NoError(err)

	return result
}
