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

	buff := prepareAddAuthorRequestData(random.RandomString(10), random.RandomString(10), random.RandomString(10))
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(http.StatusCreated, resp.Code)
}

func PostDuplicatedAuthor(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAddAuthorRequestData(random.RandomString(10), random.RandomString(10), random.RandomString(10))
	resp0 := postAuthorData(buff)

	ass.NotNil(resp0)
	ass.Equal(http.StatusCreated, resp0.Code)

	resp1 := postAuthorData(buff)

	ass.NotNil(resp1)
	ass.Equal(http.StatusBadRequest, resp1.Code)
}

func PostAuthorWithEmptyFirstNameShouldFail(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAddAuthorRequestData("", random.RandomString(10), random.RandomString(10))
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(http.StatusBadRequest, resp.Code)
}

func PostAuthorWithEmptyLastNameShouldFail(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAddAuthorRequestData(random.RandomString(10), random.RandomString(10), "")
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(http.StatusBadRequest, resp.Code)
}

func PostAuthorWithEmptyMiddleNameShouldPass(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAddAuthorRequestData(random.RandomString(10), "", random.RandomString(10))
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(http.StatusCreated, resp.Code)
}

func PostAuthorWithEmptyPropsShouldFail(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAddAuthorRequestData("", "", "")
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(http.StatusBadRequest, resp.Code)
}

func prepareAddAuthorRequestData(firstName, middleName, lastName string) *bytes.Buffer {
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
