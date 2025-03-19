package scraper

// Overall goal: Pull county GAA player data from
// https://www.finalwhistle.ie/gaelic/
// and populate the structs in structs.go

// Initial goal - scrape a single county game data from the most recent national league fixture.

import (
	"fmt"
	"net/http"
	"time"

	//"net/url"
	//"strings"
	//"time"
	"math/rand"

	"github.com/PuerkitoBio/goquery"
	//"github.com/PuerkitoBio/goquery"
)

var playerDataDomain = map[string]string{
	"base":     "https://www.finalwhistle.ie",
	"fixtures": "https://www.finalwhistle.ie/gaelic/donegal-fixtures-results",
}

type SearchResult struct {
	ResultURL   string
	ResultTitle string
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 6.1; WOW64)",
	"AppleWebKit/537.36 (KHTML, like Gecko)",
	"Chrome/44.0.2403.157 Safari/537.36",
}

func buildFixtureUrls() (string, error) {
	return playerDataDomain["fixtures"], nil
}

func selectRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func Scrape() (string, []SearchResult, error) {
	results := []SearchResult{}

	fixtureURL, err := buildFixtureUrls()
	if err != nil {
		fmt.Printf("Fixture: %v", fixtureURL)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  true,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	req, err := http.NewRequest("GET", fixtureURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to create request: %w", err)
	}

	// set headers to mimic browser
	req.Header.Set("User-Agent", selectRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("received non-200 response status: %d %s", resp.StatusCode, resp.Status)
	}

	// PArse HTML usnig GoQuery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to parse HTML: %w", err)
	}

	fmt.Println(doc.Html())

	return fixtureURL, results, nil
}
