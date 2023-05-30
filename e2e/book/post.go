package book

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
)

func Post(apiUrl string, authors []domain.AuthorID, count int, ch chan domain.BookID) {
	logger := log.NewLogger("POST /book")
	for _, authorId := range authors {
		for i := 0; i < count; i++ {
			id := post(apiUrl, authorId, logger)
			if id.IsNil() {
				continue
			}
			ch <- id
		}
	}
	close(ch)
}

func post(apiUrl string, authorId domain.AuthorID, logger *log.Logger) domain.BookID {
	url := fmt.Sprintf("%s/book", apiUrl)
	logger.Infolnf(url)

	values := prepareCreateBookData(authorId)
	jsonData, err := utils.MustMarshal(values)
	if err != nil {
		logger.Faillnf("failed to marshal values: %v", err)
		return domain.BookID(uuid.Nil.String())
	}

	client := &http.Client{}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Faillnf("failed to post: %v", err)
		return domain.BookID(uuid.Nil.String())
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			logger.Printlnf("failed to close response body: %v", errClose)
		}
	}()

	body, err := utils.GetBodyBytes(resp.Body)
	if err != nil {
		logger.Faillnf("failed to read response body: %v", err)
		return domain.BookID(uuid.Nil.String())
	}

	if resp.StatusCode != http.StatusCreated {
		logger.Printlnf("body: %s", string(body))
		logger.Faillnf("incorrect response status: expected %s, got: %s", http.StatusCreated, resp.StatusCode)
		return domain.BookID(uuid.Nil.String())
	}

	var response createBookResponse
	if err := json.Unmarshal(body, &response); err != nil {
		logger.Faillnf("failed to unmarshal response: %v", err)
		return domain.BookID(uuid.Nil.String())
	}

	logger.Printlnf("response: %+v", response)
	return response.ID
}

type createBookResponse struct {
	ID          domain.BookID   `json:"id"`
	Title       string          `json:"title"`
	ISBN        string          `json:"isbn"`
	Description string          `json:"description"`
	AuthorID    domain.AuthorID `json:"author_id"`
}

func prepareCreateBookData(authorId domain.AuthorID) map[string]interface{} {
	return map[string]interface{}{
		"title":       random.String(20),
		"isbn":        random.String(13),
		"description": random.String(100),
		"author_id":   authorId.String(),
	}
}
