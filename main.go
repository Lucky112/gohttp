package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"net/url"
	"sync"
)

var p *uint
var urls []string

func main() {
	p = flag.Uint("parallel", 10, "number of parallel requests")
	flag.Parse()

	urls = flag.Args()

	fmt.Println("Rate of parallelism:", *p)
	fmt.Println("Number of URLs: ", len(urls))
	fmt.Println()

	hashReq := NewRequestor(md5.New())

	urlCh := make(chan string)
	resCh := make(chan requestResult)
	wg := &sync.WaitGroup{}

	hashReq.Process(*p, urlCh, resCh, wg)

	go func() {
		for _, rawURL := range urls {
			urlCh <- validateSchema(rawURL)
		}

		close(urlCh)
		wg.Wait()

		close(resCh)
	}()

	for res := range resCh {
		if res.err != nil {
			fmt.Printf("%s : %v\n", res.rawURL, res.err)
		} else {
			fmt.Printf("%s %s\n", res.rawURL, res.hash)
		}
	}

}

func validateSchema(rawURL string) string {
	url, err := url.ParseRequestURI(rawURL)
	if err != nil || url.Scheme == "" {
		return "http://" + rawURL
	}

	return rawURL
}
