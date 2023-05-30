package author

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"library/domain"
	"library/e2e/log"
	"library/e2e/utils"
	"library/random"
	"net/http"
	"time"
)

func Post(apiUrl string, count int, ch chan domain.AuthorID) {
	logger := log.NewLogger("POST /author")
	for i := 0; i < count; i++ {
		id := post(apiUrl, logger)
		if id.IsNil() {
			continue
		}
		ch <- id
	}
	close(ch)
}

func post(apiUrl string, logger *log.Logger) domain.AuthorID {
	url := fmt.Sprintf("%s/author", apiUrl)
	logger.Infolnf(url)

	values := prepareCreateAuthorData()
	jsonData, err := utils.MustMarshal(values)
	if err != nil {
		logger.Faillnf("failed to marshal values: %v", err)
		return domain.AuthorID(uuid.Nil.String())
	}

	client := &http.Client{}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Faillnf("failed to post: %v", err)
		return domain.AuthorID(uuid.Nil.String())
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			logger.Printlnf("failed to close response body: %v", errClose)
		}
	}()

	body, err := utils.GetBodyBytes(resp.Body)
	if err != nil {
		logger.Faillnf("failed to read response body: %v", err)
		return domain.AuthorID(uuid.Nil.String())
	}

	if resp.StatusCode != http.StatusCreated {
		logger.Printlnf("body: %s", string(body))
		logger.Faillnf("incorrect response status: expected %s, got: %s", http.StatusCreated, resp.StatusCode)
		return domain.AuthorID(uuid.Nil.String())
	}

	var response createAuthorResponse
	if err = json.Unmarshal(body, &response); err != nil {
		logger.Faillnf("failed to unmarshal response: %v", err)
		return domain.AuthorID(uuid.Nil.String())
	}

	logger.Printlnf("response: %+v", response)
	return response.ID
}

type createAuthorResponse struct {
	ID         domain.AuthorID `json:"id"`
	FirstName  string          `json:"first_name"`
	MiddleName string          `json:"middle_name"`
	LastName   string          `json:"last_name"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func prepareCreateAuthorData() map[string]interface{} {
	return map[string]interface{}{
		"first_name":  random.String(20),
		"middle_name": random.String(20),
		"last_name":   random.String(120),
	}
}
