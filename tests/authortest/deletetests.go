package authortest

import (
	"encoding/json"
	"fmt"
	"library/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func DeleteExistingAuthor(t *testing.T) {
	ass := assert.New(t)
	clearAuthorsTable(ass)

	var values []*domain.Author
	values = append(values, createNewAuthor(ass))
	values = append(values, createNewAuthor(ass))
	values = append(values, createNewAuthor(ass))

	resp := requestDelete(values[1].ID)
	ass.NotNil(resp)
	ass.Equal(http.StatusOK, resp.Code)
	valuesWithoutDeleted := []domain.Author{*values[0], *values[2]}

	respAuthors := requestGetAll()
	var result []domain.Author
	err := json.Unmarshal(respAuthors.Body.Bytes(), &result)
	ass.NoError(err)
	ass.Equal(len(valuesWithoutDeleted), len(result))
	ass.ElementsMatch(valuesWithoutDeleted, result)
}

func DeleteNotExistingAuthor(t *testing.T) {
	ass := assert.New(t)
	clearAuthorsTable(ass)

	var values []domain.Author
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))

	resp := requestDelete(uuid.New())
	ass.NotNil(resp)
	ass.Equal(http.StatusBadRequest, resp.Code)

	respAuthors := requestGetAll()
	var result []domain.Author
	err := json.Unmarshal(respAuthors.Body.Bytes(), &result)
	ass.NoError(err)
	ass.Equal(len(values), len(result))
	ass.ElementsMatch(values, result)
}

func DeletePermanentlyAuthor(t *testing.T) {
	ass := assert.New(t)
	clearAuthorsTable(ass)

	var values []domain.Author
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))

	notDeletedValues := []domain.Author{values[0], values[2]}

	resp := requestDeletePermanently(values[1].ID)
	ass.NotNil(resp)
	ass.Equal(http.StatusOK, resp.Code)

	responseAllAuthors := requestGetAll()
	var result []domain.Author
	err := json.Unmarshal(responseAllAuthors.Body.Bytes(), &result)
	ass.NoError(err)
	ass.Equal(len(notDeletedValues), len(result))
	ass.ElementsMatch(notDeletedValues, result)
}

func requestDelete(id uuid.UUID) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/author/%s", id.String()), nil)
	return executeRequest(req)
}

func requestDeletePermanently(id uuid.UUID) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/author/delete/%s", id.String()), nil)
	return executeRequest(req)
}
