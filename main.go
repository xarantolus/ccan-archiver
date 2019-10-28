package main

import (
	"fmt"
	"log"

	"github.com/xarantolus/ccan-archiver/crawler"
	"github.com/xarantolus/ccan-archiver/zipfactory"
)

func main() {
	var output = make(chan zipfactory.Archivable, 25)
	go func() {
		// // Crawl ccan.de and return its items
		fmt.Println("Downloading CCAN items")
		errorsCCAN := crawler.CrawlCCAN(output)
		fmt.Printf("There were %d errors while downloading from ccan.de: \n", len(errorsCCAN))
		for _, err := range errorsCCAN {
			fmt.Println(err.Error())
		}

		// Crawl an archive of clonk-center.net and return its items
		fmt.Println("Downloading Clonk-Center items")
		errorsCC := crawler.CrawlClonkCenter(output)
		fmt.Printf("There were %d errors while downloading from cc-archive.lwrl.de: \n", len(errorsCC))
		for _, err := range errorsCC {
			fmt.Println(err.Error())
		}

		close(output)
	}()

	if err := zipfactory.CreateZipFileFromItems(output); err != nil {
		log.Fatalln(err)
	}

	println("Finished downloading.")
}
