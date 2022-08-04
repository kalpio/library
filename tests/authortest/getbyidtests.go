package authortest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetExistingAuthorByID(t *testing.T) {
	ass := assert.New(t)

	author := createNewAuthor(ass)

	resp := requestGetByID(author.ID)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusOK)
}

func GetNotExistingAuthorByID(t *testing.T) {
	ass := assert.New(t)

	resp := requestGetByID(2137)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusBadRequest)
}

func requestGetByID(id uint) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/author/%d", id), nil)
	return executeRequest(req)
}
