package main

import (
	"fmt"

	"./crawler"
)

func main() {
	var output = make(chan crawler.CCANItem, 25)
	go crawler.CrawlPage(output)

	var counter int
	for item := range output {
		fmt.Printf("%d: %#v\n\n", counter, item)
		counter++
	}

	// if err := zipfactory.CreateZipFileFromItems(output); err != nil {
	// 	log.Fatalln(err)
	// }
	println("Ende")
}
