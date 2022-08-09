package authortest

import (
	"encoding/json"
	"fmt"
	"library/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func DeleteExistingAuthor(t *testing.T) {
	ass := assert.New(t)
	clearAuthorsTable(ass)

	var values []models.Author
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))

	resp := requestDelete(values[1].ID)
	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusOK)
	valuesWithoutDeleted := []models.Author{values[0], values[2]}

	respAuthors := requestGetAll()
	var result []models.Author
	err := json.Unmarshal(respAuthors.Body.Bytes(), &result)
	ass.NoError(err)
	ass.ElementsMatch(valuesWithoutDeleted, result)
}

func DeleteNotExistingAuthor(t *testing.T) {
	ass := assert.New(t)
	clearAuthorsTable(ass)

	var values []models.Author
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))

	resp := requestDelete(values[2].ID + 2137)
	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusBadRequest)

	respAuthors := requestGetAll()
	var result []models.Author
	err := json.Unmarshal(respAuthors.Body.Bytes(), &result)
	ass.NoError(err)
	ass.ElementsMatch(values, result)
}

func requestDelete(id uint) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/author/%d", id), nil)
	return executeRequest(req)
}
