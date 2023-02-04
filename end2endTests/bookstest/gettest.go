package bookstest

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func GetExistingBookByID(t *testing.T) {
	ass := assert.New(t)

	bookDto, err := createBookInDB(a)
	ass.NoError(err)

	rr := getBookRR(a, bookDto.ID)
	ass.NotNil(rr)
	ass.Equal(http.StatusOK, rr.Code)

	bookResponse, err := getBookFromResponseBody(rr.Body)
	ass.NoError(err)
	assertCreatedBook(ass, *bookDto, bookResponse)
}

func GetNotExistingBookByID(t *testing.T) {
	ass := assert.New(t)

	resp := getBookRR(a, uuid.New().String())

	ass.NotNil(resp)
	ass.Equal(http.StatusBadRequest, resp.Code)
}
