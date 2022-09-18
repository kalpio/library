package authortest

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"library/domain"
	"net/http"
	"net/http/httptest"
	"testing"

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

func requestDelete(id uuid.UUID) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/author/%s", id.String()), nil)
	return executeRequest(req)
}
