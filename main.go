package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"net/url"
	"sync"
)

var p *uint
var addresses []string

func main() {
	p = flag.Uint("parallel", 10, "number of parallel requests")
	flag.Parse()

	addresses = flag.Args()

	fmt.Println("Rate of parallelism:", *p)
	fmt.Println("Number of addresses: ", len(addresses))
	fmt.Println()

	hashReq := NewRequestor(md5.New())

	addrCh := make(chan string)
	resCh := make(chan requestResult)
	wg := &sync.WaitGroup{}

	hashReq.Process(*p, addrCh, resCh, wg)

	go func() {
		for _, addr := range addresses {
			addrCh <- validateSchema(addr)
		}

		close(addrCh)
		wg.Wait()

		close(resCh)
	}()

	for res := range resCh {
		if res.err != nil {
			fmt.Printf("%s : %v\n", res.address, res.err)
		} else {
			fmt.Printf("%s %s\n", res.address, res.hash)
		}
	}

}

func validateSchema(address string) string {
	url, err := url.ParseRequestURI(address)
	if err != nil || url.Scheme == "" {
		return "http://" + address
	}

	return address
}
