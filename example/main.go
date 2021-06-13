package main

import (
	"flag"
	"fmt"
	"github.com/ktdf/mapbuilder"
)

func main() {
	site := flag.String("url", "https://calhoun.io", "site url which we are going to parse")
	depth := flag.Int("depth", 0, "How deep urls we are going to get. 0 - infinite" )
	flag.Parse()
	urls, err := mapbuilder.CollectUrls(*site, *depth)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", urls)
}
