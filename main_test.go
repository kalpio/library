package main

import (
	"library/migrations"
	"log"
	"os"
	"testing"

	"library/application"
	"library/tests/authortest"

	"gorm.io/gorm"
)

var a application.App

func TestMain(m *testing.M) {
	dsn := "test.db"
	a.Host("127.0.0.1")
	a.Port("8089")
	a.Initialize(dsn)

	code := m.Run()

	dropDatabase(a.DB, "test")

	os.Exit(code)
}

func dropDatabase(db *gorm.DB, dsn string) {
	if err := migrations.DropDatabase(db, dsn); err != nil {
		log.Fatalln(err)
	}
}

func TestAuthorAPI(t *testing.T) {
	authortest.SetApp(a)
	t.Run("authorAddNewAuthor", authortest.PostNewAuthor)
	t.Run("authorAddDuplicatedAuthor", authortest.PostDuplicatedAuthor)
	t.Run("creatingAuthorWithEmptyFirstNameShouldFail", authortest.PostAuthorWithEmptyFirstNameShouldFail)
	t.Run("creatingAuthorWithEmptyLastNameShouldFail", authortest.PostAuthorWithEmptyLastNameShouldFail)
	t.Run("creatingAuthorWithEmptyMiddleNameShouldPass", authortest.PostAuthorWithEmptyMiddleNameShouldPass)
	t.Run("creatingAuthorWithEmptyPropsShouldFail", authortest.PostAuthorWithEmptyPropsShouldFail)
}
