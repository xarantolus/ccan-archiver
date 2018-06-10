package zipfactory

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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

	var l int64 = 0
	// Loop over channel & Download & Pack
	for item := range input {
		if _, contains := downloaded[item.DownloadLink]; item.DownloadLink == "" || contains {
			println("Already have", item.DownloadLink)
			continue
		}
		// Download
		name := fmt.Sprintf("%s/%s.%s", item.Author, item.Name, getUrlExtension(item.DownloadLink))
		println("Downloading", name, "("+strconv.FormatInt(l, 10)+"/3321)")

		body, err := crawler.DoRequest(item.DownloadLink)
		if err != nil {
			println("Error:", err.Error())
			continue
		}

		f, err := w.Create(name)
		if err != nil {
			println("Error:", err.Error())
			continue
		}

		if _, err = io.Copy(f, body); err != nil {
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

		println(" > Success")
		l++
	}

	if err = w.Flush(); err != nil {
		return err
	}
	return nil
}

func getUrlExtension(url string) string {
	split := strings.Split(url, ".")
	return split[len(split)-1]
}
