package authortest

import (
	"encoding/json"
	"library/domain"
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
	if err := a.DB().GetDB().Where("1 = 1").Delete(&domain.Author{}); err != nil {
		ass.NoError(err.Error)
	}
}

// TODO: shouldn't take param Assertions. In case of exception just return error
func createNewAuthor(ass *assert.Assertions) *domain.Author {
	buff := prepareAddAuthorRequestData(random.String(10), random.String(10), random.String(10))
	resp := postAuthorData(buff)

	ass.NotNil(resp)
	ass.Equal(resp.Code, http.StatusCreated)

	var result *domain.Author
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		ass.NoError(err)
	}

	return result
}
