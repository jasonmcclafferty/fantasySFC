package scraper

// Overall goal: Pull county GAA player data from
// https://www.finalwhistle.ie/gaelic/
// and populate the structs in structs.go

// Initial goal - scrape a single county game data from the most recent national league fixture.
// This will involve scraping the fixture data, then scraping the player data from the fixture URL.
import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var playerDataDomain = map[string]string{
	"base":     "https://www.finalwhistle.ie",
	"fixtures": "https://www.finalwhistle.ie/gaelic/donegal-fixtures-results",
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 6.1; WOW64)",
	"AppleWebKit/537.36 (KHTML, like Gecko)",
	"Chrome/44.0.2403.157 Safari/537.36",
}

func selectRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]

}

// buildFixtureUrls returns the URL to scrape for fixtures
func buildFixtureUrls() (string, error) {
	return playerDataDomain["fixtures"], nil
}

// SearchResult represents a search result item
type SearchResult struct {
	ResultURL   string
	ResultTitle string
}

// Fixture represents a GAA match fixture
type Fixture struct {
	MatchTitle string
	MatchURL   string
}

// Result represents a GAA match result
type Result struct {
	MatchTitle  string
	MatchURL    string
	HomeScore   string
	AwayScore   string
	HomeTeam    string
	AwayTeam    string
	CompletedAt string
}

// MatchData combines both fixture and result information
type MatchData struct {
	Fixture Fixture
	Result  Result
}

// scrapeFixture parses a fixture element from the finalwhistle.ie website
func scrapeFixture(fixtureURL string) ([]Fixture, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create and customize request
	req, err := http.NewRequest("GET", fixtureURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set a user agent to mimic a browser
	req.Header.Set("User-Agent", selectRandomUserAgent())

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	fixtures := []Fixture{}

	// Find all fixture elements (table cells)
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		// Look for the match title element which contains the fixture link we want
		titleElement := s.Find(".sp-event-title")
		if titleElement.Length() > 0 {
			fixture := Fixture{}

			// Extract the match title and URL from the title element
			titleLink := titleElement.Find("a")
			fixture.MatchTitle = titleLink.Text()
			fixture.MatchURL = titleLink.AttrOr("href", "")

			// Only add this fixture if we found a valid URL that contains finalwhistle.ie
			if fixture.MatchURL != "" && strings.Contains(fixture.MatchURL, "finalwhistle.ie") {
				fixtures = append(fixtures, fixture)
			}
		}
	})

	return fixtures, nil
}

// scrapeResults parses match results from the finalwhistle.ie website
func scrapeResults(resultsURL string) ([]Result, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create and customize request
	req, err := http.NewRequest("GET", resultsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set a user agent to mimic a browser
	req.Header.Set("User-Agent", selectRandomUserAgent())

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	results := []Result{}

	// Find all result elements (table cells)
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		// Look for cells that have both match title and results
		titleElement := s.Find(".sp-event-title")
		resultsElement := s.Find(".sp-event-results")

		if titleElement.Length() > 0 && resultsElement.Length() > 0 {
			result := Result{}

			// Extract the match title and URL from the title element
			titleLink := titleElement.Find("a")
			result.MatchTitle = titleLink.Text()
			result.MatchURL = titleLink.AttrOr("href", "")

			// Try to parse the title to get team names
			if matchParts := strings.Split(result.MatchTitle, " v "); len(matchParts) == 2 {
				result.HomeTeam = strings.TrimSpace(matchParts[0])
				result.AwayTeam = strings.TrimSpace(matchParts[1])
			}

			// Extract scores from the results element
			resultsLink := resultsElement.Find("a")
			scoreElements := resultsLink.Find(".sp-result")

			if scoreElements.Length() >= 2 {
				result.HomeScore = strings.TrimSpace(scoreElements.First().Text())
				result.AwayScore = strings.TrimSpace(scoreElements.Eq(1).Text())
			}

			// Extract completion date
			dateElement := s.Find("time.sp-event-date")
			if dateElement.Length() > 0 {
				result.CompletedAt = dateElement.Text()
			}

			// Only add results with valid scores and URLs
			if result.MatchURL != "" &&
				strings.Contains(result.MatchURL, "finalwhistle.ie") &&
				result.HomeScore != "" &&
				result.AwayScore != "" {
				results = append(results, result)
			}
		}
	})

	return results, nil
}

// scrapeMatchData combines fixture and result scraping into one function
func scrapeMatchData(url string) ([]MatchData, error) {
	fixtures, err := scrapeFixture(url)
	if err != nil {
		return nil, fmt.Errorf("error scraping fixtures: %v", err)
	}

	results, err := scrapeResults(url)
	if err != nil {
		return nil, fmt.Errorf("error scraping results: %v", err)
	}

	// Create a map of result URLs for quick lookup
	resultMap := make(map[string]Result)
	for _, result := range results {
		resultMap[result.MatchURL] = result
	}

	// Combine fixtures with their corresponding results
	matchData := []MatchData{}
	for _, fixture := range fixtures {
		data := MatchData{
			Fixture: fixture,
		}

		// If we have result data for this fixture, add it
		if result, exists := resultMap[fixture.MatchURL]; exists {
			data.Result = result
		}

		matchData = append(matchData, data)
	}

	return matchData, nil
}

