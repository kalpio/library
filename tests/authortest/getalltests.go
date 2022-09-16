package authortest

import (
	"encoding/json"
	"library/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetExistingAuthors(t *testing.T) {
	ass := assert.New(t)
	clearAuthorsTable(ass)

	var values []domain.Author
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))
	values = append(values, *createNewAuthor(ass))

	resp := requestGetAll()
	var result []domain.Author
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	ass.NoError(err)
	ass.ElementsMatch(values, result)
}

func requestGetAll() *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/api/v1/author", nil)
	return executeRequest(req)
}
