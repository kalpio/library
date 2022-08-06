package authortest

import (
	"encoding/json"
	"library/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetExistingAuthors(t *testing.T) {
	ass := assert.New(t)
	clearAuthorsTable(ass)

	values := []models.Author{}
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))

	resp := requestGetAll()
	var result []models.Author
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	ass.NoError(err)
	ass.ElementsMatch(values, result)
}

func requestGetAll() *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/api/v1/author", nil)
	return executeRequest(req)
}
