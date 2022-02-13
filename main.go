package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"gohttp/request"
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

	hashReq := request.NewRequestor(md5.New())

	urlCh := make(chan string)
	resCh := make(chan request.HashResult)
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
		if res.Err != nil {
			fmt.Printf("%s : %v\n", res.RawURL, res.Err)
		} else {
			fmt.Printf("%s %s\n", res.RawURL, res.Hash)
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
