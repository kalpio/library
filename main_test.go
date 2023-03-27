package main_test

import (
	"library/ioc"
	"library/migrations"
	"log"
	"os"
	"testing"

	"library/application"
	"library/end2endTests/authortest"
	"library/end2endTests/bookstest"
)

var testApplication application.App

func TestMain(m *testing.M) {
	testApplication.Host("127.0.0.1")
	testApplication.Port("8089")
	testApplication.Initialize()

	code := m.Run()

	dropDatabase()

	os.Exit(code)
}

func dropDatabase() {
	migration, err := ioc.Get[migrations.Migration]()
	if err != nil {
		log.Fatalln(err)
	}
	if err = migration.DropDatabase(); err != nil {
		log.Fatalln(err)
	}
}

func TestAuthorAPI(t *testing.T) {
	authortest.SetApp(testApplication)
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

	t.Run("EditExistingAuthor", authortest.EditExistingAuthor)

	bookstest.SetApp(testApplication)
	t.Run("POST_NewBook", bookstest.PostNewBook)

	t.Run("GetExistingBookByID", bookstest.GetExistingBookByID)
	t.Run("GetNotExistingBookByID", bookstest.GetNotExistingBookByID)
}
