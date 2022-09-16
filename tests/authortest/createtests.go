package authortest

import (
	"bytes"
	"encoding/json"
	"library/application"
	"library/domain"
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
	values := map[string]string{
		"first_name":  random.String(10),
		"middle_name": random.String(10),
		"last_name":   random.String(10),
	}

	buff := prepareAddAuthorRequestData(values["first_name"], values["middle_name"], values["last_name"])
	resp := postAuthorData(buff)
	ass.NotNil(resp)
	ass.Equal(http.StatusCreated, resp.Code)

	assertCreatedAuthor(ass, values, resp.Body)
}

func PostDuplicatedAuthor(t *testing.T) {
	ass := assert.New(t)
	values := map[string]string{
		"first_name":  random.String(10),
		"middle_name": random.String(10),
		"last_name":   random.String(10),
	}

	buff := prepareAddAuthorRequestData(values["first_name"], values["middle_name"], values["last_name"])
	resp0 := postAuthorData(buff)

	ass.NotNil(resp0)
	ass.Equal(http.StatusCreated, resp0.Code)

	assertCreatedAuthor(ass, values, resp0.Body)

	resp1 := postAuthorData(buff)

	ass.NotNil(resp1)
	ass.Equal(http.StatusBadRequest, resp1.Code)
}

func PostAuthorWithEmptyFirstNameShouldFail(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAddAuthorRequestData("", random.String(10), random.String(10))
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(http.StatusBadRequest, resp.Code)
}

func PostAuthorWithEmptyLastNameShouldFail(t *testing.T) {
	ass := assert.New(t)

	buff := prepareAddAuthorRequestData(random.String(10), random.String(10), "")
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(http.StatusBadRequest, resp.Code)
}

func PostAuthorWithEmptyMiddleNameShouldPass(t *testing.T) {
	ass := assert.New(t)
	values := map[string]string{
		"first_name":  random.String(10),
		"middle_name": "",
		"last_name":   random.String(10),
	}

	buff := prepareAddAuthorRequestData(values["first_name"], values["middle_name"], values["last_name"])
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(http.StatusCreated, resp.Code)

	assertCreatedAuthor(ass, values, resp.Body)
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

func getAuthorFromResult(body *bytes.Buffer) (*domain.Author, error) {
	var result *domain.Author

	if err := json.Unmarshal(body.Bytes(), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func assertCreatedAuthor(ass *assert.Assertions, values map[string]string, body *bytes.Buffer) {
	createdAuthor, err := getAuthorFromResult(body)
	ass.NoError(err)
	ass.Equal(values["first_name"], createdAuthor.FirstName)
	ass.Equal(values["middle_name"], createdAuthor.MiddleName)
	ass.Equal(values["last_name"], createdAuthor.LastName)
}
