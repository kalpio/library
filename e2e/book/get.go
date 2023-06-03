package book

import (
	"encoding/json"
	"fmt"
	"library/domain"
	"library/e2e/log"
	"library/e2e/utils"
	"net/http"
	"sync"
)

func Get(apiUrl string, id domain.BookID, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := log.NewLogger(fmt.Sprintf("GET /book/%s", id))

	url := fmt.Sprintf("%s/book/%s", apiUrl, id)
	logger.Printlnf(url)

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		logger.Faillnf("failed to get: %v", err)
		return
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			logger.Printlnf("failed to close response body: %v", errClose)
		}
	}()

	body, err := utils.GetBodyBytes(resp.Body)
	if err != nil {
		logger.Faillnf("failed to read response body: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		logger.Printlnf("body: %s", string(body))
		logger.Faillnf("incorrect response status: expected %s, got: %s", http.StatusOK, resp.StatusCode)
		return
	}

	var response createBookResponse
	if err = json.Unmarshal(body, &response); err != nil {
		logger.Faillnf("failed to unmarshal response: %v", err)
		return
	}

	logger.Printlnf("response: %+v", response)
}
