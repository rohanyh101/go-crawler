# Website Crawler in Go

A high-performance website crawler built in Go using [GoQuery](https://github.com/PuerkitoBio/goquery) and goroutines for concurrent processing. This crawler effectively navigates through an entire domain, collecting URLs while employing various user agents to avoid getting blocked by the domain server.

## Features

- **Concurrency**: Utilizes goroutines for fast and efficient crawling.
- **Random User Agents**: Implements a strategy for random user agents per request to minimize the risk of getting blocked.
- **Domain-Wide Crawling**: Gathers all found URLs across the entire domain.
- **Customizable Scraping**: Supports additional scraping through parsing the response body to meet specific requirements.

## Getting Started

### Prerequisites
 - Basic knowledge of Go programming.
 - Understanding of goroutines, channels, and concurrency in Go.
 - Make sure you have Go installed on your machine. You can download it from [the official Go website](https://golang.org/dl/).

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/website-crawler.git
   cd website-crawler
   ```

2. Install the necessary dependencies:
 - Run the command in the project directory

   ```bash
   go mod tidy
   ```
## Usage
 ###  To run the crawler, execute the following command in your terminal:
  ```bash
  go run main.go
  ```
 - Replace the `domain` variable in `main.go` with the starting point of the domain you want to crawl.

## Configuration
**You can configure the crawler's behavior by modifying the parameters in the `main.go` file. Options include:**

 - Custom user agents
 - Depth of crawling
 - Specific parsing requirements for additional scraping

## Improvements
**Future enhancements could include:**

 - Adding a user-defined struct to scrape required fields from HTML, allowing users to specify which elements to extract.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request with your improvements or bug fixes.
