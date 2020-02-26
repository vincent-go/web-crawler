// web crawler to crawl book website for personnal reading only
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"github.com/vincent-go/web-crawler/crawler"
	"github.com/vincent-go/web-crawler/types"
	"github.com/vincent-go/web-crawler/util"
)

var (
	// the folder to contain the page result temporarily, with will deleted after combine the files into one
	tempFolder string = "temp"
	// the output folder
	bookFolder string = "books"
	// the output
	filename string = "MyBook.txt"
	// URL to start the web crawler, this URL has depth = 0
	startURL string = ""
	// the htmlTage that the user is intersted to
	htmlTag string = ""
)

var (
	// how many URL to craw at the same time
	concurreny int = 6
	// the depth the crawler can go, if depth = 0, it will crawl the all URL at the same depth with the start URL
	depth int = 1
)

func init() {
	_, err := os.Stat(tempFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(tempFolder, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	} else {
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = os.Stat(bookFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(bookFolder, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	} else {
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	runCrawler()
	outputPath := filepath.Join(bookFolder, filename)
	util.CombineFiles(tempFolder, outputPath)
}

func runCrawler() {
	currentParser := bookParser(yushuwuParserFunc)
	results := crawler.Crawl(startURL, currentParser, tempFolder, concurreny, depth)
	fmt.Println(len(results), " files in total")
}

type bookParser func(*goquery.Document) types.ScrapeResult

func yushuwuParserFunc(d *goquery.Document) types.ScrapeResult {
	var result types.ScrapeResult
	result.URL = d.Url.Path
	result.Title = d.Find("title").Text()
	result.Text = util.ExtractText(htmlTag, d)
	return result
}

func (f bookParser) ParsePage(d *goquery.Document) types.ScrapeResult {
	return f(d)
}
