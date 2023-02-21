package bookstest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"library/application"
	"library/domain"
	"library/random"
	"net/http"
	"net/http/httptest"
)

func executeRequest(app application.App, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router().ServeHTTP(rr, req)

	return rr
}

type postBookDto struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	ISBN        domain.ISBN `json:"isbn"`
	Description string      `json:"description"`
	AuthorID    string      `json:"author_id"`
}

func generateBookDto(bookAuthorID uuid.UUID) postBookDto {
	return postBookDto{
		ID:          uuid.New().String(),
		Title:       random.String(20),
		ISBN:        domain.ISBN(random.String(13)),
		Description: random.String(20),
		AuthorID:    bookAuthorID.String(),
	}
}

func postBookRR(app application.App, dto postBookDto) *httptest.ResponseRecorder {
	jsonValue, _ := json.Marshal(dto)
	buff := bytes.NewBuffer(jsonValue)
	req, _ := http.NewRequest("POST", "/api/v1/book", buff)

	return executeRequest(app, req)
}

func getBookRR(app application.App, bookID string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/book/%s", bookID), nil)

	return executeRequest(app, req)
}

func getBookFromResponseBody(body *bytes.Buffer) (*domain.Book, error) {
	var result *domain.Book
	if err := json.Unmarshal(body.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("utils_test: get book from response body: %w", err)
	}

	return result, nil
}

type postAuthorDto struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

func postAuthor(app application.App, dto postAuthorDto) (*domain.Author, error) {
	jsonValue, _ := json.Marshal(dto)
	buff := bytes.NewBuffer(jsonValue)
	req, _ := http.NewRequest("POST", "/api/v1/author", buff)
	rr := executeRequest(app, req)
	if rr.Code != http.StatusCreated {
		return nil, errors.New("utils_test: error during POST author")
	}

	var result *domain.Author
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("utils_test: deserialize response body: %w", err)
	}

	return result, nil
}

func createBookInDB(app application.App) (*postBookDto, error) {
	bookAuthor, err := postAuthor(app, postAuthorDto{
		FirstName:  random.String(20),
		MiddleName: random.String(20),
		LastName:   random.String(20),
	})

	if err != nil {
		return nil, err
	}

	bookDto := generateBookDto(bookAuthor.ID)
	rr := postBookRR(app, bookDto)

	if rr == nil || rr.Code != http.StatusCreated {
		return nil, errors.New("test utils: cannot create book")
	}

	return &bookDto, nil
}
