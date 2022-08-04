package authortest

import (
	"bytes"
	"encoding/json"
	"library/application"
	"library/random"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var a application.App

func SetApp(app application.App) {
	a = app
}

func PostNewAuthor(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAuthorRequestData(random.RandomString(10), random.RandomString(10), random.RandomString(10))
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusCreated)
}

func PostDuplicatedAuthor(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAuthorRequestData(random.RandomString(10), random.RandomString(10), random.RandomString(10))
	resp0 := postAuthorData(buff)

	ass.NotNil(resp0)
	ass.Equal(resp0.Code, http.StatusCreated)

	resp1 := postAuthorData(buff)

	ass.NotNil(resp1)
	ass.Equal(resp1.Code, http.StatusBadRequest)
}

func PostAuthorWithEmptyFirstNameShouldFail(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAuthorRequestData("", random.RandomString(10), random.RandomString(10))
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusBadRequest)
}

func PostAuthorWithEmptyLastNameShouldFail(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAuthorRequestData(random.RandomString(10), random.RandomString(10), "")
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusBadRequest)
}

func PostAuthorWithEmptyMiddleNameShouldPass(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAuthorRequestData(random.RandomString(10), "", random.RandomString(10))
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusCreated)
}

func PostAuthorWithEmptyPropsShouldFail(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAuthorRequestData("", "", "")
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusBadRequest)
}

func prepareAuthorRequestData(firstName, middleName, lastName string) *bytes.Buffer {
	values := map[string]string{
		"first_name":  firstName,
		"middle_name": middleName,
		"last_name":   lastName,
	}
	jsonValue, _ := json.Marshal(values)

	return bytes.NewBuffer(jsonValue)
}

func postAuthorData(buff *bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", "/api/v1/author", buff)
	return executeRequest(req)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}