// extractFixtureDetailsFromHTML parses the fixture details from an HTML string
func extractFixtureDetailsFromHTML(htmlContent string) (Fixture, error) {
	fixture := Fixture{}

	// Create a goquery document from the HTML string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return fixture, fmt.Errorf("error parsing HTML string: %v", err)
	}

	// Find the match title link - this is the main link we're interested in
	titleElement := doc.Find(".sp-event-title")
	if titleElement.Length() == 0 {
		return fixture, fmt.Errorf("couldn't find fixture title element")
	}

	titleLink := titleElement.Find("a")
	fixture.MatchTitle = titleLink.Text()
	fixture.MatchURL = titleLink.AttrOr("href", "")

	if fixture.MatchURL == "" {
		return fixture, fmt.Errorf("couldn't find fixture URL")
	}

	return fixture, nil
}

// extractResultDetailsFromHTML parses the match result details from an HTML string
func extractResultDetailsFromHTML(htmlContent string) (Result, error) {
	result := Result{}

	// Create a goquery document from the HTML string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return result, fmt.Errorf("error parsing HTML string: %v", err)
	}

	// Find the match title element
	titleElement := doc.Find(".sp-event-title")
	if titleElement.Length() == 0 {
		return result, fmt.Errorf("couldn't find match title element")
	}

	// Extract the match title and URL
	titleLink := titleElement.Find("a")
	result.MatchTitle = titleLink.Text()
	result.MatchURL = titleLink.AttrOr("href", "")

	// Try to parse the title to get team names
	if matchParts := strings.Split(result.MatchTitle, " v "); len(matchParts) == 2 {
		result.HomeTeam = strings.TrimSpace(matchParts[0])
		result.AwayTeam = strings.TrimSpace(matchParts[1])
	}

	// Find and extract the scores
	resultsElement := doc.Find(".sp-event-results")
	if resultsElement.Length() == 0 {
		return result, fmt.Errorf("couldn't find results element")
	}

	resultsLink := resultsElement.Find("a")
	scoreElements := resultsLink.Find(".sp-result")

	if scoreElements.Length() >= 2 {
		result.HomeScore = strings.TrimSpace(scoreElements.First().Text())
		result.AwayScore = strings.TrimSpace(scoreElements.Eq(1).Text())
	} else {
		return result, fmt.Errorf("couldn't find score elements")
	}

	// Extract completion date
	dateElement := doc.Find("time.sp-event-date")
	if dateElement.Length() > 0 {
		result.CompletedAt = dateElement.Text()
	}

	if result.MatchURL == "" {
		return result, fmt.Errorf("couldn't find match URL")
	}

	return result, nil
}

// extractMatchDataFromHTMLSnippet is a utility function to parse a single HTML snippet
// and extract both fixture and result information
func extractMatchDataFromHTMLSnippet(htmlSnippet string) (MatchData, error) {
	matchData := MatchData{}

	// Extract fixture information
	fixture, err := extractFixtureDetailsFromHTML(htmlSnippet)
	if err != nil {
		return matchData, fmt.Errorf("error extracting fixture: %v", err)
	}
	matchData.Fixture = fixture

	// Extract result information
	result, err := extractResultDetailsFromHTML(htmlSnippet)
	if err != nil {
		// Not treating this as a fatal error - we might have a fixture without results
		fmt.Printf("Note: couldn't extract result: %v\n", err)
	} else {
		matchData.Result = result
	}

	return matchData, nil
}

// Update the existing scrape function to use our new match data processing
func scrape() (string, []SearchResult, error) {
	searchResults := []SearchResult{}

	url, err := buildFixtureUrls() // Reusing this function as it returns the URL we need
	if err != nil {
		return url, searchResults, fmt.Errorf("error building URL: %v", err)
	}

	matchData, err := scrapeMatchData(url)
	if err != nil {
		return url, searchResults, fmt.Errorf("error scraping match data: %v", err)
	}

	// Convert match data to search results for compatibility with existing code
	for _, data := range matchData {
		// Create result text if available
		resultText := ""
		if data.Result.HomeScore != "" && data.Result.AwayScore != "" {
			resultText = fmt.Sprintf(" (%s %s - %s %s)",
				data.Result.HomeTeam,
				data.Result.HomeScore,
				data.Result.AwayTeam,
				data.Result.AwayScore)
		}

		searchResults = append(searchResults, SearchResult{
			ResultURL:   data.Fixture.MatchURL,
			ResultTitle: data.Fixture.MatchTitle + resultText,
		})
	}

	return url, searchResults, nil
}

// cleanTeamName removes the logo and extra whitespace from a team name
func cleanTeamName(name string) string {
	// Remove any text between <span> and </span> tags
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
