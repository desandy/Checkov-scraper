package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strings"
)

func main() {

	var url string
	flag.StringVar(&url, "url", "", "The full Checkov Guide URL")
	flag.Parse()
	// go run checkov.go -url https://docs.bridgecrew.io/docs/s3_16-enable-versioning

	// Ensure the URL is not empty
	// And contains http:// or https://
	if len(url) == 0 || (!strings.Contains(url, "http://") &&
		!strings.Contains(url, "https://")) {
		fmt.Println("Usage: main.go -url")
		fmt.Println("The -url option requires a protocol declaration (https:// or http://)")
		os.Exit(1)
	}

	severity := GetSeverityByURL(url)

	// Only return a notification for success
	if len(severity) > 0 {
		fmt.Printf("[+] Returned: [%s]", severity)
	}
}

func GetSeverityByURL(url string) string {

	text := ""
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("[-] Error getting URL: [%s]", url)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("[-] Status code: [%s] Status: [%s]", res.StatusCode, res.Status)
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Printf("[-] Error Creating Document: [%s]", err)
		return ""
	}

	doc.Find("div.markdown-body").Each(func(i int, s *goquery.Selection) {
		p := s.Find("p")
		p.Contents().Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "Severity:") {
				text = s.Text()
			}
		})
	})

	if len(text) == 0 {
		fmt.Printf("[-] No text found for selector. Url: [%s]", url)
		return ""
	}

	return strings.TrimSpace(text)
}
