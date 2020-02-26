package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseStartURL get the domain name for relative URL parsing later
func ParseStartURL(u string) string {
	parsed, _ := url.Parse(u)
	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
}

// ExtractText get the intresting part of the HTML with goquery
func ExtractText(htmlTag string, doc *goquery.Document) string {
	text := ""
	if doc != nil {
		doc.Find(htmlTag).Each(func(i int, s *goquery.Selection) {
			text, _ = s.Html()

			text = strings.ReplaceAll(text, "<br/>", "\n")
			text = rmEmptyLine(text)
		})
	}
	return text
}

// rmEmptyLine a helper function to delete empty lines in the strings.
func rmEmptyLine(text string) string {
	var newText string
	r := strings.NewReader(text)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			newText += (line + "\n")
		}
	}
	return newText
}

// ExtractLinks in the HTML with goquery
func ExtractLinks(doc *goquery.Document) []string {
	foundUrls := []string{}
	if doc != nil {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			res, _ := s.Attr("href")
			foundUrls = append(foundUrls, res)
		})
		return foundUrls
	}
	return foundUrls
}

// ResolveRelative combine the relative path with the domain rootpath
func ResolveRelative(baseURL string, hrefs []string) []string {
	internalUrls := []string{}

	for _, href := range hrefs {
		if strings.HasPrefix(href, baseURL) {
			internalUrls = append(internalUrls, href)
		}

		if strings.HasPrefix(href, "/") {
			resolvedURL := fmt.Sprintf("%s%s", baseURL, href)
			internalUrls = append(internalUrls, resolvedURL)
		}
	}
	return internalUrls
}

// URLDepth get the URL depth of the current URL compared with the input rootpath
func URLDepth(startURL string, currentURL string) (depth int) {
	rootSubs := strings.Split(strings.Trim(startURL, "/"), "/")
	currentSubs := strings.Split(strings.Trim(currentURL, "/"), "/")
	// return 0 if the currentURL is at the same level as the startURL
	if len(rootSubs) > len(currentSubs) {
		depth = -1
		return
	}
	for i := len(rootSubs) - 2; i >= 0; i-- {
		if rootSubs[i] != currentSubs[i] {
			depth = -1
			return
		}
	}
	depth = len(currentSubs) - len(rootSubs)
	return
}

// WriteToFile write string to file
func WriteToFile(folder, text string, filename string) error {
	filename = strings.TrimSpace(filename)
	f, err := os.Create(filepath.Join(folder, filename+".txt"))
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	_, err = w.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}

// CombineFiles will read the files and combine them to a single file then remove the original file
func CombineFiles(srcFolder string, outputPath string) {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	files, err := ioutil.ReadDir(srcFolder)

	w := bufio.NewWriter(outputFile)
	defer w.Flush()

	for _, file := range files {
		filepath := filepath.Join(srcFolder, file.Name())
		dat, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(dat)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove(filepath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
