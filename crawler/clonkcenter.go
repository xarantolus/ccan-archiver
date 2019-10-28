package crawler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/iorlas/whitefriday"

	"github.com/xarantolus/ccan-archiver/zipfactory"
)

const (
	urlTemplate         = "https://cc-archive.lwrl.de/download.php?act=getinfo&dl=%d"
	ccDateFormat string = "02.01.2006 15:04:05" // e.g. 30.06.2004 21:02

	maxItemID = 643 // The newest item has the id 643, and there don't seem to be more after it (404 for all above it)
)

var (
	downloadCountRe = regexp.MustCompile(`\((\d+) mal runtergeladen\)`)
)

// CCItem is an item from Clonk Center - Everything will be downloaded from the archive at https://cc-archive.lwrl.de/
type CCItem struct {
	Name          string    `json:"name"`
	Date          time.Time `json:"date"`
	Author        string    `json:"author"`
	PostedBy      string    `json:"posted_by"`
	DownloadCount int       `json:"download_count"`
	Engine        string    `json:"engine"`
	DownloadLink  string    `json:"download_link"`
	Description   string    `json:"description"`
}

// Implement zipfactory.Archivable

func (c CCItem) GetDownloadLink() string {
	return c.DownloadLink
}

func (c CCItem) GetAuthor() string {
	return c.Author
}

func (c CCItem) GetName() string {
	return c.Name
}

func (c CCItem) GetSourceName() string {
	return "Clonk-Center"
}

// CrawlClonkCenter gets all items by incrementing a number and returning the items at the corresponding urls - it doesn't close the `output` channel
func CrawlClonkCenter(output chan zipfactory.Archivable) (errorlist []error) {
	var currentItemID = 1 // 0 will return 404

	for currentItemID < maxItemID+1 {
		item, err := GetClonkCenterItem(currentItemID)
		if err != nil {
			errorlist = append(errorlist, fmt.Errorf("Error while downloading page %d: %s", currentItemID, err.Error()))
			currentItemID++
			continue
		}

		output <- item
		currentItemID++

		// Sleep one second in order to respect the server - don't ddos it
		time.Sleep(time.Second)
	}

	return
}

func GetClonkCenterItem(id int) (result CCItem, err error) {
	content, err := DoRequest(fmt.Sprintf(urlTemplate, id))
	if err != nil {
		return result, fmt.Errorf("Error while downloading page %d: %s", id, err.Error())
	}

	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return result, fmt.Errorf("Error while reading page content %d: %s", id, err.Error())
	}
	_ = content.Close()

	doc.Find("table.fullgrid > tbody > *").Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) != "tr" {
			return
		}

		if s.HasClass("header") {
			result.Name = s.Text()
			return
		}

		key := s.Children().First().Text()

		vChild := s.Children()
		if vChild == nil {
			return
		}
		value := vChild.Eq(1)
		if value == nil {
			return
		}

		switch key {
		case "Datum":
			t, err := time.Parse(ccDateFormat, value.Text())
			if err != nil {
				fmt.Printf("Error parsing date \"%s\": %s", value.Text(), err.Error())
				return
			}
			result.Date = t
		case "Autor":
			result.Author = renderWithoutTags(value.Get(0))

			// Sometimes this author (the original clonk author) will have both names, normalize to "Redwolf Design"
			if result.Author == "Matthes Bender/Redwolf Design" {
				result.Author = "Redwolf Design"
			}
		case "Gepostet von":
			result.PostedBy = renderWithoutTags(value.Get(0))
		case "Engine-Version":
			result.Engine = value.Text()
		case "Download":
			link, ok := value.Children().Eq(0).Attr("href")
			if ok {
				result.DownloadLink = "https://cc-archive.lwrl.de/" + link
			}

			matches := downloadCountRe.FindStringSubmatch(value.Text())
			if len(matches) == 2 {
				count, err := strconv.Atoi(matches[1])
				if err == nil {
					result.DownloadCount = count
				}
			}
		case "Beschreibung":
			htmlString, err := value.Html()
			if err == nil && strings.TrimSpace(htmlString) != "" {
				defer func() {
					// Sometimes whitefriday decides to panic. If this happens, we use the default html string
					if r := recover(); r != nil {
						result.Description = htmlString
					}
				}()
				// Convert this html to markdown
				result.Description = whitefriday.Convert(htmlString)
			}
		default:
			if key != "Dateigröße" && key != "Bewertung" {
				fmt.Printf("Unknown field \"%s\" encountered\n", key)
			}
		}

	})

	return result, nil
}
