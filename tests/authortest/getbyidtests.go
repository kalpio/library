package authortest

import (
	"encoding/json"
	"fmt"
	"library/models"
	"library/random"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetExistingAuthorByID(t *testing.T) {
	ass := assert.New(t)

	author := createNewAuthor(ass)

	resp := requestGetData(author.ID)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusOK)
}

func GetNotExistingAuthorByID(t *testing.T) {
	ass := assert.New(t)

	resp := requestGetData(2137)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusBadRequest)
}

func createNewAuthor(ass *assert.Assertions) *models.Author {
	buff := prepareAuthorRequestData(random.RandomString(10), random.RandomString(10), random.RandomString(10))
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusCreated)

	var result *models.Author
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		ass.NoError(err)
	}

	return result
}

func requestGetData(id uint) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/author/%d", id), nil)
	return executeRequest(req)
}
