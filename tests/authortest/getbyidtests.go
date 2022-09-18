package authortest

import (
	"fmt"
	"github.com/google/uuid"
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
	ass.Equal(http.StatusOK, resp.Code)
}

func GetNotExistingAuthorByID(t *testing.T) {
	ass := assert.New(t)

	resp := requestGetByID(uuid.New())

	ass.NotNil(resp)
	ass.Equal(http.StatusBadRequest, resp.Code)
}

func requestGetByID(id uuid.UUID) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/author/%s", id.String()), nil)
	return executeRequest(req)
}
