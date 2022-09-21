package authortest

import (
	"fmt"
	"github.com/google/uuid"
	"library/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetExistingAuthorByID(t *testing.T) {
	ass := assert.New(t)

	author, err := createNewAuthor()
	ass.NoError(err)

	resp := requestGetByID(author.ID)

	ass.NotNil(resp)
	ass.Equal(http.StatusOK, resp.Code)

	result, err := getAuthorFromResult(resp.Body)
	ass.NoError(err)
	assertAuthor(ass, author, result)
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

func assertAuthor(ass *assert.Assertions, expected, actual *domain.Author) {
	ass.Equal(expected.ID, actual.ID)
	ass.Equal(expected.FirstName, actual.FirstName)
	ass.Equal(expected.MiddleName, actual.MiddleName)
	ass.Equal(expected.LastName, actual.LastName)
	ass.Equal(expected.CreatedAt, actual.CreatedAt)
	ass.Equal(expected.UpdatedAt, actual.UpdatedAt)
	ass.Equal(expected.DeletedAt, actual.DeletedAt)
}
