package author

import (
	"library/domain"
	"library/infrastructure/repository"
	"library/infrastructure/repository/testutils"
	"library/random"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func Test_SaveAuthorSucceeded(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	id := uuid.New()
	firstName := random.String(10)
	middleName := random.String(10)
	lastName := random.String(10)
	author0 := domain.NewAuthor(id, firstName, middleName, lastName)

	ass := assert.New(t)
	result, err := repository.Save(*author0)

	ass.NoError(err)
	ass.Equal(result.ID, id)
	ass.Equal(result.FirstName, firstName)
	ass.Equal(result.MiddleName, middleName)
	ass.Equal(result.LastName, lastName)
}

func Test_SaveReturnsError_When_AuthorAlreadyExists(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	id := uuid.New()
	firstName := random.String(10)
	middleName := random.String(10)
	lastName := random.String(10)

	author0 := domain.NewAuthor(id, firstName, middleName, lastName)
	result0, err := repository.Save(*author0)

	ass.NoError(err)
	ass.Equal(result0.ID, id)
	ass.Equal(result0.FirstName, firstName)
	ass.Equal(result0.MiddleName, middleName)
	ass.Equal(result0.LastName, lastName)

	author1 := domain.NewAuthor(id, firstName, middleName, lastName)
	result1, err := repository.Save(*author1)
	ass.Error(err)
	ass.Equal(result1.ID, id)
}

func Test_SaveReturnsError_When_FirstNameIsEmpty(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	id := uuid.New()
	firstName := ""
	middleName := random.String(10)
	lastName := random.String(10)

	author := domain.NewAuthor(id, firstName, middleName, lastName)
	result, err := repository.Save(*author)
	ass.Error(err)
	ass.Equal(result.ID, domain.EmptyUUID())
}

func Test_SaveReturnsError_When_LastNameIsEmpty(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	id := uuid.New()
	firstName := random.String(10)
	middleName := random.String(10)
	lastName := ""

	author := domain.NewAuthor(id, firstName, middleName, lastName)
	result, err := repository.Save(*author)
	ass.Error(err)
	ass.Equal(result.ID, domain.EmptyUUID())
}

func TestGetByID(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	expect := createNewAuthorInDB(t)
	got, err := repository.GetByID[domain.Author](expect.ID)

	ass.NoError(err)
	assertThatTheyAreSameAuthor(ass, got, expect)
}

func TestGetAll(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)

	ass := assert.New(t)
	var expects []domain.Author
	expects = append(expects, createNewAuthorInDB(t))
	expects = append(expects, createNewAuthorInDB(t))
	expects = append(expects, createNewAuthorInDB(t))

	results, err := repository.GetAll[domain.Author]()

	ass.NoError(err)
	ass.Equal(3, len(results))

	for _, v := range expects {
		assertThatContainsAuthor(ass, results, v)
	}
}

func TestDelete(t *testing.T) {
	afterTest := testutils.BeforeTest(t)
	defer afterTest(t)
	ass := assert.New(t)

	a := createNewAuthorInDB(t)
	err := repository.Delete[domain.Author](a.ID)

	ass.NoError(err)
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
}

func createNewAuthorInDB(t *testing.T) domain.Author {
	a := domain.NewAuthor(
		uuid.New(),
		random.String(6),
		random.String(6),
		random.String(6))
	ass := assert.New(t)
	result, err := repository.Save(*a)
	ass.NoError(err)

	return result
}
