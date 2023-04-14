package author

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"library/domain"
	"library/random"
	"net/http"
	"time"
)

func Post(apiUrl string, count int, ch chan domain.AuthorID) {

}

func post(apiUrl string) map[string]interface{} {
	url := fmt.Sprintf("%s/author", apiUrl)
	log.Infof("POST %q\n", url)

	values := prepareCreateAuthorData()
	jsonData := mustMarshal(values)

	client := &http.Client{}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			log.Println(fmt.Sprintf("failed to close response body: %v", err))
		}
	}()

	body := getBodyBytes(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		log.Println(fmt.Sprintf("body: %s", string(body)))
		fail("author [post]: incorrect response status: expected %s, got: %s", http.StatusCreated, resp.StatusCode)
		return nil
	}

	var response createAuthorResponse
	if err = json.Unmarshal(body, &response); err != nil {
		fail("author [post]: failed to unmarshal response: %v", err)
		return nil
	}

	log.Println(fmt.Sprintf("response: %+v", response))
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
func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		fail("failed to marshal: %v", err)
	}
	return b
}

func getBodyBytes(body io.Reader) []byte {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		fail("failed to read body: %v", err)
	}
	return buf.Bytes()
}

func fail(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Fatalf("POST /author: %s", message)
}
