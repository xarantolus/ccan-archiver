package zipfactory

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"../crawler"
)

var (
	downloaded = make(map[string]bool)
)

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
		Timeout: 600 * time.Second,
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		currentDirectURL = req.URL.String()
		return nil
	}

	var l int64
	// Loop over channel & Download & Pack
	for item := range input {
		if _, contains := downloaded[item.DownloadLink]; item.DownloadLink == "" || contains {
			println("Already have", item.DownloadLink)
			continue
		}

		var resp, err = client.Get(item.DownloadLink)
		if err != nil {
			println("Error:", err.Error())
		}

		item.DirectLink = currentDirectURL
		// Download
		name := fmt.Sprintf("%s/%s.%s", item.Author, item.Name, getURLExtension(currentDirectURL))
		println("Downloading", name, "("+strconv.FormatInt(l, 10)+"/3321)")

		f, err := w.Create(name)
		if err != nil {
			println("Error:", err.Error())
			continue
		}

		if _, err = io.Copy(f, resp.Body); err != nil {
			println("Error:", err.Error())
			continue
		}

		if err = w.Flush(); err != nil {
			println("Error:", err.Error())
			continue
		}

		// Write info json

		infoName := fmt.Sprintf("%s.json", name)
		result, err := json.MarshalIndent(item, "", "    ")
		if err != nil {
			println("Error:", err.Error())
			continue
		}

		fj, err := w.Create(infoName)
		if err != nil {
			println("Error:", err.Error())
			continue
		}

		fj.Write(result)

		if err = w.Flush(); err != nil {
			println("Error:", err.Error())
			continue
		}

		downloaded[item.DownloadLink] = true

		currentDirectURL = ""

		println(" > Success")
		l++
	}

	if err = w.Flush(); err != nil {
		return err
	}
	return nil
}

func getURLExtension(url string) string {
	split := strings.Split(url, ".")
	return split[len(split)-1]
}
