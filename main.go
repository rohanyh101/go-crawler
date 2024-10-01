package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var tokens = make(chan struct{}, 5)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15",
	"Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (Linux; Android 10; SM-G950F Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/79.0.3945.136 Mobile Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:93.0) Gecko/20100101 Firefox/93.0",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Linux; U; Android 4.1.1; en-us; Galaxy Nexus Build/JRO03C) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
	"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 (compatible; curl/7.68.0; +https://curl.se/)",
	"curl/7.68.0",
}

func getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func getRequest(targetURL string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", getRandomUserAgent())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func discoverLinks(response *http.Response) []string {
	if response != nil {
		doc, _ := goquery.NewDocumentFromReader(response.Body)
		foundURLs := []string{}

		if doc != nil {
			doc.Find("a").Each(func(index int, s *goquery.Selection) {
				res, _ := s.Attr("href")
				foundURLs = append(foundURLs, res)
			})
		}
		return foundURLs
	}

	return []string{}
}

func checkRelative(href, baseURL string) string {
	if strings.HasPrefix(href, "/") {
		return fmt.Sprintf("%s%s", baseURL, href)
	}

	return href
}

func resolveRelativeURLs(href, baseURL string) (string, bool) {
	resultHref := checkRelative(href, baseURL)
	baseParse, _ := url.Parse(baseURL)
	resultParse, _ := url.Parse(resultHref)

	if baseParse != nil && resultParse != nil {
		if baseParse.Host == resultParse.Host {
			return resultHref, true
		}

		return "", false
	}

	return "", false
}

func Crawl(targetURL, baseURL string) []string {
	fmt.Println("Crawling", targetURL)

	// acquire a token
	tokens <- struct{}{}

	resp, _ := getRequest(targetURL)
	// release the token after the request is done
	<-tokens

	links := discoverLinks(resp)
	foundURLs := []string{}

	for _, link := range links {
		correctLink, ok := resolveRelativeURLs(link, baseURL)
		if ok {
			if correctLink != "" {
				foundURLs = append(foundURLs, correctLink)
			}
		}
	}

	return foundURLs
}

func main() {
	worklist := make(chan []string)
	var n int
	n++
	domain := "https://www.geeksforgeeks.org"

	// start with the initial domain and add it to the worklist
	go func() {
		worklist <- []string{domain}
	}()

	// slice to keep track of the links that have been seen
	seen := make(map[string]bool)

	// crawl the web concurrently
	// n will keep track of the number of links to crawl and increment it when a new link is found
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++

				go func(link, baseURL string) {
					foundLinks := Crawl(link, baseURL)
					if foundLinks != nil {
						worklist <- foundLinks
					}
				}(link, domain)
			}
		}
	}
}
