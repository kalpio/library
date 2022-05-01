package author_test

import (
	"fmt"
	"library/migrations"
	"library/models"
	"library/random"
	"library/repository/author"
	"testing"

	"github.com/matryer/is"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var db *gorm.DB

func getDBName() string {
	return fmt.Sprintf("library_%s", random.RandomString(6))
}

func beforeTest(t *testing.T) func(t *testing.T) {
	dbName := getDBName()
	var err error
	iss := is.New(t)

	dsn := "sqlserver://jz:jzsoft@localhost:1433"
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	iss.NoErr(err)

	if err := migrations.CreateAndUseDatabase(db, dbName); err != nil {
		iss.NoErr(err)
	}

	dsn = fmt.Sprintf("sqlserver://jz:jzsoft@localhost:1433?database=%s", dbName)
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	iss.NoErr(err)

	if err := migrations.UpdateDatabase(db); err != nil {
		iss.NoErr(err)
	}

	return func(t *testing.T) {
		if err := migrations.DropDatabase(db, dbName); err != nil {
			iss.NoErr(err)
		}
	}
}

func Test_Save_New_Author(t *testing.T) {
	afterTest := beforeTest(t)
	defer afterTest(t)

	firstName := "Robert"
	middleName := "C."
	lastName := "Martin"
	author0 := models.NewAuthor(firstName, middleName, lastName)
	iss := is.New(t)

	result, err := author.Save(db, author0)

	iss.NoErr(err)
	iss.True(author0.ID > 0)
	iss.Equal(author0.FirstName, result.FirstName)
	iss.Equal(author0.MiddleName, result.MiddleName)
	iss.Equal(author0.LastName, result.LastName)
}

func Test_Try_Add_Existing_Author(t *testing.T) {
	afterTest := beforeTest(t)
	defer afterTest(t)

	iss := is.New(t)
	firstName := "Robert"
	middleName := "C."
	lastName := "Martin"

	author0 := models.NewAuthor(firstName, middleName, lastName)
	result0, err := author.Save(db, author0)
	iss.NoErr(err)
	iss.Equal(author0.FirstName, result0.FirstName)
	iss.Equal(author0.MiddleName, result0.MiddleName)
	iss.Equal(author0.LastName, result0.LastName)

	author1 := models.NewAuthor(firstName, middleName, lastName)
	result1, err := author.Save(db, author1)
	iss.True(err != nil)
	iss.Equal(author1.FirstName, result1.FirstName)
	iss.Equal(author1.MiddleName, result1.MiddleName)
	iss.Equal(author1.LastName, result1.LastName)
}

func Test_Get_Author(t *testing.T) {
	afterTest := beforeTest(t)
	defer afterTest(t)

	iss := is.New(t)
	expect := createNewAuthorInDB(t)
	get, err := author.GetByID(db, expect.ID)

	iss.NoErr(err)
	iss.True(expect.ID == get.ID)
	iss.True(expect.FirstName == get.FirstName)
	iss.True(expect.MiddleName == get.MiddleName)
	iss.True(expect.LastName == get.LastName)
}

func Test_GetAll(t *testing.T) {
	afterTest := beforeTest(t)
	defer afterTest(t)

	iss := is.New(t)
	expects := []*models.Author{}
	expects = append(expects, createNewAuthorInDB(t))
	expects = append(expects, createNewAuthorInDB(t))
	expects = append(expects, createNewAuthorInDB(t))

	results, err := author.GetAll(db)

	iss.NoErr(err)
	iss.Equal(len(expects), 3)

	for _, v := range expects {
		assert_that_contains_author(t, results, v)
	}
}

func assert_that_contains_author(t *testing.T, results []*models.Author, expect *models.Author) {
	iss := is.New(t)
	for _, v := range results {
		if v.ID == expect.ID {
			iss.Equal(v.FirstName, expect.FirstName)
			iss.Equal(v.MiddleName, expect.MiddleName)
			iss.Equal(v.LastName, expect.LastName)
			iss.Equal(v.CreatedAt, expect.CreatedAt)
			iss.Equal(v.UpdatedAt, expect.UpdatedAt)
			iss.Equal(v.DeletedAt, expect.DeletedAt)

			return
		}
	}

	iss.Fail()
}

func createNewAuthorInDB(t *testing.T) *models.Author {
	a := models.NewAuthor(random.RandomString(6), random.RandomString(6), random.RandomString(6))
	iss := is.New(t)
	result, err := author.Save(db, a)
	iss.NoErr(err)

	return result
}
