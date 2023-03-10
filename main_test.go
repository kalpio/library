package main_test

import (
	"library/migrations"
	"log"
	"os"
	"testing"

	"library/application"
	"library/end2endTests/authortest"
	"library/end2endTests/bookstest"
)

var a application.App

func TestMain(m *testing.M) {
	a.Host("127.0.0.1")
	a.Port("8089")
	a.Initialize()

	code := m.Run()

	dropDatabase("test")

	os.Exit(code)
}

func dropDatabase(dsn string) {
	if err := migrations.DropDatabase(dsn); err != nil {
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

	t.Run("EditExistingAuthor", authortest.EditExistingAuthor)

	bookstest.SetApp(a)
	t.Run("POST_NewBook", bookstest.PostNewBook)

	t.Run("GetExistingBookByID", bookstest.GetExistingBookByID)
	t.Run("GetNotExistingBookByID", bookstest.GetNotExistingBookByID)
}
