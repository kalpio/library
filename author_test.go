package main

import (
	"bytes"
	"encoding/json"
	"library/random"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func addNewAuthor(t *testing.T) {
	iss := is.New(t)

	buff := prepareAuthorRequestData(random.RandomString(10), random.RandomString(10), random.RandomString(10))
	resp := postAuthorData(buff)

	iss.True(resp != nil)
	iss.Equal(resp.Code, http.StatusCreated)
}

func addDuplicatedAuthor(t *testing.T) {
	iss := is.New(t)

	buff := prepareAuthorRequestData(random.RandomString(10), random.RandomString(10), random.RandomString(10))
	resp0 := postAuthorData(buff)

	iss.True(resp0 != nil)
	iss.Equal(resp0.Code, http.StatusCreated)

	resp1 := postAuthorData(buff)

	iss.True(resp1 != nil)
	iss.Equal(resp1.Code, http.StatusBadRequest)
}

func creatingAuthorWithEmptyFirstNameShouldFail(t *testing.T) {
	iss := is.New(t)

	buff := prepareAuthorRequestData("", random.RandomString(10), random.RandomString(10))
	resp := postAuthorData(buff)

	iss.True(resp != nil)
	iss.Equal(resp.Code, http.StatusBadRequest)
}

func creatingAuthorWithEmptyLastNameShouldFail(t *testing.T) {
	iss := is.New(t)

	buff := prepareAuthorRequestData(random.RandomString(10), random.RandomString(10), "")
	resp := postAuthorData(buff)

	iss.True(resp != nil)
	iss.Equal(resp.Code, http.StatusBadRequest)
}

func creatingAuthorWithEmptyMiddleNameShouldPass(t *testing.T) {
	iss := is.New(t)

	buff := prepareAuthorRequestData(random.RandomString(10), "", random.RandomString(10))
	resp := postAuthorData(buff)

	iss.True(resp != nil)
	iss.Equal(resp.Code, http.StatusCreated)
}

func creatingAuthorWithEmptyPropsShouldFail(t *testing.T) {
	iss := is.New(t)

	buff := prepareAuthorRequestData("", "", "")
	resp := postAuthorData(buff)

	iss.True(resp != nil)
	iss.Equal(resp.Code, http.StatusBadRequest)
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
