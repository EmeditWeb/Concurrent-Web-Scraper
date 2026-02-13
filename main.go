package main

import (
	"context"
	"log"
	"os"
	"sync"

)

// the worker takes the urls rom the jobs channel and sends the results to the results channel
func worker(ctx context.Context, wg *sync.WaitGroup, s *Scraper, jobs <-chan string, results chan<- ScrapeResult) {
	defer wg.Done()
	for url := range jobs {
		results <- s.Fetch(ctx, url)
	}
}

func main() {

	// setup the logging
	logFile, _ := os.OpenFile("scraper.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()
	log.SetOutput(logFile) // redirects log output to the file
	log.Println("Scraper started")

	// list of urls to scrape
	urls := []string{
		"https://google.com",
		"https://golang.org",
		"https://github.com",
		"https://stackoverflow.com",
		"https://chatengine.io",
		"https://learn2earn.ng",
		"https://www.github.com/EmeditWeb",
		"https://leetcode.com",
	}

	// initialize the pointer to the scraper and analyzer.
	bot := NewScraper()
	analyst := NewAnalyzer()

	// create channels for jobs and results
	jobs := make(chan string, len(urls))
	results := make(chan ScrapeResult, len(urls))

	// create a wait group to wait for all workers to finish
	var wg sync.WaitGroup

	// start the workers
	numWorkers := 20
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(context.Background(), &wg, bot, jobs, results)
	}

	// send the urls to the jobs channel
	for _, url := range urls {
		jobs <- url
	}
	close(jobs) // close the jobs channel after sending all urls

	// wait for all workers to finish
	go func() {
		wg.Wait()
		close(results) // close the results channel after all workers are done
	}()

	// process the results as they come in and analyse them jhhh
	for res := range results {
		analyst.Process(res)
	}

	// export the results to a file
	analyst.SaveOutputs()
	log.Println("Session ended")
}
