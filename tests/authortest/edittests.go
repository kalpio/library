package authortest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"library/random"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func EditExistingAuthor(t *testing.T) {
	ass := assert.New(t)

	author := createNewAuthor(ass)
	firstName, middleName, lastname := random.RandomString(10), random.RandomString(10), random.RandomString(10)

	buff := prepareEditAuthorRequestData(author.ID, firstName, middleName, lastname)
	resp := patchAuthorData(buff, author.ID)

	ass.NotNil(resp)
	ass.Equal(http.StatusOK, resp.Code)

	values, err := unmarshalEditResponse(resp.Body)
	ass.NoError(err)

	ass.Equal(author.ID, uint(values["id"].(float64)))
	ass.Equal(firstName, values["first_name"])
	ass.Equal(middleName, values["middle_name"])
	ass.Equal(lastname, values["last_name"])
	ass.Equal(author.CreatedAt.Format(time.RFC3339Nano), values["created_at"])

	updatedAt, err := time.Parse(time.RFC3339Nano, values["updated_at"].(string))
	ass.NoError(err)
	ass.Greater(updatedAt, author.CreatedAt)
}

func unmarshalEditResponse(body *bytes.Buffer) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal(body.Bytes(), &result)

	return result, err
}

func prepareEditAuthorRequestData(id uint, firstName, middleName, lastName string) *bytes.Buffer {
	values := map[string]any{
		"id":          id,
		"first_name":  firstName,
		"middle_name": middleName,
		"last_name":   lastName,
	}
	jsonValue, _ := json.Marshal(values)

	return bytes.NewBuffer(jsonValue)
}

func patchAuthorData(buff *bytes.Buffer, id uint) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/api/v1/author/%d", id), buff)
	return executeRequest(req)
}
