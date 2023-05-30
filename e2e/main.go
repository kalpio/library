package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"library/domain"
	"library/e2e/author"
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

	var wg sync.WaitGroup
	var ch = make(chan domain.AuthorID, 1000)

	go author.Post(apiURL, cap(ch), ch)

	for c := range ch {
		wg.Add(2)
		author.Get(apiURL, c, &wg)
		author.Delete(apiURL, c, &wg)
	}

	wg.Wait()
}
