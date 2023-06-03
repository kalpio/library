package book

import (
	"encoding/json"
	"fmt"
	"library/domain"
	"library/e2e/log"
	"library/e2e/utils"
	"net/http"
)

func GetAll(apiUrl string) []*domain.Book {
	logger := log.NewLogger("GET /book")

	url := fmt.Sprintf("%s/book", apiUrl)
	logger.Printlnf(url)

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		logger.Faillnf("failed to get: %v", err)
		return nil
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			logger.Printlnf("failed to close response body: %v", errClose)
		}
	}()

	body, err := utils.GetBodyBytes(resp.Body)
	if err != nil {
		logger.Faillnf("failed to read response body: %v", err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		logger.Printlnf("body: %s", string(body))
		logger.Faillnf("incorrect response status: expected %s, got: %s", http.StatusOK, resp.StatusCode)
		return nil
	}

	var response getAllBooksResponse
	if err = json.Unmarshal(body, &response); err != nil {
		logger.Faillnf("failed to unmarshal response: %v", err)
		return nil
	}

	logger.Printlnf("response: found %+v books", len(response.Books))

	return response.Books
}

type getAllBooksResponse struct {
	Books []*domain.Book `json:"books"`
}
