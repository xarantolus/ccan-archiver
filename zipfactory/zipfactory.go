package zipfactory

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type Archivable interface {
	GetDownloadLink() string
	GetAuthor() string
	GetName() string
	GetSourceName() string
}

type createError struct {
	Error    string     `json:"error_message"`
	Metadata Archivable `json:"item"`
}

var (
	downloaded   = make(map[string]bool)
	failedEntrys = []createError{}
)

func appendPrintError(what string, err error, item Archivable) {
	errMessage := fmt.Sprintf("%s: %s", what, err.Error())
	fmt.Printf(" > Error %s\n", errMessage)
	failedEntrys = append(failedEntrys, createError{
		Error:    errMessage,
		Metadata: item,
	})
}

// CreateZipFileFromItems streams the items in input to a zip file called 'result.zip'
func CreateZipFileFromItems(input chan Archivable) error {
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
		if _, contains := downloaded[item.GetDownloadLink()]; item.GetDownloadLink() == "" || contains {
			println("Already have", item.GetDownloadLink)
			continue
		}
		// Set default value
		currentDirectURL = item.GetDownloadLink()

		var resp, err = client.Get(item.GetDownloadLink())
		if err != nil {
			appendPrintError("while downloading item", err, item)
			continue
		}

		// Generate name and show user
		name := fmt.Sprintf("%s/%s/%s.%s", item.GetSourceName(), cleanFilename(item.GetAuthor()), cleanFilename(item.GetName()), getURLExtension(currentDirectURL))
		fmt.Printf("Downloading %s (#%d)", name, itemCount)

		// Create in zip file
		f, err := w.Create(name)
		if err != nil {
			appendPrintError("while creating file", err, item)
			continue
		}

		// Copy to zip file
		if _, err = io.Copy(f, resp.Body); err != nil {
			appendPrintError("while copying file stream to archive", err, item)
			continue
		}

		if err = w.Flush(); err != nil {
			appendPrintError("while flushing downloaded file", err, item)
			continue
		}

		// Write info json
		infoName := fmt.Sprintf("%s.json", name)
		result, err := json.MarshalIndent(item, "", "    ")
		if err != nil {
			appendPrintError("while generating json data", err, item)
			continue
		}

		fj, err := w.Create(infoName)
		if err != nil {
			appendPrintError("while creating json file", err, item)
			continue
		}

		_, err = fj.Write(result)
		if err != nil {
			appendPrintError("while writing json file", err, item)
			continue
		}

		if err = w.Flush(); err != nil {
			appendPrintError("while flushing json file", err, item)
			continue
		}

		// add the current link to the links we already downloaded
		downloaded[item.GetDownloadLink()] = true

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
	GenerateReadme(rm, itemCount, int64(len(failedEntrys)))
	println("\nGenerated README.")

	if len(failedEntrys) > 0 {
		ff, err := w.Create("failed.json")
		if err != nil {
			return err
		}
		byt, err := json.MarshalIndent(&failedEntrys, "", "    ")
		if err != nil {
			return err
		}
		ff.Write(byt)
	}

	if err = w.Flush(); err != nil {
		return err
	}

	return nil
}

func getURLExtension(url string) string {
	ext := path.Ext(url)
	if len(ext) == 0 {
		return ""
	}
	return ext[1:]
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
			_, err := b.WriteRune(item)
			if err != nil {
				panic(err) // This should never happen
			}
		}
	}

	return b.String()
}
