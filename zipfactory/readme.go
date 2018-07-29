package zipfactory

//go:generate go-bindata -pkg zipfactory -o readme_template.go README.md.tmpl

import (
	"io"
	"text/template"
	"time"
)

const (
	readmeOutputDateFormat = "January 02, 2006"
)

type readmeData struct {
	Count      int64
	DateString string
}

// GenerateReadme Writes the README file to `w`
func GenerateReadme(w *io.Writer, itemCount int64) {
	readme, err := readmeMdTmplBytes()
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("readme").Parse(string(readme))
	if err != nil {
		panic(err)
	}

	ds := time.Now().Format(readmeOutputDateFormat)

	if err := tmpl.Execute(*w, readmeData{
		Count:      itemCount,
		DateString: ds,
	}); err != nil {
		panic(err)
	}
}
