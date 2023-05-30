package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"library/domain"
	"library/e2e/author"
	"library/e2e/book"
	"net"
	"sync"
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
	var authorChan = make(chan domain.AuthorID, 1000)
	const bookCountPerAuthor = 10
	var bookChan = make(chan domain.BookID, cap(authorChan)*bookCountPerAuthor)

	go author.Post(apiURL, cap(authorChan), authorChan)

	for c := range authorChan {
		authors = append(authors, c)
	}

	book.Post(apiURL, authors, bookCountPerAuthor, bookChan)

	for _, authorId := range authors {
		wg.Add(2)
		author.Get(apiURL, authorId, &wg)
		author.Delete(apiURL, authorId, &wg)
	}

	wg.Wait()
}
