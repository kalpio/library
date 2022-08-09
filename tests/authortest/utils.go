package authortest

import (
	"encoding/json"
	"library/models"
	"library/random"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router().ServeHTTP(rr, req)

	return rr
}

func clearAuthorsTable(ass *assert.Assertions) {
	if err := a.DB().Where("1 = 1").Delete(&models.Author{}); err != nil {
		ass.NoError(err.Error)
	}
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
