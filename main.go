package main

import (
	"flag"
	"fmt"
)

var p *uint
var addresses []string

func main() {
	p = flag.Uint("parallel", 10, "number of parallel requests")
	flag.Parse()

	addresses = flag.Args()

	fmt.Println("Rate of parallelism:", *p)
	fmt.Println("Number of addresses: ", len(addresses))
	fmt.Println("Addresses:")
	for _, addr := range addresses {
		fmt.Println(addr)
	}
}
