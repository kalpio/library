package authortest

import (
	"encoding/json"
	"fmt"
	"library/models"
	"library/random"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func GetExistingAuthorByID(t *testing.T) {
	iss := is.New(t)

	author := createNewAuthor(iss)

	resp := requestGetData(author.ID)
	iss.True(resp != nil)
	iss.Equal(resp.Code, http.StatusOK)
}

func GetNotExistingAuthorByID(t *testing.T) {
	iss := is.New(t)

	resp := requestGetData(2137)
	iss.True(resp != nil)
	iss.Equal(resp.Code, http.StatusBadRequest)
}

func createNewAuthor(iss *is.I) *models.Author {
	buff := prepareAuthorRequestData(random.RandomString(10), random.RandomString(10), random.RandomString(10))
	resp := postAuthorData(buff)
	iss.True(resp != nil)
	iss.Equal(resp.Code, http.StatusCreated)

	var result *models.Author
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		iss.NoErr(err)
	}

	return result
}

func requestGetData(id uint) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/author/%d", id), nil)
	return executeRequest(req)
}
