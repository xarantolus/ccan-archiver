package main

import (
	"fmt"
	"log"

	"./crawler"
	"./zipfactory"
)

func main() {
	var finished = make(chan bool)

	var output = make(chan zipfactory.Archivable, 25)
	go func() {
		defer close(output)

		// Crawl ccan.de and return its items
		errorsCCAN := crawler.CrawlCCAN(output)
		fmt.Printf("There were %d errors while downloading from ccan.de: \n", len(errorsCCAN))
		for _, err := range errorsCCAN {
			fmt.Println(err.Error())
		}

		finished <- true
	}()

	if err := zipfactory.CreateZipFileFromItems(output); err != nil {
		log.Fatalln(err)
	}

	<-finished // This should already have happened as zipfactory will need to process the outputs
	println("Finished downloading.")
}
