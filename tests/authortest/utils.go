package authortest

import (
	"encoding/json"
	"errors"
	"fmt"
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

func createNewAuthor() (*domain.Author, error) {
	buff := prepareAddAuthorRequestData(random.String(10), random.String(10), random.String(10))
	resp := postAuthorData(buff)

	if resp == nil {
		return nil, errors.New("createNewAuthor: POST response is nil")
	}

	if resp.Code != http.StatusCreated {
		return nil, fmt.Errorf("createNewAuthor: POST response status code is %q", resp.Code)
	}

	var result *domain.Author
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("createNewAuthor: deserialize response body: %w", err)
	}

	return result, nil
}
