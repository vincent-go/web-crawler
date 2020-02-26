package crawler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vincent-go/web-crawler/types"
	"github.com/vincent-go/web-crawler/util"
)

func getRequest(url string) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Baiduspider; +http://www.baidu.com/bot.html)")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return res, fmt.Errorf("%v returned empty response", url)
	}
	return res, nil
}

// CrawlPage crawl the targetURL concurrently
func CrawlPage(baseURL, targetURL string, parser types.Parser, token chan struct{}) ([]string, types.ScrapeResult) {

	token <- struct{}{}
	fmt.Println("Requesting: ", targetURL)
	resp, _ := getRequest(targetURL)
	<-token

	doc, _ := goquery.NewDocumentFromResponse(resp)
	pageResults := parser.ParsePage(doc)
	links := util.ExtractLinks(doc)
	foundUrls := util.ResolveRelative(baseURL, links)

	return foundUrls, pageResults
}

// Crawl gets the scrapeResult with CrawlPage function until it hits the depth limitation
func Crawl(startURL string, parser types.Parser, dstFolder string, concurrency int, depth int) []types.ScrapeResult {
	results := []types.ScrapeResult{}
	worklist := make(chan []string)
	var n int
	n++
	var tokens = make(chan struct{}, concurrency)

	go func() {
		worklist <- []string{startURL}
	}()
	seen := make(map[string]bool)
	baseDomain := util.ParseStartURL(startURL)

	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if seen[link] {
				continue
			}
			seen[link] = true
			n++
			go func(baseDomain, link, dstFolder string, parser types.Parser, token chan struct{}) {
				foundLinks, pageResults := CrawlPage(baseDomain, link, parser, token)
				prefix := func(link string) string {
					splits := strings.Split(link, "/")
					prefix := strings.Replace(splits[len(splits)-1], ".html", "-", 1)
					return prefix
				}(link)
				err := util.WriteToFile(dstFolder, pageResults.Text, prefix+pageResults.Title)
				if err != nil {
					log.Panic(err)
				}
				results = append(results, pageResults)
				if foundLinks == nil {
					return
				}
				var inDepthLinks []string
				for _, url := range foundLinks {
					if d := util.URLDepth(startURL, url); d > -1 && d <= depth {
						inDepthLinks = append(inDepthLinks, url)
					}
				}
				if inDepthLinks != nil {
					worklist <- inDepthLinks
				}
			}(baseDomain, link, dstFolder, parser, tokens)
		}
	}
	return results
}
