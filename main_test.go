package main

import (
	"library/migrations"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/gorm"
)

var a App

func TestMain(m *testing.M) {
	dsn := "test.db"
	a.Host("127.0.0.1")
	a.Port("8089")
	a.Initialize(dsn)

	code := m.Run()

	dropDatabase(a.db, "test")

	os.Exit(code)
}

func dropDatabase(db *gorm.DB, dsn string) {
	if err := migrations.DropDatabase(db, dsn); err != nil {
		log.Fatalln(err)
	}
}

func TestAuthorAPI(t *testing.T) {
	t.Run("authorAddNewAuthor", addNewAuthor)
	t.Run("authorAddDuplicatedAuthor", addDuplicatedAuthor)
	t.Run("creatingAuthorWithEmptyFirstNameShouldFail", creatingAuthorWithEmptyFirstNameShouldFail)
	t.Run("creatingAuthorWithEmptyLastNameShouldFail", creatingAuthorWithEmptyLastNameShouldFail)
	t.Run("creatingAuthorWithEmptyMiddleNameShouldPass", creatingAuthorWithEmptyMiddleNameShouldPass)
	t.Run("creatingAuthorWithEmptyPropsShouldFail", creatingAuthorWithEmptyPropsShouldFail)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.router.ServeHTTP(rr, req)

	return rr
}
