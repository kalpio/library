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

	var (
		values []domain.Author
		model  *domain.Author
		err    error
	)
	for i := 0; i < 3; i++ {
		if model, err = createNewAuthor(); err != nil {
			ass.FailNow(err.Error())
		}
		values = append(values, *model)
	}

	resp := requestGetAll()
	var result []domain.Author
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	ass.NoError(err)
	ass.Equal(len(values), len(result))
	ass.ElementsMatch(values, result)
}

func requestGetAll() *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/api/v1/author", nil)
	return executeRequest(req)
}
