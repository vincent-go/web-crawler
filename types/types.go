package types

import (
	"github.com/PuerkitoBio/goquery"
)

type ScrapeResult struct {
	URL   string
	Title string
	Text  string
}

type Parser interface {
	ParsePage(*goquery.Document) ScrapeResult
}
