package scraper

// Overall goal: Pull county GAA player data from
// https://www.finalwhistle.ie/gaelic/
// and populate the structs in structs.go

// Initial goal - scrape a single county game data from the most recent national league fixture.

import (
	"fmt"
	"net/http"
	"strings"
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

type Fixture struct {
	Date        time.Time
	HomeTeam    string
	AwayTeam    string
	Time        string
	Competition string
	Venue       string
	MatchDay    string
	FixtureURL  string
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
	fixtures := []Fixture{}

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

	// Print the page title for debugging
	fmt.Println("Page title:", doc.Find("title").Text())

	// Find the fixtures table
	fixturesTable := doc.Find("table.sp-event-list")
	fmt.Printf("Found %d fixtures tables\n", fixturesTable.Length())

	// Extract fixtures from the table
	fixturesTable.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		// Create a new fixture
		fixture := Fixture{}

		// Extract date
		dateText := row.Find("td.data-date a").Text()
		dateStr := strings.TrimSpace(row.Find("td.data-date date").Text())
		if dateStr == "" {
			// If the date element isn't found directly, try to parse it from the text
			dateStr = strings.TrimSpace(dateText)
			// You might need to convert the date format
		}

		// Parse the date string to time.Time
		// The format depends on the actual date string format
		// Example format: "2025-03-23 15:45:00"
		parsedDate, err := time.Parse("2006-01-02 15:04:05", dateStr)
		if err != nil {
			fmt.Printf("Error parsing date '%s': %v\n", dateStr, err)
		} else {
			fixture.Date = parsedDate
		}

		// Extract home team
		fixture.HomeTeam = strings.TrimSpace(row.Find("td.data-home a").Text())

		// Extract time/results
		fixture.Time = strings.TrimSpace(row.Find("td.data-time a").Text())

		// Extract away team
		fixture.AwayTeam = strings.TrimSpace(row.Find("td.data-away a").Text())

		// Extract competition
		fixture.Competition = strings.TrimSpace(row.Find("td.data-league a").Text())

		// Extract venue
		fixture.Venue = strings.TrimSpace(row.Find("td.data-venue a").Text())

		// Extract match day
		fixture.MatchDay = strings.TrimSpace(row.Find("td.data-day").Text())

		// Extract fixture URL
		fixture.FixtureURL, _ = row.Find("td.data-date a").Attr("href")

		// Clean up the texts to remove any extra whitespace and logos
		fixture.HomeTeam = cleanTeamName(fixture.HomeTeam)
		fixture.AwayTeam = cleanTeamName(fixture.AwayTeam)

		fixtures = append(fixtures, fixture)

		fmt.Printf("Extracted fixture: %s vs %s on %s at %s\n",
			fixture.HomeTeam, fixture.AwayTeam, fixture.Date.Format("2006-01-02"), fixture.Venue)
	})

	return fixtureURL, results, nil
}

// cleanTeamName removes the logo and extra whitespace from a team name
func cleanTeamName(name string) string {
	// Remove any text between <span> and </span> tags (simplified approach)
	cleaned := strings.Split(name, "<span")[0]

	// Trim any remaining whitespace
	return strings.TrimSpace(cleaned)
}

/* For testing/debugging
func main() {
	fmt.Println("Starting scraper...")
	fixtureURL, fixtures, err := scrape()
	if err != nil {
		log.Fatalf("Error scraping fixtures: %v", err)
	}

	fmt.Printf("Successfully scraped %d fixtures from %s\n", len(fixtures), fixtureURL)

	// If you want to further scrape player data from a fixture
	if len(fixtures) > 0 {
		fmt.Println("Scraping player data from the first fixture...")
		players, err := scrapeFixtureDetails(fixtures[0].FixtureURL)
		if err != nil {
			log.Fatalf("Error scraping player data: %v", err)
		}

		fmt.Printf("Successfully scraped %d players\n", len(players))
		for _, player := range players {
			fmt.Printf("Player: %s, County: %s, Position: %s, Score: %d\n",
				player.Name, player.County.Name, player.Pos, player.Score)
		}
	}
}
*/
