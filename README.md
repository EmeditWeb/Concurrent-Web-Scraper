# 🌐 Web Scraper built in Go

![Web Scraper Screenshot](Screenshot%20from%202026-02-13%2015-09-41.png)

A high-performance, concurrent web scraper built in **Go** that fetches and analyzes metadata from a list of target websites. The scraper extracts key SEO and structural information including page titles, headings, sub-headers, and meta descriptions and exports the results as structured JSON and a human-readable summary report.

---

## 📌 Purpose

The primary goal of this project is to **automate the collection and analysis of web page metadata** across multiple websites simultaneously. It was built to:

- **Scrape websites concurrently** using Go's goroutine-based worker pool pattern for maximum speed.
- **Extract meaningful page metadata** such as the `<title>`, `<h1>`, `<h2>` headings, and `<meta name="description">` tags from each page.
- **Handle unreliable networks gracefully** with built-in retry logic and exponential backoff.
- **Produce structured output** in both JSON (`results.json`) and plain-text summary (`summary.txt`) formats for easy downstream consumption.

This tool is useful for SEO auditing, competitive analysis, website monitoring, and any scenario where batch metadata extraction from web pages is needed.

---

## 🏗️ Architecture Overview

The project follows a clean, modular architecture split across three Go source files:

```
┌─────────────┐       ┌──────────────┐       ┌──────────────┐
│   main.go   │──────▶│  scraper.go  │──────▶│ analyzer.go  │
│ (Orchestrator)      │  (Fetcher)   │       │ (Processor)  │
└─────────────┘       └──────────────┘       └──────────────┘
       │                     │                      │
       │  Worker Pool        │  HTTP + goquery      │  JSON + TXT
       │  (goroutines)       │  (parsing)           │  (export)
       ▼                     ▼                      ▼
   Channels            ScrapeResult            results.json
   (jobs/results)       (struct)              summary.txt
```

### Data Flow

1. **`main.go`** defines the target URLs and spins up a pool of concurrent workers via goroutines.
2. URLs are dispatched through a **jobs channel** to the workers.
3. Each worker calls **`scraper.go`**'s `Fetch()` method, which performs an HTTP GET request, parses the HTML using [goquery](https://github.com/PuerkitoBio/goquery), and returns a `ScrapeResult`.
4. Results flow back through a **results channel** to the main goroutine.
5. **`analyzer.go`**'s `Process()` method collects successful results and `SaveOutputs()` exports everything to disk.

---

## ✨ Key Features

| Feature | Description |
|---|---|
| **Concurrent Scraping** | Uses a pool of 20 goroutine workers to scrape multiple URLs in parallel |
| **Exponential Backoff** | Retries failed requests up to 3 times with exponential wait (1s → 2s) |
| **Metadata Extraction** | Extracts `<title>`, `<h1>`, all `<h1>`/`<h2>` tags, and `<meta description>` |
| **Custom User-Agent** | Sends a Chrome-like User-Agent header to avoid bot detection |
| **URL Safety** | Automatically handles URL prefix validation to prevent malformed requests |
| **Structured JSON Output** | Exports all scraped data to a well-formatted `results.json` |
| **Summary Report** | Generates a `summary.txt` with a timestamped overview of all results |
| **File Logging** | Logs all activity (successes, failures, session start/end) to `scraper.log` |
| **Context Support** | Uses Go's `context.Context` for request lifecycle management |
| **Graceful Channel Management** | Properly closes channels after work completion to prevent goroutine leaks |

---

## 📁 Project Structure

```
scraperproject/
├── main.go          # Entry point — sets up logging, worker pool, channels, and orchestrates the pipeline
├── scraper.go       # Scraper struct & Fetch() — HTTP client, HTML parsing, retry logic with backoff
├── analyzer.go      # Analyzer struct — processes results and exports to JSON and TXT
├── go.mod           # Go module definition and dependencies
├── go.sum           # Dependency checksums
├── results.json     # Output — structured JSON with all scraped metadata
├── summary.txt      # Output — human-readable summary report
├── scraper.log      # Log file — records session activity, successes, and failures
└── README.md        # This file
```

---

## 🔧 Prerequisites

