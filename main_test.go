package main

import (
	"library/domain"
	"library/migrations"
	"log"
	"os"
	"testing"

	"library/application"
	"library/tests/authortest"
)

var a application.App

func TestMain(m *testing.M) {
	dsn := "test.db"
	a.Host("127.0.0.1")
	a.Port("8089")
	a.Initialize(dsn)

	code := m.Run()

	dropDatabase(a.DB(), "test")

	os.Exit(code)
}

func dropDatabase(db domain.Database, dsn string) {
	if err := migrations.DropDatabase(db, dsn); err != nil {
		log.Fatalln(err)
	}
}

func TestAuthorAPI(t *testing.T) {
	authortest.SetApp(a)
	t.Run("PostNewAuthor", authortest.PostNewAuthor)
	t.Run("PostDuplicatedAuthor", authortest.PostDuplicatedAuthor)
	t.Run("PostAuthorWithEmptyFirstNameShouldFail", authortest.PostAuthorWithEmptyFirstNameShouldFail)
	t.Run("PostAuthorWithEmptyLastNameShouldFail", authortest.PostAuthorWithEmptyLastNameShouldFail)
	t.Run("PostAuthorWithEmptyMiddleNameShouldPass", authortest.PostAuthorWithEmptyMiddleNameShouldPass)
	t.Run("PostAuthorWithEmptyPropsShouldFail", authortest.PostAuthorWithEmptyPropsShouldFail)

	t.Run("GetExistingAuthorByID", authortest.GetExistingAuthorByID)
	t.Run("GetNotExistingAuthorByID", authortest.GetNotExistingAuthorByID)

	t.Run("GetExistingAuthors", authortest.GetExistingAuthors)

	t.Run("DeleteExistingAuthor", authortest.DeleteExistingAuthor)
	t.Run("DeleteNotExistingAuthor", authortest.DeleteNotExistingAuthor)
}
