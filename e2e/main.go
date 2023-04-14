package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"library/domain"
	"library/random"
	"net"
	"net/http"
	"sync"
	"time"
)

func main() {
	var host string
	var port string
	flag.StringVar(&host, "host", "localhost", "host")
	flag.StringVar(&port, "port", "8089", "port")

	flag.Parse()
	baseURL := fmt.Sprintf("http://%s", net.JoinHostPort(host, port))
	log.Println(fmt.Sprintf("baseURL: %s", baseURL))

	apiURL := fmt.Sprintf("%s/api/v1", baseURL)

	var authors []domain.AuthorID
	var wg sync.WaitGroup
	var ch = make(chan domain.AuthorID, 1000)
	go postAuthors(apiURL, cap(ch), ch)
	for c := range ch {
		authors = append(authors, c)
	}

	for _, author := range authors {
		wg.Add(1)
		go getAuthor(apiURL, author, &wg)
	}

	wg.Wait()
}

func postAuthors(apiURL string, count int, ch chan domain.AuthorID) {
	for i := 0; i < count; i++ {
		ch <- postAuthor(apiURL)
	}
	close(ch)
}

func postAuthor(apiURL string) domain.AuthorID {
	url := fmt.Sprintf("%s/author", apiURL)
	log.Println(fmt.Sprintf("POST %q", url))

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
		return domain.AuthorID(domain.EmptyUUID().String())
	}

	var response createAuthorResponse
	if err = json.Unmarshal(body, &response); err != nil {
		fail("author [post]: failed to unmarshal response: %v", err)
		return domain.AuthorID(domain.EmptyUUID().String())
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

func getAuthor(apiURL string, id domain.AuthorID, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("%s/author/%s", apiURL, id)
	log.Println(fmt.Sprintf("GET %q", url))

	client := &http.Client{}
	resp, err := client.Get(url)
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			log.Println(fmt.Sprintf("failed to close response body: %v", err))
		}
	}()

	body := getBodyBytes(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Println(fmt.Sprintf("body: %s", string(body)))
		fail("author [get]: incorrect response status: expected %s, got: %s", http.StatusOK, resp.StatusCode)
		return
	}

	var response createAuthorResponse
	if err = json.Unmarshal(body, &response); err != nil {
		fail("author [get]: failed to unmarshal response: %v", err)
		return
	}

	log.Println(fmt.Sprintf("response: %+v", response))
}

func fail(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	log.Println(fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05.000"), message))
}
