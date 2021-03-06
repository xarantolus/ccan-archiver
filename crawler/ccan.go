package crawler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xarantolus/ccan-archiver/zipfactory"
	"golang.org/x/net/html"
)

const (
	// url is the url to the listing that is embedded on http://ccan.de
	url string = "https://ccan.de/cgi-bin/ccan/ccan-view.pl?a=&sc=tm&so=d&nr=250&ac=ty-ti-ni-tm-ca-dc-ev-vo-si&reveal=1&pg=%d"

	// dateFormat is the date format used in the listing
	ccanDateFormat string = "02.01.06 15:04"
)

// CCANItem is an item from the listing at `url`
type CCANItem struct {
	Name          string    `json:"name"`
	Date          time.Time `json:"date"`
	DownloadCount int       `json:"download_count"`
	Author        string    `json:"author"`
	Votes         int       `json:"votes"`
	Category      string    `json:"category"`
	Engine        string    `json:"engine"`
	DownloadLink  string    `json:"download_link"`
}

// Implement zipfactory.Archivable

func (c CCANItem) GetDownloadLink() string {
	return c.DownloadLink
}

func (c CCANItem) GetAuthor() string {
	return c.Author
}

func (c CCANItem) GetName() string {
	return c.Name
}

func (c CCANItem) GetSourceName() string {
	return "CCAN"
}

// CrawlCCAN crawls the entire listing and returns items in the channel - it will not be closed
func CrawlCCAN(output chan zipfactory.Archivable) (errorlist []error) {
	var totalItemsLoaded int
	var pageCounter int

	// Add items that aren't listed on ccan.de, but might be needed - See items.go (they are part of this crawler as the files will be in the right directory to find them easily)
	for _, nonlistedItem := range additionalItems {
		output <- nonlistedItem
	}

	var errorCount = 0
	for {
		var currentPageItemCount int

		fmt.Printf("Fetching %d. page\n", pageCounter+1)
		var pageContent, err = DoRequest(fmt.Sprintf(url, pageCounter))
		if err != nil {
			errorCount++
			errorlist = append(errorlist, fmt.Errorf("Error while downloading listing page %d (try %d/6): %s", pageCounter, errorCount+1, err.Error()))

			if errorCount > 5 {
				log.Fatalln(err)
			}

			log.Println(err)

			time.Sleep(5 * time.Second)
			continue
		}

		node, err := html.Parse(pageContent)
		// Close content after parsing, but ignore errors
		_ = pageContent.Close()

		if err != nil {
			errorCount++
			errorlist = append(errorlist, fmt.Errorf("Failed to parse html, attempt %d: %s", errorCount, err.Error()))

			if errorCount > 5 {
				log.Fatalln(err)
			}

			log.Println(err)

			time.Sleep(5 * time.Second)
			continue
		}
		errorCount = 0

		// First should be ignored
		var tbodyFirst = node.LastChild.LastChild.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild.NextSibling

		for tbodyFirst.NextSibling != nil {
			// print(renderNode(tbodyFirst) + "\n\n")

			currentResult := CCANItem{}

			currentNode := tbodyFirst.FirstChild

			var counter int
			var resValid = true
			for currentNode.NextSibling != nil {

				switch counter {
				case 1:
					{
						var nameNode = currentNode.FirstChild.FirstChild

						currentResult.Name = strings.Trim(renderWithoutTags(nameNode), " ")
						// println("Name:", currentResult.Name)
					}
				case 2:
					{
						if currentNode == nil || currentNode.FirstChild == nil {
							resValid = false
							break // because we don't get a download link
						}

						currLink := ""

						for _, a := range currentNode.FirstChild.Attr {
							if a.Key == "href" {
								currLink = "https://ccan.de/cgi-bin/ccan/" + a.Val
								break
							}
						}

						if currLink == "" || currLink == "https://ccan.de/cgi-bin/ccan/" {
							resValid = false
							break
						}
						currentResult.DownloadLink = currLink
						// println("Link:", currentResult.DownloadLink)
					}
				case 3:
					{
						currentResult.Category = strings.Trim(renderWithoutTags(currentNode.FirstChild.FirstChild), " ")
						// println("Category:", currentResult.Category)
					}
				case 4:
					{
						currentResult.Author = strings.Trim(renderWithoutTags(currentNode.FirstChild.FirstChild), " ")
						// println("Author:", currentResult.Author)
					}
				case 5:
					{
						if currentNode == nil || currentNode.FirstChild == nil || currentNode.FirstChild.FirstChild == nil {
							resValid = false
							break // because we don't get a engine
						}
						currentResult.Engine = strings.Trim(renderWithoutTags(currentNode.FirstChild.FirstChild), " ")
						// println("Engine:", currentResult.Engine)
					}
				case 6:
					{
						votescount, err := strconv.ParseInt(strings.Trim(renderWithoutTags(currentNode.FirstChild), " "), 10, 64)
						if err != nil {
							resValid = false
							break
						}
						currentResult.Votes = int(votescount)
						// println("Votes:", currentResult.Votes)
					}
				case 7:
					{
						dlCount, err := strconv.ParseInt(strings.Trim(renderWithoutTags(currentNode.FirstChild), " "), 10, 64)
						if err != nil {
							resValid = false
							break
						}
						currentResult.DownloadCount = int(dlCount)
						// println("Downloads:", currentResult.DownloadCount)
					}
				case 9:
					{
						date, err := parseCCANDate(strings.Trim(renderNode(currentNode.FirstChild), " "))
						if err != nil {
							resValid = false
							break
						}
						currentResult.Date = date
						// println("Date:", date.Format("Jan 2, 2006 um 15:04"))
					}
				}

				if !resValid {
					break
				}

				currentNode = currentNode.NextSibling
				counter++
			}

			if resValid && currentResult.Author != "" && currentResult.Name != "" && currentResult.Category != "" && currentResult.DownloadLink != "" && currentResult.Engine != "" {
				output <- currentResult
				totalItemsLoaded++
				currentPageItemCount++
			}

			tbodyFirst = tbodyFirst.NextSibling
		}

		// Exit if the page we just loaded was empty
		if currentPageItemCount == 0 {
			break
		}

		pageCounter++
	}
	return
}

// renderNode renders the text of a html node and ignores errors
func renderNode(n *html.Node) string {
	if n == nil {
		return ""
	}
	var buf bytes.Buffer
	w := io.Writer(&buf)
	_ = html.Render(w, n)
	return html.UnescapeString(buf.String())
}

var noTagsRegex = regexp.MustCompile(`<[^>]*>`)

// renderWithoutTags removes all html tags from the rendered string
func renderWithoutTags(node *html.Node) string {
	return noTagsRegex.ReplaceAllString(renderNode(node), "")
}

// parseDate parses the `input` with the assumption that it is formatted as `dateFormat`
// All dates in the listing are formatted as `dateFormat`
func parseCCANDate(input string) (output time.Time, err error) {
	output, err = time.Parse(ccanDateFormat, input)
	return
}

// DoRequest opens the file at the specified url
func DoRequest(url string) (io.ReadCloser, error) {
	client := http.Client{
		Timeout: time.Hour,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Do request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode > 399 {
		return nil, fmt.Errorf("Error in http request: StatusCode is %d", res.StatusCode)
	}

	return res.Body, nil
}
