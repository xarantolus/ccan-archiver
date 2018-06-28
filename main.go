package main

import (
	"log"

	"./crawler"
	"./zipfactory"
)

func main() {
	// TODO: follow external links
	var output = make(chan crawler.CCANItem, 25)
	go crawler.CrawlPage(output)

	if err := zipfactory.CreateZipFileFromItems(output); err != nil {
		log.Fatalln(err)
	}
	println("Ende")
}