- **Go** 1.24+ installed ([Download Go](https://go.dev/dl/))
- Internet connection (to scrape target URLs)

***

## 🚀 Installation & Usage

### 1. Clone the Repository

```bash
git clone https://github.com/EmeditWeb/Concurrent-Web-Scraper.git
cd scraperproject
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Run the Scraper

```bash
go run .
```

### 4. View the Results

After execution, check the generated output files:

```bash
# View structured JSON results
cat results.json

# View the summary report
cat summary.txt

# View the activity log
cat scraper.log
```

---

## 📊 Output Formats

### `results.json`

A structured JSON file containing all successfully scraped page metadata:

```json
{
  "results": {
    "https://golang.org": {
      "title": "The Go Programming Language",
      "h1": "Build simple, secure, scalable systems with Go",
      "all_sub_headers": [
        "Build simple, secure, scalable systems with Go",
        "Companies using Go",
        "Try Go",
        "What's possible with Go",
        "Get started with Go"
      ],
      "description": "Go is an open source programming language that makes it simple to build secure, scalable systems."
    }
  },
  "total_count": 4
}
```

### `summary.txt`

A timestamped plain-text report listing each URL and its primary heading:

```
Scraper Summary Report
Time: Fri, 13 Feb 2026 14:42:31 +01
Total Successful Scrapes: 4

Detailed Results:
- https://golang.org: Build simple, secure, scalable systems with Go
- https://github.com: Search code, repositories, users, issues, pull requests...
- https://learn2earn.ng: Become an AI-NativeFull-Stack Developer
```

### `scraper.log`

Activity log recording session lifecycle events:

```
2026/02/13 14:42:25 Scraper started
2026/02/13 14:42:26 Processed result for[https://google.com] H1:
2026/02/13 14:42:28 Processed result for[https://golang.org] H1: Build simple, secure, scalable systems with Go
2026/02/13 14:42:31 Failure: https://chatengine.io was unreachable after 3 attempts
2026/02/13 14:42:31 Results saved to results.json and summary.txt
2026/02/13 14:42:31 Session ended
```

---

## ⚙️ Configuration

To scrape different URLs, edit the `urls` slice in `main.go`:

```go
urls := []string{
    "https://google.com",
    "https://golang.org",
    "https://github.com",
    // Add your URLs here
}
```

You can also adjust these parameters:

| Parameter | Location | Default | Description |
|---|---|---|---|
| `numWorkers` | `main.go` | `20` | Number of concurrent worker goroutines |
| `Timeout` | `scraper.go` | `15s` | HTTP request timeout per attempt |
| Retry count | `scraper.go` | `3` | Maximum number of fetch attempts per URL |
| Backoff base | `scraper.go` | `2^i` seconds | Exponential backoff between retries |

---

## 📦 Dependencies

| Package | Purpose |
|---|---|
| [goquery](https://github.com/PuerkitoBio/goquery) | jQuery-like HTML parsing and DOM traversal |
| [cascadia](https://github.com/andybalholm/cascadia) | CSS selector engine (goquery dependency) |
| [golang.org/x/net](https://pkg.go.dev/golang.org/x/net) | Extended networking support (goquery dependency) |

---

## 🧠 How It Works (Step by Step)

1. **Logging Setup** — A log file (`scraper.log`) is opened in append mode to record all activity.
2. **URL List** — A list of target URLs is defined in `main.go`.
3. **Worker Pool** — 20 goroutine workers are started, each listening on a shared `jobs` channel.
4. **Job Dispatch** — Each URL is sent into the `jobs` channel; the channel is closed after all URLs are dispatched.
5. **Fetching** — Each worker calls `Scraper.Fetch()`, which:
   - Validates the URL prefix
   - Sends an HTTP GET request with a custom User-Agent
   - On success (HTTP 200), parses the HTML with goquery to extract metadata
   - On failure, retries up to 3 times with exponential backoff (1s, 2s)
6. **Result Processing** — The `Analyzer.Process()` method receives each `ScrapeResult` and stores it if the fetch was successful.
7. **Export** — `Analyzer.SaveOutputs()` writes `results.json` (structured data) and `summary.txt` (readable report).
8. **Session End** — The session is logged and the program exits.

---

## 📝 License

This project is open source. Feel free to use, modify, and distribute.

---

## 👤 Author

**Emmanuel Itighise** — [https://www.github.com/EmeditWeb](https://www.github.com/EmeditWeb)
