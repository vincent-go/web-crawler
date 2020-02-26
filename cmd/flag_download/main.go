// web crawler to crawl book website for personnal reading only
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vincent-go/web-crawler/crawler"
	"github.com/vincent-go/web-crawler/types"
	"github.com/vincent-go/web-crawler/util"
)

var (
	// the folder to contain the page result temporarily, with will deleted after combine the files into one
	tempFolder = "temp"
	// the output folder
	bookFolder = "books"
	// the output
	filename = "MyBook.txt"
	// the htmlTage that the user is intersted to
	htmlTag = "#ur1"
)

var (
	// how many URL to craw at the same time
	concurreny = 6
	// the depth the crawler can go, if depth = 0, it will crawl the all URL at the same depth with the start URL
	depth = 1
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
	fmt.Printf("Please input the URL: \n->")
	reader := bufio.NewReader(os.Stdin) //create new reader, assuming bufio imported
	startURL, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	startURL = strings.TrimSpace(startURL)
	fmt.Printf("The start URL is: %v\n--------------------------Crawler running--------------------------\n", startURL)
	currentParser := bookParser(yushuwuParserFunc)
	results := crawler.Crawl(startURL, currentParser, tempFolder, concurreny, depth)
	fmt.Println(len(results), " files in total")
	outputPath := filepath.Join(bookFolder, filename)
	util.CombineFiles(tempFolder, outputPath)
	fmt.Println("------------------------Everything is done------------------------")
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
