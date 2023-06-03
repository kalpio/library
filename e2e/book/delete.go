package book

import (
	"fmt"
	"library/domain"
	"library/e2e/log"
	"net/http"
	"sync"
)

func Delete(apiUrl string, id domain.BookID, wg *sync.WaitGroup) {
	defer wg.Done()

	logger := log.NewLogger(fmt.Sprintf("DELETE /book/%s", id))
	url := fmt.Sprintf("%s/book/%s", apiUrl, id)
	logger.Printlnf(url)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		logger.Faillnf("failed to create request: %v", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Faillnf("failed to delete: %v", err)
		return
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			logger.Printlnf("failed to close response body: %v", errClose)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logger.Faillnf("incorrect response status: expected %s, got: %s", http.StatusOK, resp.StatusCode)
		return
	}

	logger.Printlnf("response: %+v", resp)
}
