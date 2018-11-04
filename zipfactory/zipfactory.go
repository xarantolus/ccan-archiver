package zipfactory

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"../crawler"
)

var (
	downloaded = make(map[string]bool)
)

// CreateZipFileFromItems streams the items in input to a zip file called 'result.zip'
func CreateZipFileFromItems(input chan crawler.CCANItem) error {
	// Create Zip
	f, err := os.Create("result.zip")
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	// This holds the current direct url to a file
	var currentDirectURL string

	// Create http client
	var client = http.Client{
		Timeout: 30 * time.Minute, // long timeout as downloads can be big
	}

	// If we get redirected, we set the direct url
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		currentDirectURL = req.URL.String()
		return nil
	}

	var itemCount int64 = 1

	// Loop over channel & Download & Pack
	for item := range input {
		// Check if we already this item or there is no download link
		if _, contains := downloaded[item.DownloadLink]; item.DownloadLink == "" || contains {
			println("Already have", item.DownloadLink)
			continue
		}
		// Set default value
		currentDirectURL = item.DownloadLink

		var resp, err = client.Get(item.DownloadLink)
		if err != nil {
			println("Error:", err.Error())
			continue
		}

		// Check if there was a redirect
		if currentDirectURL != "" {
			item.DirectLink = currentDirectURL
		}

		// Generate name and show user
		name := fmt.Sprintf("%s/%s.%s", cleanFilename(item.Author), cleanFilename(item.Name), getURLExtension(item.DirectLink))
		fmt.Printf("Downloading %s (#%d)", name, itemCount)

		// Create in zip file
		f, err := w.Create(name)
		if err != nil {
			println(" > Error while creating file:", err.Error())
			continue
		}

		// Copy to zip file
		if _, err = io.Copy(f, resp.Body); err != nil {
			println(" > Error while downloading:", err.Error())
			continue
		}

		if err = w.Flush(); err != nil {
			println(" > Error while flushing file:", err.Error())
			continue
		}

		// Write info json
		infoName := fmt.Sprintf("%s.json", name)
		result, err := json.MarshalIndent(item, "", "    ")
		if err != nil {
			println(" > Error while generating JSON:", err.Error())
			continue
		}

		fj, err := w.Create(infoName)
		if err != nil {
			println(" > Error while creating JSON file:", err.Error())
			continue
		}

		_, err = fj.Write(result)
		if err != nil {
			println(" > Error while writing JSON file:", err.Error())
			continue
		}

		if err = w.Flush(); err != nil {
			println(" > Error while flushing JSON file:", err.Error())
			continue
		}

		// add the current link to the links we already downloaded
		downloaded[item.DownloadLink] = true

		// Reset the url so the check above won't generate unexpected results
		currentDirectURL = ""

		println(" > Success")
		itemCount++
	}

	// Generate a README.md file
	rm, err := w.Create("README.md")
	if err != nil {
		return err
	}
	GenerateReadme(&rm, itemCount)
	println("\nGenerated README.")

	if err = w.Flush(); err != nil {
		return err
	}
	return nil
}

func getURLExtension(url string) string {
	split := strings.Split(url, ".")
	return split[len(split)-1]
}

var allowedChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789" + " -_.,()[]+" + "ÄÖÜßäöü"

func isAllowedChar(char rune) (contains bool) {
	for _, r := range allowedChars {
		if r == char {
			return true
		}
	}
	return false
}

func cleanFilename(in string) (out string) {
	var b = strings.Builder{}

	for _, item := range in {
		if isAllowedChar(item) {
			b.WriteRune(item)
		}
	}

	return b.String()
}
