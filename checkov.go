package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
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

	var text []string
	// Limit the domain to guide results
	c := colly.NewCollector(
		colly.AllowedDomains("docs.bridgecrew.io"),
	)

	// Grab the 2nd paragraph after the div class markdown-body
	c.OnHTML("div.markdown-body p+p", func(e *colly.HTMLElement) {
		s := strings.Split(e.Text, "\n")
		text = append(text, s...)

	})

	err := c.Visit(url)
	if err != nil {
		fmt.Errorf("[-] Error: [%s] Url: [%s]", err, url)
		return ""
	}
	
	// Ensure the returned string contains the correct parsing
	if !strings.Contains(text[2], "Severity:") {
		fmt.Errorf("[-] Bad parsing for URL: [%s]", url)
		return ""
	}

	return text[2]
}
