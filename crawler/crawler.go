package crawler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	//url string = "https://ccan.de/cgi-bin/ccan/ccan-view.pl?a=&sc=tm&so=d&nr=3567&pg=0&ac=ty-ti-ni-tm-ca-dc-ev-vo-si"
	url         string = "https://ccan.de/cgi-bin/ccan/ccan-view.pl?a=&sc=tm&so=d&nr=30&ac=ty-ti-ni-tm-ca-dc-ev-vo-si&pg=%d"
	count              = 3567
	date_format string = "02.01.06 15:04"
)

type CCANItem struct {
	Name          string
	Date          time.Time
	DownloadCount int
	Author        string
	Votes         int
	Category      string
	Engine        string
	DownloadLink  string
}

func CrawlPage(output chan CCANItem) {
	var totalItemsLoaded int

	var pageCounter int
	for {
		var pageContent, err = DoRequest(fmt.Sprintf(url, pageCounter))
		if err != nil {
			log.Fatalln(err)
		}

		var node, ok = html.Parse(pageContent)
		if ok != nil {
			log.Fatalln(err)
		}
		// First should be ignored
		var tbodyFirst = node.LastChild.LastChild.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild.NextSibling

		for tbodyFirst.NextSibling != nil {
			// print(renderNode(tbodyFirst) + "\n\n")

			currentResult := CCANItem{}

			currentNode := tbodyFirst.FirstChild

			var counter int
			for currentNode.NextSibling != nil {

				switch counter {
				case 1:
					{
						currentResult.Name = strings.Trim(renderNode(currentNode.FirstChild.FirstChild), " ")
						// println("Name:", currentResult.Name)
					}
				case 2:
					{
						if currentNode == nil || currentNode.FirstChild == nil {
							continue // because we don't get a download link
						}

						currLink := ""

						for _, a := range currentNode.FirstChild.Attr {
							if a.Key == "href" {
								currLink = "https://ccan.de/cgi-bin/ccan/" + a.Val
								break
							}
						}
						currentResult.DownloadLink = currLink
						// println("Link:", currentResult.DownloadLink)
					}
				case 3:
					{
						currentResult.Category = strings.Trim(renderNode(currentNode.FirstChild.FirstChild), " ")
						// println("Category:", currentResult.Category)
					}
				case 4:
					{
						currentResult.Author = strings.Trim(renderNode(currentNode.FirstChild.FirstChild), " ")
						// println("Author:", currentResult.Author)
					}
				case 5:
					{
						if currentNode == nil || currentNode.FirstChild == nil || currentNode.FirstChild.FirstChild == nil {
							continue // because we don't get a download link
						}
						currentResult.Engine = strings.Trim(renderNode(currentNode.FirstChild.FirstChild), " ")
						// println("Engine:", currentResult.Engine)
					}
				case 6:
					{
						votescount, err := strconv.ParseInt(strings.Trim(renderNode(currentNode.FirstChild), " "), 10, 64)
						if err != nil {
							panic(err)
						}
						currentResult.Votes = int(votescount)
						// println("Votes:", currentResult.Votes)
					}
				case 7:
					{
						dlCount, err := strconv.ParseInt(strings.Trim(renderNode(currentNode.FirstChild), " "), 10, 64)
						if err != nil {
							panic(err)
						}
						currentResult.DownloadCount = int(dlCount)
						// println("Downloads:", currentResult.DownloadCount)
					}
				case 9:
					{
						date, err := parseDate(strings.Trim(renderNode(currentNode.FirstChild), " "))
						if err != nil {
							panic(err)
						}
						currentResult.Date = date
						// println("Date:", date.Format("Jan 2, 2006 um 15:04"))
					}
				}

				currentNode = currentNode.NextSibling
				counter++
			}

			output <- currentResult

			tbodyFirst = tbodyFirst.NextSibling
		}

		// Exit condition
		if totalItemsLoaded == count {
			break
		}

		pageCounter++
	}
	close(output)
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return html.UnescapeString(buf.String())
}

// Parse date format
func parseDate(input string) (output time.Time, err error) {
	output, err = time.Parse(date_format, input)
	return
}

// Request helper
func DoRequest(link string) (io.Reader, error) {

	client := http.Client{
		Timeout: time.Second * 600,
	}

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}

	// Do request
	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	return res.Body, nil
}
