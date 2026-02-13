package main

import (
	"context"
	"math"
	"net/http"
	"strings" // Added to handle URL cleaning
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ScrapeResult stores the outcome of a single URL attempt.
type ScrapeResult struct {
	URL       string
	Status    int
	IsActive  bool
	PageTitle string
	H1        string
	AllHeaders []string //this slice will hold all the h1 and h2 tags
	Description string
}

// Scraper provides the shared HTTP client for all 50 workers.
type Scraper struct {
	Client *http.Client
}

// NewScraper initializes the scraper with a pointer to a shared client.
func NewScraper() *Scraper {
	return &Scraper{
		Client: &http.Client{Timeout: 15 * time.Second},
	}
}

// Fetch performs the scrape with exponential backoff logic
func (s *Scraper) Fetch(ctx context.Context, url string) ScrapeResult {

	// 🛠️ URL Safety Check: Prevents double "https://" prefixes
	fullURL := url
	if !strings.HasPrefix(url, "http") {
		fullURL = "https://" + url
	}

	for i := 0; i < 3; i++ {
		req, _ := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/121.0.0.0")

		resp, err := s.Client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
	
			// Extract title using goquery
			doc, _ := goquery.NewDocumentFromReader(resp.Body)

			//title := "collecting all H1 and H2 tags
			var subHeaders []string
			doc.Find("h1, h2").Each(func(i int, s *goquery.Selection) {
				text := strings.TrimSpace(s.Text())
				if text != "" {
					subHeaders = append(subHeaders, text)
				}
			})

			title := doc.Find("title").First().Text()
			h1 := doc.Find("h1").First().Text()
			Description, _ := doc.Find("meta[name=description]").Attr("content")

			return ScrapeResult{
				URL:       fullURL,
				Status:    resp.StatusCode,
				IsActive:  true,
				PageTitle: strings.TrimSpace(title),
				H1:        strings.TrimSpace(h1),
				AllHeaders: subHeaders, // the full lists for website like google.com
				Description: strings.TrimSpace(Description),
			}
		}

		// Exponential Backoff: Wait 1s, then 2s before failing
		if i < 2 {
			waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
			time.Sleep(waitTime)
		}
	}

	return ScrapeResult{URL: fullURL, IsActive: false}
}