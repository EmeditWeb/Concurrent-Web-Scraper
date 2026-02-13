package main

import (
	"encoding/json"
	"log"
	"os"
	"fmt"
	"time"
)


// this struct holds the data for each sites
type SiteInfo struct {
	Title string `json:"title"`
	H1 string `json:"h1"`
	AllHeaders []string `json:"all_sub_headers"`  // collects other subheaders found in the sites descriptons
	Description string `json:"description"`
}


// `json:"results"` stores the struct
type Analyzer struct {
	    Results        map[string]SiteInfo`json:"results"`
		TotalCount	  int               `json:"total_count"`
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		Results: make(map[string]SiteInfo),
	}
}

// Process updates the internal state because we are modifying the TotalCount and Results map, we need to use a pointer receiver (*Analyzer) to ensure that the changes are reflected in the original Analyzer instance.
func (a *Analyzer) Process(result ScrapeResult) {
	    if result.IsActive {
			    a.TotalCount++
			    a.Results[result.URL] = SiteInfo{
					Title: result.PageTitle,
					H1: result.H1,
					AllHeaders: result.AllHeaders,  // pass the whole lists.
					Description: result.Description,
				}
			log.Printf("Processed result for [%s] | H1: %s\n", result.URL, result.H1)
		} else {
			log.Printf("Failure: %s was unreachable after 3 attempts\n", result.URL)
		}
	
}

// SaveOuputs handles the .json and .txt files
func (a *Analyzer) SaveOutputs() {
	// Save JSON output
	jsonData, _ := json.MarshalIndent(a, "", "  ")
	os.WriteFile("results.json", jsonData, 0644)


	// save the summary files too
	summary := fmt.Sprintf("Scraper Summary Report\nTime: %s\nTotal Successful Scrapes: %d\n\nDetailed Results:\n", time.Now().Format(time.RFC1123), a.TotalCount)
	for url, siteInfo := range a.Results {
		summary += fmt.Sprintf("- %s: %s\n", url, siteInfo.H1)
	}
	os.WriteFile("summary.txt", []byte(summary), 0644)

	log.Println("Results saved to results.json and summary.txt")
	fmt.Println("Results saved to results.json and summary.txt")
}
