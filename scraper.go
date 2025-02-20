package main

// Overall goal: Pull county GAA player data from
// https://www.finalwhistle.ie/gaelic/
// and populate the structs in structs.go

// Initial goal - scrape a single county game data from the most recent national league fixture.

import (
	"fmt"
	//"net/http"
	//"net/url"
	//"strings"
	//"time"
	"math/rand"
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

func selectRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func scrape() (string, []SearchResult, error) {
	results := []SearchResult{}

	fixtureURL, err := buildFixtureUrls()
	if err != nil {
		fmt.Printf("Fixture: ", fixtureURL)
	}

	return fixtureURL, results, nil
}

func buildFixtureUrls() (string, error) {
	return playerDataDomain["fixtures"], nil
}
